package models

type Teacher struct {
	Id        string `protobuf:"id,omitempty" bson:"_id,omitempty"`
	FirstName string `protobuf:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  string `protobuf:"last_name,omitempty" bson:"_last_name,omitempty"`
	Email     string `protobuf:"email,omitempty" bson:"_email,omitempty"`
	Subject   string `protobuf:"subject,omitempty" bson:"_subject,omitempty"`
	Class     string `protobuf:"class,omitempty" bson:"_class,omitempty"`
}
