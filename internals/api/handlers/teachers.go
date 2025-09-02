package handlers

import (
	"grpc-school-mgnt/internals/models"
	"grpc-school-mgnt/internals/repositories/mongodb"
	"grpc-school-mgnt/pkg/utils"
	pb "grpc-school-mgnt/proto/gen"
	"reflect"

	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Server) AddTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {

	client := mongodb.Client()

	newTeachers := make([]*models.Teacher, len(req.GetTeachers()))
	for i, pbTeacher := range req.GetTeachers() {
		modelTeacher := models.Teacher{}
		pbVal := reflect.ValueOf(pbTeacher).Elem()
		modelVal := reflect.ValueOf(&modelTeacher).Elem()

		for j := range pbVal.NumField() {
			pbField := pbVal.Field(j)
			fieldName := pbVal.Type().Field(j).Name

			modelField := modelVal.FieldByName(fieldName)

			if modelField.IsValid() && modelField.CanSet() {
				modelField.Set(pbField)
			}
		}
		newTeachers[i] = &modelTeacher

	}

	var addedTeachers []*pb.Teacher
	for _, newTeacher := range newTeachers {
		res, err := client.Database("school-management").Collection("teachers").InsertOne(ctx, newTeacher)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error adding data to database")
		}

		objectId, ok := res.InsertedID.(primitive.ObjectID)
		if ok {
			newTeacher.Id = objectId.Hex()
		}

		pbTeacher := &pb.Teacher{}

		modelVal := reflect.ValueOf(*newTeacher)
		pbVal := reflect.ValueOf(pbTeacher).Elem()

		for i := range modelVal.NumField() {
			modelField := modelVal.Field(i)
			modelFieldType := modelVal.Type().Field(i)

			pbField := pbVal.FieldByName(modelFieldType.Name)
			if pbField.IsValid() && pbField.CanSet() {
				pbField.Set(modelField)
			}
		}

		addedTeachers = append(addedTeachers, pbTeacher)
	}

	return &pb.Teachers{Teachers: addedTeachers}, nil
}
