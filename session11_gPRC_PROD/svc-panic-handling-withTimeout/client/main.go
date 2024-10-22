package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	svc "github.com/Hocheung1997/svc-panic-handling-withTimeout/service"
	"google.golang.org/grpc"
)

func setupGrpcConn(addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithChainUnaryInterceptor(
			loggingUnaryInterceptor,
			metadataUnaryInterceptor,
		),
		grpc.WithChainStreamInterceptor(
			loggingStreamingInterceptor,
			metadataStreamingInterceptor,
		),
	)
}

func getUserServiceClient(conn *grpc.ClientConn) svc.UsersClient {
	return svc.NewUsersClient(conn)
}

func getUser(
	client svc.UsersClient,
	u *svc.UserGetRequest,
) (*svc.UserGetReply, error) {
	return client.GetUser(context.Background(), u)
}

func createHelpStream(c svc.UsersClient) (svc.Users_GetHelpClient, error) {
	return c.GetHelp(
		context.Background(),
		grpc.WaitForReady(true),
	)
}

func setupChat(c svc.UsersClient) (err error) {
	var clientConn = make(chan svc.Users_GetHelpClient)
	var done = make(chan bool)

	stream, err := createHelpStream(c)
	defer stream.CloseSend()
	if err != nil {
		return err
	}

	go func() {
		for {
			clientConn <- stream
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true
				break
			}
			if err != nil {
				log.Printf("Recreating stream.")
				stream, err = createHelpStream(c)
				if err != nil {
					close(clientConn)
					done <- true
					break
				}
			} else {
				fmt.Printf("Response: %s\n", resp.Response)
				if resp.Response == "hello-10" {
					done <- true
					break
				}
			}
		}
	}()

	requestMsg := "hello"
	msgCount := 1
	for {
		if msgCount > 10 {
			break
		}
		stream = <-clientConn
		if stream == nil {
			break
		}
		request := svc.UserHelpRequest{
			Request: fmt.Sprintf("%s-%d", requestMsg, msgCount),
		}
		err := stream.Send(&request)
		if err != nil {
			log.Printf("Send error: %v. Will retry.\n", err)
		} else {
			log.Printf("Request sent: %d\n", msgCount)
			msgCount += 1
		}
	}
	<-done
	return stream.CloseSend()
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal(
			"Specify a gRPC server and method to call",
		)
	}
	serverAddr := os.Args[1]
	methodName := os.Args[2]

	if methodName == "GetUser" && len(os.Args) != 4 {
		log.Fatal("Specify an email address for the user")
	}

	conn, err := setupGrpcConn(serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := getUserServiceClient(conn)

	switch methodName {
	case "GetUser":
		userEmail := os.Args[3]
		result, err := getUser(
			c,
			&svc.UserGetRequest{Email: userEmail},
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(
			os.Stdout, "User: %s %s\n",
			result.User.FirstName,
			result.User.LastName,
		)
	case "GetHelp":
		err = setupChat(c)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Unrecognized method name")
	}
}
