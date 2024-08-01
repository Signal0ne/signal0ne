# Signal0ne integrations - Contribution Guide

## Integration metadata file

Create yaml file with following format `<integration_name>.yaml` in the
`pkg/integrations/.assets/metadata` directory

File should contain fields `typeName`(unique), `imageUri` with uri to the
integrated tool logo and `config` with demanded config field for integration to
work

## Integration root creation

1. Create directory with the name of your integration in `pkg/integrations/`
   path for example: `pkg/integrations/backstage`
2. Create main integration file within your created directory with name having
   following format `<integration_name>_integration.go` for example
   `backstage_integration.go`
3. Create config file within your created directory with following name
   `config.go`

## Integration structure definition

In the `<integration_name>_integration.go` create struct which implements
`IIntegration` interface and contains fields defined `models.Integration` +
custom integration fields defined in
`pkg/integrations/<integration_name>/config.go`

Add integration type to types dictionary `InstallableIntegrationTypesLibrary`
located in `pkg/integrations/load.go`

## Integration functions

Implement functions provided by integration. Insert them to the
`functions map[string]func(T any, dryRun bool) (any, error)` with proper
function name as key.
