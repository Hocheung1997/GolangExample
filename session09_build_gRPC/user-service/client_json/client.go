package main

import (
	"context"
	"fmt"
	"log"
	"os"

	users "github.com/Hocheung1997/user-service/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func setupGrpcConnection(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(
		addr,
		grpc.WithInsecure(),
	)

}

func getUserServiceClient(conn *grpc.ClientConn) users.UsersClient {
	return users.NewUsersClient(conn)
}

func getUser(client users.UsersClient, u *users.UserGetRequest) (*users.UserGetReply, error) {
	return client.GetUser(context.Background(), u)
}

func createUserRequest(jsonQuery string) (*users.UserGetRequest, error) {
	u := users.UserGetRequest{}
	input := []byte(jsonQuery)
	return &u, protojson.Unmarshal(input, &u)
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Must specify a gRPC server address and search query")
	}
	serverAddr := os.Args[1]
	u, err := createUserRequest(os.Args[2])
	if err != nil {
		log.Fatalf("Bad user input: %v", err)
	}

	conn, err := setupGrpcConnection(serverAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	c := getUserServiceClient(conn)

	result, err := getUser(c, u)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(
		os.Stdout, "User: %s %s\n",
		result.User.FirstName,
		result.User.LastName,
	)
}
