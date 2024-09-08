package models

type Trigger struct {
	WebhookTrigger   `json:",inline"`
	SchedulerTrigger `json:",inline"`
}

type WebhookTrigger struct {
	Webhook   Webhook `json:"webhook"`
	Condition string  `json:"condition"`
}

type Webhook struct {
	Output map[string]string `json:"output"`
}

type SchedulerTrigger struct {
	Scheduler Scheduler `json:"scheduler"`
}

type Scheduler struct {
	IntervalInMinutes int64             `json:"interval"`
	Output            map[string]string `json:"output"`
}
