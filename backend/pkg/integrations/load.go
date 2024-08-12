package integrations

import (
	"embed"
	"fmt"
	"reflect"
	"signal0ne/pkg/integrations/backstage"
	"signal0ne/pkg/integrations/jaeger"
	"signal0ne/pkg/integrations/opensearch"
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
	"backstage":  reflect.TypeOf(backstage.BackstageIntegration{}),
	"jaeger":     reflect.TypeOf(jaeger.JaegerIntegration{}),
	"opensearch": reflect.TypeOf(opensearch.OpenSearchIntegration{}),
	"signal0ne":  reflect.TypeOf(signal0ne.Signal0neIntegration{}),
	"slack":      reflect.TypeOf(slack.SlackIntegration{}),
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
				installableIntegrationsLib[integration["typeName"].(string)] = integration
			}
		}
	})

	return installableIntegrationsLib, globalErrorHandle
}
