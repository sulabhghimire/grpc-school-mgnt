package handlers

import pb "grpc-school-mgnt/proto/gen"

type Server struct {
	pb.UnimplementedExecsServiceServer
	pb.UnimplementedStudentsServiceServer
	pb.UnimplementedTeachersServiceServer
}
