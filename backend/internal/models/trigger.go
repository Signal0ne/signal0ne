package models

type Trigger struct {
	Data map[string]interface{} `json:"-"`
}

type WebhookTrigger struct {
	Output map[string]interface{} `json:"-"`
}

type SchedulerTrigger struct {
	Interval string `json:"interval"`
}
