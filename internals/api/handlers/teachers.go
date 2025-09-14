package handlers

import (
	"grpc-school-mgnt/internals/models"
	"grpc-school-mgnt/internals/repositories/mongodb"
	"grpc-school-mgnt/pkg/utils"
	pb "grpc-school-mgnt/proto/gen"

	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AddTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {

	newTeachers := make([]*models.Teacher, 0, len(req.GetTeachers()))
	for _, pbTeacher := range req.GetTeachers() {

		if pbTeacher.Id != "" {
			return nil, status.Error(codes.InvalidArgument, "incorrect payload. non-empty field Id are not allowed.")
		}

		newTeacher := mapTeacherPbToModel(pbTeacher)
		newTeachers = append(newTeachers, newTeacher)
	}

	addedTeachers, err := mongodb.AddTeachersToDB(ctx, newTeachers)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Teachers{Teachers: addedTeachers}, nil
}

func (s *Server) GetTeachers(ctx context.Context, req *pb.GetTeachersRequest) (*pb.TeachersResponse, error) {

	filter, err := buildFilter(req.Teacher, &models.Teacher{})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, utils.ErrorHandler(err, "invalid teacher id").Error())
	}

	sortOptions := buildSortOptions(req.GetSortBy())
	skip, limit, page := buildPaginationOptions(req.PageNumber, req.PageSize)

	teachers, totalCount, err := mongodb.GetTeachersFromDB(ctx, sortOptions, skip, limit, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.TeachersResponse{
		Teachers:   &pb.Teachers{Teachers: teachers},
		Total:      int32(totalCount),
		PageSize:   int32(limit),
		PageNumber: page,
		TotalPages: int32((totalCount + limit - 1) / limit),
	}, nil
}

func (c *Server) UpdateTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {

	var modelTeachers []*models.Teacher
	for _, teacher := range req.Teachers {
		if teacher.Id == "" {
			return nil, status.Error(codes.InvalidArgument, "please provide id for every teacher")
		}
		modelTeacher := mapTeacherPbToModel(teacher)
		modelTeachers = append(modelTeachers, modelTeacher)
	}

	updatedTeachers, err := mongodb.ModifyTeachersDB(ctx, modelTeachers)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Teachers{Teachers: updatedTeachers}, nil

}
