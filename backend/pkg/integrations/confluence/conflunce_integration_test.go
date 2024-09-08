package confluence

import (
	"signal0ne/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Search(t *testing.T) {
	mockConn := utils.ConnectToSocket()
	defer mockConn.Close()
	mockedSearchInput := map[string]string{
		"query":           "title~build",
		"similarity_case": "2024-09-07T18:30:19.6825159Z ##[error]Username and password required",
	}

	inventory := NewConfluenceIntegrationInventory(mockConn)
	integration := ConfluenceIntegration{
		Inventory: inventory,
		Config: Config{
			Url:    "https://signaloneai.atlassian.net",
			Email:  "contact@signaloneai.com",
			APIKey: "<your api key here>",
		},
	}

	_, err := search(mockedSearchInput, integration)
	assert.NoError(t, err)
}
