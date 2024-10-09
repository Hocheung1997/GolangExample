package main

import (
	"context"
	"fmt"
	"log"

	svr "github.com/Hocheung1997/multiple-services/service"
	"google.golang.org/protobuf/encoding/protojson"
)

func getRepo(client svr.RepoClient, r *svr.RepoGetRequest) (*svr.RepoGetReply, error) {
	return client.GetRepo(context.Background(), r)
}

func createRepoRequest(jsonQuery string) (*svr.RepoGetRequest, error) {
	r := svr.RepoGetRequest{}
	input := []byte(jsonQuery)
	return &r, protojson.Unmarshal(input, &r)
}

func HandleRepoService(c svr.RepoClient, j string) {
	r, err := createRepoRequest(j)
	if err != nil {
		log.Fatalf("Bad user input: %v", err)
	}
	result, err := getRepo(c, r)
	if err != nil {
		log.Fatal(err)
	}
	for _, repo := range result.Repo {
		fmt.Println("repo Id:", repo.Id)
		fmt.Println("repo name:", repo.Name)
	}
}
