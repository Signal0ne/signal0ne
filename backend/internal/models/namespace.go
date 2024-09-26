package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type NamespaceUserRef struct {
	Id       primitive.ObjectID `json:"id" bson:"_id"`
	Accepted bool               `json:"accepted" bson:"accepted"`
	Username string             `json:"username" bson:"username"`
}

type Namespace struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Users     []NamespaceUserRef `json:"users" bson:"users"`
	Workflows []string           `json:"workflows" bson:"workflows"`
}
