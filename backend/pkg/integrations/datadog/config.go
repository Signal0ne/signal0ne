package datadog

var TriggerStateMapping = map[string]string{
	// $ALERT_TRANSITION datadog build-in variable, https://docs.datadoghq.com/integrations/webhooks/#usage
	"Triggered":    "open",
	"Re-Triggered": "open",
	"Renotify":     "open",
	"Warn":         "open",
	"Re-Warn":      "open",
	"Recovered":    "inactive",
}

type Config struct {
	ApiKey string `json:"apiKey "bson:"apiKey"`
	Url    string `json:"url" bson:"url"`
}
