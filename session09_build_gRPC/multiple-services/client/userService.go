package main

import (
	"context"
	"fmt"
	"log"
	"os"

	svr "github.com/Hocheung1997/multiple-services/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

func createUserRequest(jsonQuery string) (*svr.UserGetRequest, error) {
	u := svr.UserGetRequest{}
	input := []byte(jsonQuery)
	return &u, protojson.Unmarshal(input, &u)
}

func getUser(client svr.UsersClient, u *svr.UserGetRequest) (*svr.UserGetReply, error) {
	return client.GetUser(context.Background(), u)
}

func HandleUserService(c svr.UsersClient, j string) {
	u, err := createUserRequest(j)
	if err != nil {
		log.Fatalf("Bad user input: %v", err)
	}
	result, err := getUser(c, u)
	s := status.Convert(err)
	if s.Code() != codes.OK {
		log.Fatalf("Request failed: %v -%v\n", s.Code(), s.Message())
	}
	fmt.Fprintf(
		os.Stdout, "User: %s %s\n",
		result.User.FirstName,
		result.User.LastName,
	)
}
