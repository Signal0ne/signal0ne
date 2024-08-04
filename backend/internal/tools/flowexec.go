package tools

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"signal0ne/internal/models"
	"time"

	"github.com/gin-gonic/gin"
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

	switch v := payload.(type) {
	case map[string]any:
		for key, value := range v {
			if key == mapping {
				return value
			}
			TraverseOutput(value, desiredKey, mapping)
		}
	case []any:
		for _, value := range v {
			TraverseOutput(value, desiredKey, mapping)
		}
	default:
		return v
	}
	return nil
}

func WebhookTriggerExec(ctx *gin.Context, workflow *models.Workflow) (map[string]any, error) {
	var incomingTriggerPayload map[string]any
	var desiredPropertiesWithValues = map[string]any{}
	body, err := ctx.GetRawData()
	if err != nil || len(body) == 0 {
		return desiredPropertiesWithValues, fmt.Errorf("cannot get body %s", err)
	}

	err = json.Unmarshal(body, &incomingTriggerPayload)
	if err != nil {
		return desiredPropertiesWithValues, fmt.Errorf("cannot decode body %s", err)
	}

	for key, mapping := range workflow.Trigger.WebhookTrigger.Webhook.Output {
		desiredPropertiesWithValues[key] = TraverseOutput(incomingTriggerPayload, key, mapping)
	}

	return desiredPropertiesWithValues, nil
}

func RecordExecution(
	ctx context.Context,
	localErrorMessage string,
	workflowsCollection *mongo.Collection,
	filter bson.M) error {

	var status string
	var log string
	if localErrorMessage == "" {
		status = "Success"
		log = "Successfully executed"
	} else {
		status = "Failure"
		log = localErrorMessage
	}

	newExecution := models.Execution{
		Status:    status,
		Log:       log,
		Timestamp: time.Now().Unix(),
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
