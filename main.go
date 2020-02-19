package main

import (
	"context"
	echopb "echo-service/api/echo"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type hello struct {

}

func (h *hello) PingStream(stream echopb.EchoService_PingStreamServer) error {
	messageCount := 5
	var messages string
	for {
		req, err := stream.Recv()

		log.Printf("Recieved a message: %s ", req.GetMessage())
		if err == io.EOF {
			messages = messages + " " + req.GetMessage()
			messageCount--

			if messageCount < 1 {
				_ = stream.SendAndClose(&echopb.PongResponse{
					Message: messages,
				})

			}
		}
		if err != nil {
			return err
		}


	}


	return nil
}

func (h *hello) Ping(_ context.Context, req *echopb.PingRequest) (*echopb.PongResponse, error) {

	if req.GetMessage() == "error" {
		return nil, errors.New("yeah thats right, get an error")
	}

	return &echopb.PongResponse{
		Message:              req.GetMessage(),
	}, nil
}

func main() {
	fmt.Println("yeet")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", "0.0.0.0", "50051"))
	if err != nil {
		log.Fatalf("shit happned %v", err)
	}
	var opts []grpc.ServerOption
	gRPCServer := grpc.NewServer(opts...)

	echopb.RegisterEchoServiceServer(gRPCServer, &hello{})

	err = gRPCServer.Serve(listener)
}
