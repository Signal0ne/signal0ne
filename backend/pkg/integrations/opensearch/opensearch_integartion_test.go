package opensearch

import (
	"encoding/json"
	"log"
	"net"
	"signal0ne/cmd/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_GetLogOccurrences tests the function getLogOccurrences

// You need to run lab01 -  https://portal.azure.com/#@szymonst2808gmail.onmicrosoft.com/resource/subscriptions/fb775820-301c-4a7d-af99-83285b864825/resourceGroups/rg01-lab01/providers/Microsoft.Compute/virtualMachines/lab01/overview
// In order to run this test without mocking
func Test_GetLogOccurrences(t *testing.T) {
	cfg := config.GetInstance()
	var mockConn, err = net.DialTimeout("unix", cfg.IPCSocket, (15 * time.Second))
	if err != nil {
		log.Fatalf("Error connecting to the socket: %s", err)
	}

	mockedGetLogOccurrencesInput := GetLogOccurrencesInput{
		Query: `{
					"query": {
						"match_all": {}
					}
				  }`,
		CompareBy: "resource.service.name",
	}

	mockedGetLogOccurrencesInputStringified, _ := json.Marshal(mockedGetLogOccurrencesInput)

	input := mockedGetLogOccurrencesInputStringified

	//Create integration object
	integration := OpenSearchIntegration{
		Inventory: OpenSearchIntegrationInventory{
			PyInterface: mockConn,
		},
		Config: Config{
			Host:  "20.127.192.216",
			Index: "otel",
			Port:  "9200",
			Ssl:   false,
		},
	}

	output, err := getLogOccurrences(input, integration)
	assert.NoError(t, err)
	assert.NotNil(t, output)

}
