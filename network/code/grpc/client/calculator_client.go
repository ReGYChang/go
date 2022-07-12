package main

import (
	"calculator/calculator"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	defer conn.Close()

	client := calculator.NewCalculatorServiceClient(conn)

	doUnary(client)
}

func doUnary(client calculator.CalculatorServiceClient) {
	fmt.Println("Staring to do a Unary RPC")
	req := &calculator.CalculatorRequest{
		A: 3,
		B: 10,
	}

	res, err := client.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling CalculatorService: %v \n", err)
	}

	log.Printf("Response from CalculatorService: %v", res.Result)
}
