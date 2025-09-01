package main

import (
	"fmt"
	"grpc-school-mgnt/internals/api/handlers"
	"grpc-school-mgnt/internals/repositories/mongodb"
	"grpc-school-mgnt/pkg/config"
	"grpc-school-mgnt/pkg/utils"
	pb "grpc-school-mgnt/proto/gen"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	cfg, err := config.GetAppConfig()
	if err != nil {
		panic(err)
	}

	// make mongodb connection
	err = mongodb.Connect(cfg.DB.URI)
	if err != nil {
		panic(err)
	}
	defer mongodb.Disconnect()

	grpcServer := grpc.NewServer()
	pb.RegisterExecsServiceServer(grpcServer, &handlers.Server{})
	pb.RegisterStudentsServiceServer(grpcServer, &handlers.Server{})
	pb.RegisterTeachersServiceServer(grpcServer, &handlers.Server{})

	reflection.Register((grpcServer))

	listner, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		utils.ErrorHandler(err, fmt.Sprintf("Error listening on the specified port %s", cfg.Server.Port))
		return
	}

	log.Printf("The gRPC server is running on port %s\n", cfg.Server.Port)
	err = grpcServer.Serve(listner)
	if err != nil {
		utils.ErrorHandler(err, "Error running gRPC server.")
		return
	}

}
