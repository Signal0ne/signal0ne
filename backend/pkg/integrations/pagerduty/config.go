package pagerduty

type Config struct {
	Url    string `json:"url" bson:"url"`
	ApiKey string `json:"apiKey bson:"apiKey"`
}
