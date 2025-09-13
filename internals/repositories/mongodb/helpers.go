package mongodb

import (
	"reflect"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func decodeEntity[T any, M any](entity *T, model *M) *T {

	modelVal := reflect.ValueOf(model).Elem()
	if modelVal.Kind() != reflect.Struct {
		panic("model must be a pointer to a struct")
	}

	pbVal := reflect.ValueOf(entity).Elem()

	for i := range modelVal.NumField() {
		modelField := modelVal.Field(i)
		modelFieldType := modelVal.Type().Field(i)

		pbField := pbVal.FieldByName(modelFieldType.Name)
		if !pbField.IsValid() || !pbField.CanSet() {
			continue
		}

		// Special case: convert ObjectID → string
		if modelField.Type() == reflect.TypeOf(bson.ObjectID{}) {
			oid := modelField.Interface().(bson.ObjectID)
			pbField.SetString(oid.Hex())
		} else {
			// Otherwise copy directly
			pbField.Set(modelField)
		}
	}
	return entity
}
