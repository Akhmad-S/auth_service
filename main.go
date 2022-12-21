package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/uacademy/e_commerce/auth_service/config"
	ecom "github.com/uacademy/e_commerce/auth_service/proto-gen/e_commerce"
	srv "github.com/uacademy/e_commerce/auth_service/services/auth"
	"github.com/uacademy/e_commerce/auth_service/storage"
	"github.com/uacademy/e_commerce/auth_service/storage/postgres"
)

func main() {
	cfg := config.Load()

	psqlConString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	)

	var stg storage.StorageI
	stg, err := postgres.InitDb(psqlConString)
	if err != nil {
		panic(err)
	}

	println("gRPC server tutorial in Go")

	listener, err := net.Listen("tcp", ":9003")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	ecom.RegisterAuthServiceServer(s, srv.NewAuthService(cfg, stg))
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
