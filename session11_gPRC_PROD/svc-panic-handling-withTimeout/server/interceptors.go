package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type wrappedServerStream struct {
	RecvMsgTimeout time.Duration
	grpc.ServerStream
}

func (s wrappedServerStream) SendMsg(m interface{}) error {
	return s.ServerStream.SendMsg(m)
}

func (s wrappedServerStream) RecvMsg(m interface{}) error {
	ch := make(chan error)
	t := time.NewTimer(s.RecvMsgTimeout)
	go func() {
		log.Printf("Waiting to receive a message: %T", m)
		ch <- s.ServerStream.RecvMsg(m)
	}()
	select {
	case <-t.C:
		return status.Error(
			codes.DeadlineExceeded,
			"Deadline exceeded",
		)
	case err := <-ch:
		return err
	}
}

func timeoutStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	serverStream := wrappedServerStream{
		RecvMsgTimeout: 500 * time.Millisecond,
		ServerStream:   stream,
	}
	err := handler(srv, serverStream)
	return err
}

func loggingUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	resp, err := handler(ctx, req)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Print("No metadata")
	}
	log.Printf("Method:%s, Error:%v, Request-Id:%s",
		info.FullMethod,
		err,
		md.Get("Request-Id"),
	)
	return resp, err
}

func loggingStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	serverStream := wrappedServerStream{
		ServerStream: stream,
	}
	err := handler(srv, serverStream)
	ctx := stream.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Print("No metadata")
	}
	log.Printf("Method:%s, Error:%v, Request-Id:%s",
		info.FullMethod,
		err,
		md.Get("Request-Id"),
	)
	return err
}

func metricUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	end := time.Now()
	log.Printf("Method:%s, Duration:%s",
		info.FullMethod,
		end.Sub(start),
	)
	return resp, err
}

func metricStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {

	start := time.Now()
	err := handler(srv, stream)
	end := time.Now()
	log.Printf("Method:%s, Duration:%s",
		info.FullMethod,
		end.Sub(start),
	)
	return err
}

func panicUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered: %v\n", r)
			err = status.Error(
				codes.Internal,
				"Unexpected error happened",
			)
		}
	}()
	resp, err = handler(ctx, req)
	return
}

func panicStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered: %v\n", r)
			err = status.Error(
				codes.Internal,
				"Unexpected error happened",
			)
		}
	}()
	serverStream := wrappedServerStream{
		ServerStream: stream,
	}
	err = handler(srv, serverStream)
	return
}

func timeoutUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	var resp interface{}
	var err error

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	ch := make(chan error)
	go func() {
		resp, err = handler(ctxWithTimeout, req)
	}()
	select {
	case <-ctxWithTimeout.Done():
		cancel()
		err = status.Error(codes.DeadlineExceeded,
			fmt.Sprintf("%s: Deadline exeeded", info.FullMethod))
		return resp, err
	case <-ch:
		fmt.Println("Received error from channel 'ch'")
	}
	return resp, err
}

func clientDisconnectUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handle grpc.UnaryHandler,
) (interface{}, error) {
	var resp interface{}
	var err error

	ch := make(chan error)

	go func() {
		resp, err = handle(ctx, req)
		ch <- err
	}()

	select {
	case <-ctx.Done():
		err = status.Error(codes.Canceled, fmt.Sprintf("%s: Request canceled", info.FullMethod))
	case <-ch:
		log.Println("error in UnaryHandler")
	}
	return resp, err
}

func clientDisconnectStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) (err error) {
	ch := make(chan error)
	go func() {
		err = handler(srv, stream)
		ch <- err
	}()

	select {
	case <-stream.Context().Done():
		err = status.Error(
			codes.Canceled,
			fmt.Sprintf("%s: Request canceled", info.FullMethod),
		)
	case <-ch:
		log.Println("error in Stream Handler")
	}

	return
}
