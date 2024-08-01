package models

type Step struct {
	Name        string                 `json:"name" bson:"name"`
	Function    string                 `json:"function" bson:"function"`
	Integration string                 `json:"integration" bson:"integration"`
	Input       map[string]interface{} `json:"input" bson:"input"`
	Output      StepOutput             `json:"output" bson:"output"`
	Condition   string                 `json:"condition" bson:"condition"`
}

type StepOutput struct {
	Data map[string]interface{} `json:"-" bson:"-"`
}

func (s *Step) ParseCondition() {

}
