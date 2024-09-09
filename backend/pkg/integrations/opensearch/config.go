package opensearch

type Config struct {
	Url   string `json:"url" bson:"url"`
	Index string `json:"index" bson:"index"`
}
