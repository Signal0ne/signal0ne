package github

import (
	"signal0ne/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetFileContent(t *testing.T) {

	input := GetFileContentInput{
		ContentUrl: "https://api.github.com/repos/Signal0ne/local-log-farm/actions/runs/10753566051/logs",
		Type:       "logs",
	}

	integration := GithubIntegration{
		Integration: models.Integration{
			Name: "github",
			Type: "github",
		},
		Config: Config{
			ApiKey: "<your api key here>",
		},
	}

	output, err := getContent(input, integration)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}
