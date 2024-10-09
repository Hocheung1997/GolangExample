package main

import (
	"context"
	"log"
	"net"
	"testing"

	users "github.com/Hocheung1997/user-service/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type dummyUserService struct {
	users.UnimplementedUsersServer
}

func (s *dummyUserService) GetUser(ctx context.Context, in *users.UserGetRequest) (*users.UserGetReply, error) {
	u := users.User{
		Id:        "user-123-a",
		FirstName: "jane",
		LastName:  "doe",
		Age:       36,
	}
	return &users.UserGetReply{User: &u}, nil
}

func startServer(s *grpc.Server, l *bufconn.Listener) error {
	return s.Serve(l)
}

func startTestGprcServer() (*grpc.Server, *bufconn.Listener) {
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	users.RegisterUsersServer(s, &dummyUserService{})
	go func() {
		err := startServer(s, l)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return s, l
}

func TestGetUser(t *testing.T) {
	// start up the Test server with buffconn, return the listener for client
	s, l := startTestGprcServer()
	defer s.GracefulStop()
	// grpc.WithContextDialer require a function return net.conn, Dailing address in it
	bufconnDialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return l.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(bufconnDialer))
	if err != nil {
		t.Fatal(err)
	}

	c := getUserServiceClient(conn)
	result, err := getUser(
		c,
		&users.UserGetRequest{Email: "jane@doe.com"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if result.User.FirstName != "jane" || result.User.LastName != "doe" {
		t.Fatalf("Expected: jane doe, Got: %s %s", result.User.FirstName, result.User.LastName)
	}
}
