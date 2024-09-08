package confluence

type Config struct {
	Url    string `json:"url" bson:"url"`
	Email  string `json:"email" bson:"email"`
	APIKey string `json:"apiKey" bson:"apiKey"`
}
