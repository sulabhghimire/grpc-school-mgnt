package mongodb

import (
	"context"
	"grpc-school-mgnt/internals/models"
	"grpc-school-mgnt/pkg/utils"
	"reflect"

	pb "grpc-school-mgnt/proto/gen"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func AddTeachersToDB(ctx context.Context, newTeachers []*models.Teacher) ([]*pb.Teacher, error) {

	var addedTeachers []*pb.Teacher

	client := Client()
	for _, newTeacher := range newTeachers {
		res, err := client.Database("school-management").Collection("teachers").InsertOne(ctx, newTeacher)
		if err != nil {
			return nil, utils.ErrorHandler(err, "something went wrong")
		}

		objectId, ok := res.InsertedID.(bson.ObjectID)
		if ok {
			newTeacher.Id = objectId.Hex()
		}

		pbTeacher := mapTeacherModelToPb(newTeacher)
		addedTeachers = append(addedTeachers, pbTeacher)
	}
	return addedTeachers, nil
}

func mapTeacherModelToPb(newTeacher *models.Teacher) *pb.Teacher {
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
	return pbTeacher
}
