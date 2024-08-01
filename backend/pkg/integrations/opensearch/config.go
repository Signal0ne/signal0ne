package opensearch

type Config struct {
	Host string `json:"host" bson:"host"`
	Port string `json:"port" bson:"port"`
	Ssl  bool   `json:"ssl" bson:"ssl"`
}
