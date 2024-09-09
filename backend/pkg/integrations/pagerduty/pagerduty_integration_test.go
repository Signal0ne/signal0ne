package pagerduty

import (
	"encoding/json"
	"signal0ne/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/yaml.v3"
)

func Test_CreateIncident(t *testing.T) {
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
- name: Get relevant logs
  integration: opensearch
  function: get_log_occurrences
  input:
    service: "{{.TriggerProperties.service}}"
    query: '{"query": {"bool": {"must": [{"match": {"resource.service.name": "{{.TriggerProperties.service}}"}},{"range": {"@timestamp": {"gte": "{{date .TriggerProperties.timestamp "-15m" "rfc"}}","lte": "{{date .TriggerProperties.timestamp "+15m" "rfc"}}"}}}],"must_not": [{"match": {"severity.text": "INFO"}},{"match": {"severity.text": "Information"}}]}}}'
    compare_by: "body, resource.host.name"
  output:
    output_source: "output_source"
    count: "count"
    body: "body"
    hostname: "resource.host.name"
- name: Create incident
  integration: pagerduty
  function: create_incident
  input:
    type: "incident"
    title: "Error rate is high"
    service_name: "{{.TriggerProperties.service}}"
    parsable_context_object: "{{root}}"
`

	var temporaryWorkflow map[string]any
	err := yaml.Unmarshal([]byte(testWorkflowYaML), &temporaryWorkflow)
	if err != nil {
		t.Error(err)
	}

	jsonWorkflow, _ := json.Marshal(temporaryWorkflow)

	json.Unmarshal(jsonWorkflow, &mockedWorkflow)

	mockedEnrichedAlert := models.EnrichedAlert{
		Id: primitive.NewObjectID(),
		TriggerProperties: map[string]any{
			"service":   "Default Service",
			"timestamp": "2021-09-01T00:00:00Z",
			"category":  "core",
		},
		AdditionalContext: map[string]models.Outputs{
			"opensearch_get_log_occurrences": models.Outputs{
				Output: []any{
					map[string]any{
						"occurrences": 1,
						"body":        "Test body",
						"hostname":    "Test hostname",
					},
				},
			},
		},
	}

	mockedEnrichedAlertStringified, _ := json.Marshal(mockedEnrichedAlert)

	input := CreateIncidentInput{
		Type:                  "incident",
		Title:                 "Test incident",
		ServiceName:           "Default Service",
		ParsableContextObject: string(mockedEnrichedAlertStringified),
	}

	inventory := NewPagerdutyIntegrationInventory(&mockedWorkflow)

	integration := PagerdutyIntegration{
		Inventory: inventory,
		Integration: models.Integration{
			Name: "pagerduty",
			Type: "pagerduty",
		},
		Config: Config{
			Url:    "https://api.pagerduty.com",
			ApiKey: "<paste_api_key_here>",
		},
	}

	_, err = createIncident(input, integration)
	assert.NoError(t, err)
}
