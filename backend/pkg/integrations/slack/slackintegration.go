package slack

type SlackIntegartion struct {
	Webhook             string `json:"webhook"`
	ConditionExpression string `json:"conditionExpression"`
}
