package alertmanager

var TriggerStateMapping = map[string]string{
	"firing":   "active",
	"resolved": "inactive",
}

type Config struct {
	Url string `json:"url" bson:"url"`
}
