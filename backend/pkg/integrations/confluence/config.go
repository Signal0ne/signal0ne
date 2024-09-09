package confluence

type Config struct {
	Url    string `json:"url" bson:"url"`
	Email  string `json:"email" bson:"email"`
	ApiKey string `json:"apiKey" bson:"apiKey"`
}
