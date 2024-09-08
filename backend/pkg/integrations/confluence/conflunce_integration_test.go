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
			APIKey: "ATATT3xFfGF02Xn91wGZy7z3kd8IQmP9jLYOb2OzgKZ1dRPrCSoWleBgK7iGtKw4rmid8Kc6lXbdjoIP4PotCA8OtLkv5qOtErrz3fcvXkON_8OBZSq2wxPkUEJoVvmPpmEVy-B7dR0XZZ0P5nbVnlZS88DF-OC5vVfVA6V0Zpa7eb7YPphNiIw=E062F0AA",
		},
	}

	output, err := search(mockedSearchInput, integration)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}
