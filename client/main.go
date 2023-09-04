package main

import (
	"auth-project/data"
	pb "auth-project/protobuf"
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

var addr = flag.String("addr", "localhost:2003", "the address to connect to")

func CallGetMe(client pb.UserServiceClient, id string) {
	ctx,cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := client.GetMe(ctx, &pb.GetMeRequest{Id: id})
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	fmt.Println(res)
}

func main() {
	flag.Parse()

	// Set up the credentials for the connection.
	perRPC := oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(fetchToken())}
	creds, err := credentials.NewClientTLSFromFile(data.Path("x509/ca_cert.pem"), "x.test.example.com")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(perRPC),
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.Dial(*addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	rgc := pb.NewUserServiceClient(conn)

	CallGetMe(rgc, "1")
}

func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "some-secret-token",
	}
}