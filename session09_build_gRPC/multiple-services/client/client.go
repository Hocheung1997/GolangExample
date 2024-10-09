package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	svr "github.com/Hocheung1997/multiple-services/service"
	"google.golang.org/grpc"
)

func getUserServiceClient(conn *grpc.ClientConn) svr.UsersClient {
	return svr.NewUsersClient(conn)
}

func getRepoServiceClient(conn *grpc.ClientConn) svr.RepoClient {
	return svr.NewRepoClient(conn)
}

func setupGrpcConnection(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(
		addr,
		grpc.WithInsecure(),
	)

}

type callerConfig struct {
	address  string
	jsonData string
	service  string
}

func (c *callerConfig) creatFlagSet(w io.Writer) *flag.FlagSet {

	fs := flag.NewFlagSet("grpc_client", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&c.service, "service", "", "the URL of RPC server")
	fs.StringVar(&c.address, "address", "", "the URL of RPC server")
	fs.StringVar(&c.jsonData, "jsonData", "", "the data you submit")
	fs.Usage = func() {
		var usageString = `
grpc: A gRPC client.

grpc: <options> server`
		fmt.Fprint(w, usageString)
		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}
	return fs
}

func main() {
	config := callerConfig{}
	fs := config.creatFlagSet(os.Stdout)
	err := fs.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	serverAddr := config.address

	conn, err := setupGrpcConnection(serverAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	switch config.service {
	case "repo":
		c := getRepoServiceClient(conn)
		HandleRepoService(c, config.jsonData)
	case "user":
		c := getUserServiceClient(conn)
		HandleUserService(c, config.jsonData)
	}

}
