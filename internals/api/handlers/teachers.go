package handlers

import (
	"grpc-school-mgnt/internals/models"
	"grpc-school-mgnt/internals/repositories/mongodb"
	"grpc-school-mgnt/pkg/utils"
	pb "grpc-school-mgnt/proto/gen"
	"reflect"
	"strings"

	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func (s *Server) GetTeachers(ctx context.Context, req *pb.GetTeachersRequest) (*pb.TeachersResponse, error) {

	filter := buildGetTeachersFilter(req.Teacher)
	sortOptions := buildSortOptions(req.GetSortBy())
	skip, limit, page := buildPaginationOptions(req.PageNumber, req.PageSize)

	client := mongodb.Client()
	collection := client.Database("school-management").Collection("teachers")

	findOptions := options.Find().SetSort(sortOptions).SetSkip(skip).SetLimit(limit)
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, utils.ErrorHandler(err, "something went wrong.")
	}
	defer cursor.Close(ctx)

	var teachers []*pb.Teacher
	for cursor.Next(ctx) {
		var teacher models.Teacher
		err = cursor.Decode(&teacher)
		if err != nil {
			return nil, utils.ErrorHandler(err, "something went wrong.")
		}
		teachers = append(teachers, mongodb.MapTeacherModelToPb(&teacher))
	}

	totalCount, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, utils.ErrorHandler(err, "something went wrong.")
	}

	return &pb.TeachersResponse{
		Teachers:   &pb.Teachers{Teachers: teachers},
		Total:      int32(totalCount),
		PageSize:   int32(limit),
		PageNumber: page,
		TotalPages: int32((totalCount + limit - 1) / limit),
	}, nil
}

func buildGetTeachersFilter(reqfilter *pb.Teacher) bson.M {
	filter := bson.M{}

	var teacher models.Teacher
	modelVal := reflect.ValueOf(&teacher).Elem()
	modelType := modelVal.Type()

	filterVal := reflect.ValueOf(reqfilter).Elem()
	filterType := filterVal.Type()

	for i := range filterVal.NumField() {
		fieldVal := filterVal.Field(i)
		fieldName := filterType.Field(i).Name

		if fieldVal.IsValid() && !fieldVal.IsZero() {
			modelField := modelVal.FieldByName(fieldName)
			if modelField.IsValid() && modelField.CanSet() {
				modelField.Set(fieldVal)
			}
		}
	}

	for i := range modelVal.NumField() {
		fieldVal := modelVal.Field(i)

		if fieldVal.IsValid() && !fieldVal.IsZero() {
			bsonTag := strings.TrimSuffix(modelType.Field(i).Tag.Get("bson"), ",omitempty")
			if strings.TrimSpace(bsonTag) != "" {
				filter[bsonTag] = fieldVal.Interface()
			}
		}
	}

	return filter
}

func buildSortOptions(sortFields []*pb.SortField) bson.D {

	var sortOptions bson.D

	for _, sortField := range sortFields {
		order := 1
		if sortField.GetOrder() == pb.Order_DESC {
			order = -1
		}
		sortOptions = append(sortOptions, bson.E{Key: sortField.Field, Value: order})
	}

	return sortOptions

}

func buildPaginationOptions(page int32, size int32) (skip int64, limit int64, pageNumber int32) {

	limit = int64(size)
	if limit < 1 {
		limit = 10
	}

	if page < 1 {
		page = 1
	}

	skip = (int64(page) - 1) * (limit)

	return skip, limit, page

}
