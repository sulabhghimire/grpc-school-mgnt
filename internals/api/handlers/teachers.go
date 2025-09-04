package handlers

import (
	"grpc-school-mgnt/internals/models"
	"grpc-school-mgnt/internals/repositories/mongodb"
	pb "grpc-school-mgnt/proto/gen"
	"reflect"

	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AddTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {

	newTeachers := make([]*models.Teacher, len(req.GetTeachers()))
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

func mapTeacherPbToModel(pbTeacher *pb.Teacher) *models.Teacher {
	modelTeacher := models.Teacher{}
	pbVal := reflect.ValueOf(pbTeacher).Elem()
	modelVal := reflect.ValueOf(&modelTeacher).Elem()

	for i := range pbVal.NumField() {
		pbField := pbVal.Field(i)
		fieldName := pbVal.Type().Field(i).Name

		modelField := modelVal.FieldByName(fieldName)

		if modelField.IsValid() && modelField.CanSet() {
			modelField.Set(pbField)
		}
	}

	return &modelTeacher
}
