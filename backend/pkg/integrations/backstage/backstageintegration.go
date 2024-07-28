package backstage

type BackstageIntegartion struct {
	Host                string `json:"host"`
	Port                string `json:"port"`
	ApiKey              string `json:"apiKey"`
	ConditionExpression string `json:"conditionExpression"`
}

type GetOwnershipStepConfig struct {
	Filter string
}

type GetOwnershipStepResult struct {
	Ownership string
}

type CustomQueryStepConfig struct {
	Query string
}

type CustomQueryStepResult struct {
	Data interface{}
}

func (b *BackstageIntegartion) GetOwnership(config GetOwnershipStepConfig) (result GetOwnershipStepResult) {
	return result
}

func (b *BackstageIntegartion) CustomQuery(config CustomQueryStepConfig) (result GetOwnershipStepResult) {
	return result
}
