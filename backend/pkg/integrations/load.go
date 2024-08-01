package integrations

import (
	"embed"
	"fmt"
	"reflect"
	"signal0ne/pkg/integrations/backstage"
	"signal0ne/pkg/integrations/slack"
	"sync"

	"gopkg.in/yaml.v3"
)

//go:embed .assets/metadata/*
var integrationsMDFS embed.FS
var once sync.Once
var installableIntegrationsLib map[string]map[string]string
var globalErrorHandle error = nil

var InstallableIntegrationTypesLibrary = map[string]reflect.Type{
	"backstage": reflect.TypeOf(backstage.BackstageIntegration{}),
	"slack":     reflect.TypeOf(slack.SlackIntegartion{}),
}

func GetInstallableIntegrationsLib() (map[string]map[string]string, error) {
	once.Do(func() {
		entries, err := integrationsMDFS.ReadDir("metadata")
		if err != nil {
			globalErrorHandle = fmt.Errorf("FATAL: %s", err)
			return
		}
		for _, integrationMDFSObject := range entries {
			var integration map[string]string
			if !integrationMDFSObject.IsDir() {
				rawBytes, err := integrationsMDFS.ReadFile(integrationMDFSObject.Name())
				if err != nil {
					fmt.Printf("Warning: failed to read integartion metadata from: %s", integrationMDFSObject.Name())
				}
				err = yaml.Unmarshal(rawBytes, integration)
				if err != nil {
					fmt.Printf("Warning: failed to read integartion metadata from: %s", integrationMDFSObject.Name())
				}
				installableIntegrationsLib[integration["typeName"]] = integration
			}
		}
	})

	return installableIntegrationsLib, globalErrorHandle
}
