package main

import (
	"log"
	"net"

	pb "github.com/cloudbiter/guess-number-grpc/proto"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	guessNumberBot := pb.NewGuessNumberBot()
	pb.RegisterGuessNumberGameServer(grpcServer, guessNumberBot)

	grpcServer.Serve(lis)
}
