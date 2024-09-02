package pagerduty

import (
	"encoding/json"
	"signal0ne/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_CreateIncident(t *testing.T) {
	mockedWorkflow := models.Workflow{}

	mockedEnrichedAlert := models.EnrichedAlert{
		Id:                primitive.NewObjectID(),
		TriggerProperties: map[string]any{},
		AdditionalContext: map[string]models.Outputs{},
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
		Inventory:   inventory,
		Integration: models.Integration{},
		Config:      Config{},
	}

	output, err := createIncident(input, integration)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}
