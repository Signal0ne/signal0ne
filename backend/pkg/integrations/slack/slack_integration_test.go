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

	input := PostMessageInput{
		ParsableContextObject: string(mockedEnrichedAlertStringified),
		SlackChannel:          "incidents",
	}

	//Create integration object
	integration := SlackIntegration{
		Inventory: SlackIntegrationInventory{
			AlertTitle: "Test Alert",
		},
		Config: Config{
			Url:         "http://localhost:8091",
			WorkspaceID: "workspace_123",
		},
	}

	_, err := postMessage(input, integration)
	assert.NoError(t, err)
}

func Test_CreateChannel(t *testing.T) {
	input := CreateChannelInput{
		ChannelName: "new_channel",
		IsPrivate:   "false",
	}

	integration := SlackIntegration{
		Config: Config{
			Url:         "http://localhost:8091",
			WorkspaceID: "workspace_123",
		},
	}

	output, err := createChannel(input, integration)
	assert.NoError(t, err)

	// Type assert output[0] to map[string]any
	if len(output) > 0 {
		result, ok := output[0].(map[string]any)
		assert.True(t, ok, "Expected output[0] to be of type map[string]any")
		assert.Equal(t, "success", result["status"])
	}
}

func Test_AddUsersToTheChannel(t *testing.T) {
	input := AddUsersToTheChannelInput{
		ChannelName: "new_channel",
		UserHandles: "vidhumathur2002@gmail.com",
	}

	integration := SlackIntegration{
		Config: Config{
			Url:         "http://localhost:8091", // Ensure your service is running on this port
			WorkspaceID: "workspace_123",
		},
	}

	output, err := addUsersToTheChannel(input, integration)
	assert.NoError(t, err)
	if len(output) > 0 {
		result, ok := output[0].(map[string]any)
		assert.True(t, ok, "Expected output[0] to be of type map[string]any")
		assert.Equal(t, "success", result["status"])
	}
}
