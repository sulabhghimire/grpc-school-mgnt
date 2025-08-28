package main

import (
	"grpc-school-mgnt/internals/api/handlers"
	pb "grpc-school-mgnt/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	grpcServer := grpc.NewServer()
	pb.RegisterExecsServiceServer(grpcServer, &handlers.Server{})
	pb.RegisterStudentsServiceServer(grpcServer, &handlers.Server{})
	pb.RegisterTeachersServiceServer(grpcServer, &handlers.Server{})

	reflection.Register((grpcServer))

}
