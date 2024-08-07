package opensearch

type Config struct {
	Host  string `json:"host" bson:"host"`
	Index string `json:"index" bson:"inex"`
	Port  string `json:"port" bson:"port"`
	Ssl   bool   `json:"ssl" bson:"ssl"`
}
