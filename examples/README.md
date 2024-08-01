# Examples reference

| Name | Desciribtion | Used Integrations |
|----------|----------|----------|
| ErrorRateByService.yaml | Workflow to enrich SLO error rate breach with logs and other necessary info to resolve issue much faster skipping manual analysis. | Slack, OpenAI, OpenSearch, Backstage, Jeager |


## Trigger definition high level reference
```yaml
trigger:
  webhook: # webhook | schedule - enrichment can happen on event and on timely basis
    output:
      # map your expected incoming values to general workflow fields
```

## Step definition high level reference
```yaml
- name: StepName #(required)
  integration: typeOfIntegartion #(required), type of integartion to use in the step
  function: integrationFunctionToUse #(required), function exposed by integration to use in the step
  input: #(required if enforced by function), input to the function
  output: #(optional), output from the function it is always returned by can be catched by user or not(see ./ErrorRateByService.yaml)
  condition: #(optional), condition on which step should run, step is always run if condition is empty
```