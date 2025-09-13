package mongodb

import (
	"context"
	"fmt"
	"grpc-school-mgnt/internals/models"
	"grpc-school-mgnt/pkg/utils"

	pb "grpc-school-mgnt/proto/gen"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func AddTeachersToDB(ctx context.Context, newTeachers []*models.Teacher) ([]*pb.Teacher, error) {

	var addedTeachers []*pb.Teacher

	client := Client()
	if client == nil {
		return nil, fmt.Errorf("mongo client is nil, did you connect?")
	}
	for _, newTeacher := range newTeachers {

		if newTeacher == nil {
			return nil, fmt.Errorf("received nil teacher")
		}

		res, err := client.Database("school-management").Collection("teachers").InsertOne(ctx, *newTeacher)
		if err != nil {
			return nil, utils.ErrorHandler(err, "something went wrong")
		}

		objectId, ok := res.InsertedID.(bson.ObjectID)
		if ok {
			newTeacher.Id = objectId
		}

		pbTeacher := decodeEntity(&pb.Teacher{}, &newTeacher)
		addedTeachers = append(addedTeachers, pbTeacher)
	}
	return addedTeachers, nil
}

func GetTeachersFromDB(ctx context.Context, sortOptions bson.D, skip int64, limit int64, filter bson.M) ([]*pb.Teacher, int64, error) {
	client := Client()
	collection := client.Database("school-management").Collection("teachers")

	findOptions := options.Find()
	if len(sortOptions) > 0 {
		findOptions.SetSort(sortOptions)
	}
	findOptions.SetSkip(skip).SetLimit(limit)

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, utils.ErrorHandler(err, "something went wrong.")
	}
	defer cursor.Close(ctx)

	var teachers []*pb.Teacher
	for cursor.Next(ctx) {
		var teacher models.Teacher
		err = cursor.Decode(&teacher)
		if err != nil {
			return nil, 0, utils.ErrorHandler(err, "something went wrong.")
		}
		pbTeacher := &pb.Teacher{}
		result := decodeEntity(pbTeacher, &teacher)
		teachers = append(teachers, result)
	}

	err = cursor.Err()
	if err != nil {
		return nil, 0, utils.ErrorHandler(err, "something went wrong.")
	}

	totalCount, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, utils.ErrorHandler(err, "something went wrong.")
	}
	return teachers, totalCount, nil
}
