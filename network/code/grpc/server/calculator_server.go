package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"calculator/calculator"

	"google.golang.org/grpc"
)

type Server struct{}

func (*Server) Sum(ctx context.Context, req *calculator.CalculatorRequest) (*calculator.CalculatorResponse, error) {
	fmt.Printf("Sum function is invoked with %v \n", req)

	a := req.GetA()
	b := req.GetB()

	res := &calculator.CalculatorResponse{
		Result: a + b,
	}

	return res, nil
}

func main() {
	fmt.Println("starting gRPC server...")

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v \n", err)
	}

	grpcServer := grpc.NewServer()
	calculator.RegisterCalculatorServiceServer(grpcServer, &Server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v \n", err)
	}
}
