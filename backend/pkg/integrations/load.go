package integrations

import (
	"embed"
	"fmt"
	"reflect"
	"signal0ne/pkg/integrations/alertmanager"
	"signal0ne/pkg/integrations/backstage"
	"signal0ne/pkg/integrations/confluence"
	"signal0ne/pkg/integrations/github"
	"signal0ne/pkg/integrations/jaeger"
	"signal0ne/pkg/integrations/openai"
	"signal0ne/pkg/integrations/opensearch"
	"signal0ne/pkg/integrations/pagerduty"
	"signal0ne/pkg/integrations/servicenow"
	"signal0ne/pkg/integrations/signal0ne"
	"signal0ne/pkg/integrations/slack"
	"sync"

	"gopkg.in/yaml.v3"
)

//go:embed .assets/metadata/*
var integrationsMDFS embed.FS
var once sync.Once
var installableIntegrationsLib = map[string]map[string]any{}
var globalErrorHandle error = nil

var InstallableIntegrationTypesLibrary = map[string]reflect.Type{
	"alertmanager": reflect.TypeOf(alertmanager.AlertmanagerIntegration{}),
	"backstage":    reflect.TypeOf(backstage.BackstageIntegration{}),
	"confluence":   reflect.TypeOf(confluence.ConfluenceIntegration{}),
	"github":       reflect.TypeOf(github.GithubIntegration{}),
	"jaeger":       reflect.TypeOf(jaeger.JaegerIntegration{}),
	"openai":       reflect.TypeOf(openai.OpenaiIntegration{}),
	"opensearch":   reflect.TypeOf(opensearch.OpenSearchIntegration{}),
	"pagerduty":    reflect.TypeOf(pagerduty.PagerdutyIntegration{}),
	"servicenow":   reflect.TypeOf(servicenow.ServicenowIntegration{}),
	"signal0ne":    reflect.TypeOf(signal0ne.Signal0neIntegration{}),
	"slack":        reflect.TypeOf(slack.SlackIntegration{}),
}

func GetInstallableIntegrationsLib() (map[string]map[string]any, error) {
	once.Do(func() {
		entries, err := integrationsMDFS.ReadDir(".assets/metadata")
		if err != nil {
			globalErrorHandle = fmt.Errorf("FATAL: while reading directory, error: %s", err)
			return
		}
		for _, integrationMDFSObject := range entries {
			var integration map[string]any
			if !integrationMDFSObject.IsDir() {
				rawBytes, err := integrationsMDFS.ReadFile(fmt.Sprintf(".assets/metadata/%s", integrationMDFSObject.Name()))
				if err != nil {
					fmt.Printf("Warning: failed to read integration metadata from: %s, error: %s", integrationMDFSObject.Name(), err)
				}
				err = yaml.Unmarshal(rawBytes, &integration)
				if err != nil {
					fmt.Printf("Warning: failed to read integration metadata from: %s, error: %s", integrationMDFSObject.Name(), err)
				}
				installableIntegrationsLib[integration["type"].(string)] = integration
			}
		}
	})

	return installableIntegrationsLib, globalErrorHandle
}
