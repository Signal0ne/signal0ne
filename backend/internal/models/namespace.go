package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type NamespaceUserRef struct {
	Id       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Accepted bool               `json:"accepted" bson:"accepted"`
}

type Namespace struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Workflows []string           `json:"workflows" bson:"workflows"`
	Users     []NamespaceUserRef `json:"users" bson:"users"`
}
