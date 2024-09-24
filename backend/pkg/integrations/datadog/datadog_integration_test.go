package datadog

import (
	"context"
	"encoding/json"
	"fmt"
	"signal0ne/cmd/config"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/internal/utils"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

// Test_GetRelevantLogs tests the function getRelevantLogs

// You need to run lab01
// In order to run this test without mocking
func Test_GetRelevantLogs(t *testing.T) {
	var mockedWorkflow models.Workflow
	testWorkflowYaML := `
name: ErrorRateByService - Incident
lookback: 15m
trigger:
  webhook:
    output:
      service: "job"
      span: "span_name"
      timestamp: "startsAt"
      category: "category"
steps:
- displayName: Get relevant logs
  name: get_relevant_logs_dd
  integration: datadog
  function: get_relevant_logs
  input:
    service: "{{.TriggerProperties.service}}"
    query: '*:message 3'
    compare_by: "message, host"
  output:
    output_source: "output_source"
    count: "count"
    body: "body"
`

	var temporaryWorkflow map[string]any
	err := yaml.Unmarshal([]byte(testWorkflowYaML), &temporaryWorkflow)
	if err != nil {
		t.Error(err)
	}

	jsonWorkflow, _ := json.Marshal(temporaryWorkflow)

	json.Unmarshal(jsonWorkflow, &mockedWorkflow)

	cfg := config.GetInstance()
	if cfg == nil {
		panic("CRITICAL: unable to load config")
	}

	ctx := context.Background()

	mockConn := utils.ConnectToSocket()
	defer mockConn.Close()

	mongoConn, err := tools.InitMongoClient(ctx, cfg.MongoUri)
	if err != nil {
		panic(
			fmt.Sprintf("Failed to establish connection to %s, error: %s",
				strings.Split(cfg.MongoUri, "/")[2],
				err),
		)
	}
	defer mongoConn.Disconnect(ctx)

	alertsCollection := mongoConn.Database("signalone").Collection("alerts")

	mockedGetRelevantLogsInput := map[string]any{
		"service":    "adservice",
		"query":      "*:message 3",
		"compare_by": "attributes.message, attributes.host",
		"limit":      10,
	}

	//Create integration object
	inventory := NewDataDogIntegrationInventory(
		mockConn,
		alertsCollection,
		&mockedWorkflow)

	integration := DataDogIntegration{
		Inventory: inventory,
		Config: Config{
			Url:            "https://api.us5.datadoghq.com",
			ApiKey:         "<API_KEY>",
			ApplicationKey: "<APPLICATION_KEY>",
		},
	}

	output, err := getRelevantLogs(mockedGetRelevantLogsInput, integration)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}
