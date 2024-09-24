package alertmanager

var TriggerStateMapping = map[string]string{
	"firing":   "open",
	"resolved": "inactive",
}

type Config struct {
	Url string `json:"url" bson:"url"`
}
