package tools

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"signal0ne/cmd/config"
	"signal0ne/internal/errors" //only internal import allowed
	"signal0ne/internal/models" //only internal import allowed
	"strconv"
	"strings"
	"text/template"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GenerateWebhookSalt() (string, error) {
	randomBytes := make([]byte, 9)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	salt := base64.URLEncoding.EncodeToString(randomBytes)

	return salt, nil
}

func TraverseOutput(
	payload any,
	desiredKey string,
	mapping string) any {

	mappings := strings.Split(mapping, ".")
	currentMapping := mappings[0]

	switch v := payload.(type) {
	case map[string]any:
		for key, value := range v {
			if key == currentMapping || key == mapping {
				_, isMap := value.(map[string]any)
				_, isSlice := value.([]map[string]any)
				if len(mappings) <= 1 || (!isMap && !isSlice) {
					return value
				} else {
					mapping = strings.Join(mappings[1:], ".")
					return TraverseOutput(value, desiredKey, mapping)
				}
			}
		}
		return nil
	case []map[string]any:
		for _, elem := range v {
			for key, value := range elem {
				if key == currentMapping || key == mapping {
					_, isMap := value.(map[string]any)
					_, isSlice := value.([]map[string]any)
					if len(mappings) <= 1 || (!isMap && !isSlice) {
						return value
					} else {
						mapping = strings.Join(mappings[1:], ".")
						return TraverseOutput(value, desiredKey, mapping)
					}
				}
			}
		}
		return nil
	default:
		return nil
	}
}

func WebhookTriggerExec(payload map[string]any, workflow *models.Workflow) (map[string]any, error) {
	var desiredPropertiesWithValues = map[string]any{}

	for key, mapping := range workflow.Trigger.WebhookTrigger.Webhook.Output {
		desiredPropertiesWithValues[key] = TraverseOutput(payload, key, mapping)
	}

	alertWithTriggerProperties := models.EnrichedAlert{
		TriggerProperties: desiredPropertiesWithValues,
	}

	if !EvaluateCondition(workflow.Trigger.WebhookTrigger.Webhook.Condition,
		alertWithTriggerProperties) {
		return desiredPropertiesWithValues, errors.ErrConditionNotSatisfied
	}

	return desiredPropertiesWithValues, nil
}

func RecordExecution(
	ctx context.Context,
	newExecution models.StepExecution,
	workflowsCollection *mongo.Collection,
	filter bson.M) error {

	var workflow models.Workflow
	cfg := config.GetInstance()

	workflowOutput := workflowsCollection.FindOne(ctx, filter)
	workflowOutput.Decode(&workflow)

	executionsHistoryBacklog := len(workflow.Executions)
	excessiveBacklog := int64(executionsHistoryBacklog) - cfg.WorkflowsExecutionsHistoryLimit

	if excessiveBacklog > 0 {

		for i := 0; i < int(excessiveBacklog); i++ {
			_, err := workflowsCollection.UpdateOne(ctx, filter, bson.M{
				"$pop": bson.M{
					"executions": -1,
				},
			})
			if err != nil {
				return err
			}
		}
	}

	_, err := workflowsCollection.UpdateOne(ctx, filter, bson.M{
		"$push": bson.M{
			"executions": newExecution,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func ExecutionResultWrapper(intermediateResults []any,
	output map[string]string,
	outputTags []string) []map[string]any {
	var results []map[string]any

	for _, result := range intermediateResults {
		var traverseResult = map[string]any{}
		for key, mapping := range output {
			mappings := strings.Split(mapping, ",")
			for _, localMapping := range mappings {
				localMapping = strings.TrimSpace(localMapping)
				traverseResult[key] = TraverseOutput(result, key, localMapping)
				if traverseResult[key] != nil {
					break
				}
			}
		}
		traverseResult["tags"] = outputTags
		results = append(results, traverseResult)
	}

	return results
}

func EvaluateCondition(conditionExpression string, alert models.EnrichedAlert) bool {
	var satisfied = true
	buf := new(bytes.Buffer)

	if conditionExpression == "" {
		return satisfied
	}

	parsedTemplate, err := template.New("").Parse(conditionExpression)
	if err != nil {
		return satisfied
	}

	err = parsedTemplate.Execute(buf, alert)
	if err != nil {
		return satisfied
	}

	satisfied, err = strconv.ParseBool(buf.String())
	if err != nil {
		return satisfied
	}

	return satisfied
}

func MapAlertState(payload map[string]any, stateKey string, triggerStateMapping map[string]string) (models.AlertStatus, error) {
	stateValue, exists := payload[stateKey].(string)
	if !exists {
		return "", fmt.Errorf("cannot find state key in alert payload")
	}

	mappedStateValue, exists := triggerStateMapping[stateValue]
	if !exists {
		return "", fmt.Errorf("cannot find mapping for state value %s", stateValue)
	}

	return models.AlertStatus(mappedStateValue), nil
}

func GetStartTime(payload map[string]any, startTimeKey string) (time.Time, error) {
	var parsedTime time.Time
	startTime, exists := payload[startTimeKey].(string)
	if !exists {
		return parsedTime, fmt.Errorf("cannot find start time key in alert payload")
	}

	parsedTime, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return parsedTime, fmt.Errorf("invalid time format: %v", err)
	}

	return parsedTime, nil
}
