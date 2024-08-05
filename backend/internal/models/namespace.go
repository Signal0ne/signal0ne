package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Namespace struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Workflows []string           `json:"workflows" bson:"workflows"`
	Users     []string           `json:"users" bson:"users"`
}
