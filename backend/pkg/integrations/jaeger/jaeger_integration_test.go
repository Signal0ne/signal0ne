package jaeger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_GetDependencies tests the function getDependencies

// You need to run lab01
// In order to run this test without mocking
func Test_GetDependencies(t *testing.T) {

	mockedGetDependenciesInput := map[string]string{
		"service": "adservice",
	}

	inventory := NewJaegerIntegrationInventory(nil)
	integration := JaegerIntegration{
		Inventory: inventory,
		Config: Config{
			Url: "http://20.127.192.216:16686",
		},
	}

	output, err := getDependencies(mockedGetDependenciesInput, integration)

	fmt.Printf("output: %v\n", output)

	assert.NoError(t, err)
	assert.NotNil(t, output)

}
