package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"testing"

	svc "github.com/Hocheung1997/server-streaming/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func startTestGrpcServer() (*grpc.Server, *bufconn.Listener) {
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	registerServices(s)
	go func() {
		err := startServer(s, l)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return s, l
}

func TestRepoService(t *testing.T) {
	s, l := startTestGrpcServer()
	defer s.GracefulStop()
	bufconnDialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return l.Dial()
	}

	client, err := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(bufconnDialer))

	if err != nil {
		t.Fatal(err)
	}
	repoClient := svc.NewRepoClient(client)
	stream, err := repoClient.GetRepos(
		context.Background(),
		&svc.RepoGetRequest{
			CreatorId: "user-123",
			Id:        "repo-123",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	var repos []*svc.Repository
	for {
		repoMsg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		repos = append(repos, repoMsg.Repo)
	}
	if len(repos) != 5 {
		t.Fatalf("Expected to get back 5 repos, got back: %d repos", len(repos))
	}

	for idx, repo := range repos {
		gotRepoName := repo.Name
		expectedRepoName := fmt.Sprintf("repo-%d", idx+1)
		if gotRepoName != expectedRepoName {
			t.Errorf(
				"Expected Repo Name to be: %s, Got %s",
				expectedRepoName,
				gotRepoName,
			)
		}
	}
}
