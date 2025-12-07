package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"syscall"

	grcpsServer "google.golang.org/grpc"
)

var GrpcAddr = ":9093"

func main() {
	inmemRepo := repository.NewInmemRepository()
	svc := service.NewService(inmemRepo)
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()
	Listener, err := net.Listen("tcp", GrpcAddr)

	if err != nil {
		log.Fatalf("could not start server: %v", err)
	}
	server := grcpsServer.NewServer()
	grpc.RegisterTripHandler(server, svc)
	go func() {
		if err := server.Serve(Listener); err != nil {
			log.Fatalf("could not start server: %v", err)
		}
	}()
	ctx.Done()
	log.Printf("shutting down server...")

	server.GracefulStop()
	log.Printf("server stopped Successfully")

}
