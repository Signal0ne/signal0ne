package slack

import (
	"encoding/json"
	"signal0ne/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Test_PostMessage tests the function PostMessage
// You need to run signal0ne-slackapp on port "8091" before running this test
func Test_PostMessage(t *testing.T) {
	mockedEnrichedAlert := models.EnrichedAlert{
		Id: primitive.NewObjectID(),
	}

	mockedEnrichedAlertStringified, _ := json.Marshal(mockedEnrichedAlert)

	mockedPostMessagePayload := `{"text": "Hello, World!"}`

	input := PostMessageInput{
		ParsableContextObject: string(mockedEnrichedAlertStringified),
		PostMessagePayload:    mockedPostMessagePayload,
		SlackChannel:          "incidents",
	}

	//Create integration object
	integration := SlackIntegration{
		Inventory: SlackIntegrationInventory{
			AlertTitle: "Test Alert",
		},
		Config: Config{
			Url:         "localhost:8091",
			WorkspaceID: "workspace_123",
		},
	}

	_, err := postMessage(input, integration)
	assert.NoError(t, err)
}
