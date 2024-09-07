package confluence

type Config struct {
	Email  string `json:"email" bson:"email"`
	APIKey string `json:"apiKey" bson:"apiKey"`
}
