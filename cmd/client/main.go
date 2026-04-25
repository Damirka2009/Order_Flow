package main

import (
	"context"
	"log"
	"time"

	api "master/pkg/api"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot connect: %v", err)
	}
	defer conn.Close()

	// 2. Создаём клиент
	client := api.NewOrderServiceClient(conn)

	// 3. Контекст (таймаут — очень важно)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ===== CREATE =====
	createResp, err := client.CreateOrder(ctx, &api.CreateOrderRequest{
		Item:  "pizza",
		Price: 100,
	})
	if err != nil {
		log.Fatalf("CreateOrder error: %v", err)
	}
	log.Println("Created:", createResp.Order)

	orderID := createResp.Order.Id

	// ===== GET =====
	getResp, err := client.GetOrder(ctx, &api.GetOrderRequest{
		Id: orderID,
	})
	if err != nil {
		log.Fatalf("GetOrder error: %v", err)
	}
	log.Println("Got:", getResp.Order)

	// ===== LIST =====
	listResp, err := client.OrdersList(ctx, &api.Empty{})
	if err != nil {
		log.Fatalf("ListOrders error: %v", err)
	}
	log.Println("All orders:")
	for _, o := range listResp.Orders {
		log.Printf(" - %+v\n", o)
	}

	// ===== DELETE =====
	_, err = client.DeleteOrder(ctx, &api.DeleteOrderRequest{
		Id: orderID,
	})
	if err != nil {
		log.Fatalf("DeleteOrder error: %v", err)
	}
	log.Println("Deleted order:", orderID)

	// ===== CHECK DELETE =====
	_, err = client.GetOrder(ctx, &api.GetOrderRequest{
		Id: orderID,
	})
	if err != nil {
		log.Println("Expected error after delete:", err)
	}
}
