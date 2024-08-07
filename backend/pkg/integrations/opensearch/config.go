package opensearch

type Config struct {
	Host  string `json:"host" bson:"host"`
	Index string `json:"index" bson:"index"`
	Port  string `json:"port" bson:"port"`
	Ssl   bool   `json:"ssl" bson:"ssl"`
}
