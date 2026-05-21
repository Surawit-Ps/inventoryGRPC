package main

import (
	"fmt"
	"inventoryService/adapter/grpc"
	"inventoryService/adapter/repository"
	// "inventoryService/core/entity"
	"inventoryService/core/service"
	"inventoryService/proto/pb"
	"log"
	"net"
	// "time"

	// "github.com/go-redis/redis/v8"
	grpcServer "google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Initialize SQLite Database
	db, err := gorm.Open(sqlite.Open("inventory.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.AutoMigrate(&repository.InventoryModel{})

	// Initialize
	// Create Repository
	inventoryRepo := repository.NewInventoryRepository(db)

	// inv := []entity.RoomInventory{
	// 	{RoomTypeId: "deluxe", Date: time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC), Available: 5, Total: 5},
	// 	{RoomTypeId: "deluxe", Date: time.Date(2024, 7, 2, 0, 0, 0, 0, time.UTC), Available: 5, Total: 5},
	// 	{RoomTypeId: "suite", Date: time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC), Available: 3, Total: 3},
	// 	{RoomTypeId: "suite", Date: time.Date(2024, 7, 2, 0, 0, 0, 0, time.UTC), Available: 3, Total: 3},
	// }

	// for _, inventory := range inv {
	// 	err := inventoryRepo.CreateInventory(inventory)
	// 	if err != nil {
	// 		log.Printf("Failed to create inventory: %v", err)
	// 	}
	// }

	// Create Service
	inventoryService := service.NewInventoryService(inventoryRepo, nil) // Pass nil for Redis client for now

	// Create gRPC Handler
	handler := grpc.NewInventoryGRPCHandler(inventoryService)

	// Create gRPC Server
	grpcServer := grpcServer.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, handler)

	// Listen on port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}

	fmt.Println("gRPC Server is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
