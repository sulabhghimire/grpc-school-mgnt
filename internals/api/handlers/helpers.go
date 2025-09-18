package handlers

import (
	"errors"
	"fmt"
	"grpc-school-mgnt/pkg/utils"
	pb "grpc-school-mgnt/proto/gen"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func buildFilter(reqfilter any, model any) (bson.M, error) {
	filter := bson.M{}

	if reqfilter == nil || reflect.ValueOf(reqfilter).IsNil() {
		return filter, nil
	}

	modelVal := reflect.ValueOf(model).Elem()
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
		fieldName := modelType.Field(i).Name

		if fieldVal.IsValid() && !fieldVal.IsZero() {
			bsonTag := strings.TrimSuffix(modelType.Field(i).Tag.Get("bson"), ",omitempty")
			if strings.TrimSpace(bsonTag) != "" {
				if bsonTag == "_id" {
					_id, err := bson.ObjectIDFromHex(filterVal.FieldByName(fieldName).Interface().(string))
					if err != nil {
						return nil, errors.New("invalid _id for teacher")
					}
					filter[bsonTag] = _id
				} else {
					filter[bsonTag] = fieldVal.Interface()

				}
			}
		}
	}

	return filter, nil
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

func mapPbToModel[P any, M any](pbStruct P, newModel func() *M) *M {
	model := newModel()
	pbVal := reflect.ValueOf(pbStruct).Elem()
	modelVal := reflect.ValueOf(&model).Elem()

	for i := range pbVal.NumField() {
		pbField := pbVal.Field(i)
		fieldName := pbVal.Type().Field(i).Name

		modelField := modelVal.FieldByName(fieldName)

		if !modelField.IsValid() || !modelField.CanSet() {
			continue
		}

		// If model field is ObjectID and pb field is string → convert
		if modelField.Type() == reflect.TypeOf(bson.ObjectID{}) && pbField.Kind() == reflect.String {
			if oid, err := bson.ObjectIDFromHex(pbField.String()); err == nil {
				modelField.Set(reflect.ValueOf(oid))
			}
		} else {
			// Direct assignment if types match
			if pbField.Type().AssignableTo(modelField.Type()) {
				modelField.Set(pbField)
			}
		}
	}

	return model
}

func stringToObjectId(val string) (*bson.ObjectID, error) {
	objectId, err := bson.ObjectIDFromHex(val)
	if err != nil {
		return nil, utils.ErrorHandler(err, fmt.Sprintf("please provide valid teacher id %d", val))
	}
	return &objectId, nil
}

// func objectIdToString(oid bson.ObjectID) string {
// 	return oid.Hex()
// }
