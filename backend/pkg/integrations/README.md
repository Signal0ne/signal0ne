# Signal0ne integrations - Contribution Guide

## Integartion metadata file

Create yaml file with following format `<integartion_name>.yaml` in the `pkg/integrations/.assets/metadata` directory 

File should contain fields `typeName`(unique), `imageUri` with uri to the integrated tool logo and `config` with demanded config field for integartin to work

## Integration root creation
1. Create directory with the name of your integration in `pkg/integrations/` path for example: `pkg/integrations/backstage`
2. Create main integration file within your created directory with name having following format `<integration_name>integration.go` for example `backstageintegration.go`
3. Create config file within your created directory with following name `config.go`

## Integration structure definiton

In the `<integartion_name>integration.go` create struct which implements `IIntegartion` interface and contains fields defined `models.Integartion` + custom integarion fields defined in `pkg/integrations/<integartion_name>/config.go`

Add integration type to types dictionary `InstallableIntegrationTypesLibrary` located in `pkg/integartions/load.go`


## Integartion functions

Implement functions provided by integartion. Insert them to the `functions map[string]func(T any, dryRun bool) (any, error)` with proper function name as key.