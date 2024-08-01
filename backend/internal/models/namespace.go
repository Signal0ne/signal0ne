package models

type Namespace struct {
	Name      string   `json:"name" bson:"name"`
	Workflows []string `json:"workflows" bson:"workflows"`
	Users     []string `json:"users" bson:"users"`
}
