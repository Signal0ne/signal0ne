package github

type Config struct {
	ApiKey string `json:"apiKey" bson:"apiKey"`
	Url    string `json:"url" bson:"url"`
}
