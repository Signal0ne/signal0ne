package opensearch

import (
	"signal0ne/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_GetLogOccurrences tests the function getLogOccurrences

// You need to run lab01
// In order to run this test without mocking
func Test_GetLogOccurrences(t *testing.T) {
	mockConn := utils.ConnectToSocket()
	defer mockConn.Close()
	mockedGetLogOccurrencesInput := map[string]string{
		"service":    "adservice",
		"query":      "{\"query\": {\"match_all\": {}}}",
		"compare_by": "resource.service.name",
	}

	//Create integration object
	inventory := NewOpenSearchIntegrationInventory(mockConn)
	integration := OpenSearchIntegration{
		Inventory: inventory,
		Config: Config{
			Index: "otel",
			Url:   "http://20.127.192.216:9200",
		},
	}

	output, err := getLogOccurrences(mockedGetLogOccurrencesInput, integration)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}
