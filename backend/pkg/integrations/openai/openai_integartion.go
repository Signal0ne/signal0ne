package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations/helpers"
)

type OpenaiIntegration struct {
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
}

var functions = map[string]models.WorkflowFunctionDefinition{
	"propose_resolution_steps": models.WorkflowFunctionDefinition{
		Function: proposeResolutions,
		Input:    ProposeResolutionsInput{},
	},
	"summarize_context": models.WorkflowFunctionDefinition{
		Function: summarizeContext,
		Input:    SummarizeContextInput{},
	},
}

func (integration OpenaiIntegration) Execute(
	input any,
	output map[string]string,
	functionName string) ([]map[string]any, error) {

	var results []map[string]any

	function, ok := functions[functionName]
	if !ok {
		return results, fmt.Errorf("%s.%s: cannot find requested function", integration.Name, functionName)
	}

	intermediateResults, err := function.Function(input, integration)
	if err != nil {
		return results, fmt.Errorf("%s.%s:%v", integration.Name, functionName, err)
	}

	results = tools.ExecutionResultWrapper(intermediateResults, output)

	return results, nil
}

func (integration OpenaiIntegration) Validate() error {
	// Implement your config validation here
	return nil
}

func (integration OpenaiIntegration) ValidateStep(
	input any,
	functionName string,
) error {
	function, exists := functions[functionName]
	if !exists {
		return fmt.Errorf("cannot find selected function")
	}

	//Validate input parameters for the chosen function
	err := helpers.ValidateInputParameters(input, function.Input, functionName)
	if err != nil {
		return err
	}

	return nil
}

type SummarizeContextInput struct {
	Context string `json:"context"`
}

type ProposeResolutionsInput struct {
	AdditionalContext string `json:"additional_context"`
	Logs              string `json:"logs"`
}

func proposeResolutions(input any, integration any) ([]any, error) {
	var parsedInput ProposeResolutionsInput
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "propose_resolution_steps")
	if err != nil {
		return output, err
	}

	assertedIntegration := integration.(OpenaiIntegration)

	fmt.Printf("###\nExecuting OpenAi integration function...\n")
	model := assertedIntegration.Model
	apiKey := assertedIntegration.ApiKey
	prompt := fmt.Sprintf(`You are on-call engineer Based on the logs and additional context like documentation or runbooks propose resolutions.
		Response must contain up to 3 steps with resolutions.
		Logs: %s
		Additional Context %s`, parsedInput.Logs, parsedInput.AdditionalContext)

	resolutions, err := callOpenAiApi(prompt, model, apiKey)
	if err != nil {
		return output, err
	}

	output = append(output, map[string]any{
		"content": resolutions,
	})

	return output, nil
}

func summarizeContext(input any, integration any) ([]any, error) {
	var parsedInput SummarizeContextInput
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "summarize_context")
	if err != nil {
		return output, err
	}

	assertedIntegration := integration.(OpenaiIntegration)

	fmt.Printf("###\nExecuting OpenAi integration function...\n")
	model := assertedIntegration.Model
	apiKey := assertedIntegration.ApiKey
	prompt := fmt.Sprintf(`You are on-call engineer Based on the full context from the investigation summarize investigation context for other engineers.
		Response must contain one short paragraph of explanation of the probable root causes in full sentences. Try to correlate different context for an holistic overview.
		The full issue context: %s`, parsedInput.Context)

	summary, err := callOpenAiApi(prompt, model, apiKey)
	if err != nil {
		return output, err
	}

	output = append(output, map[string]any{
		"summary": summary,
	})

	return output, nil
}

func callOpenAiApi(prompt string, model string, apiKey string) (string, error) {
	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	type Request struct {
		Model       string    `json:"model"`
		Messages    []Message `json:"messages"`
		Temperature float64   `json:"temperature"`
	}

	type Response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	var response Response
	var apiURL = "https://api.openai.com/v1/chat/completions"

	client := &http.Client{}

	messages := []Message{
		{Role: "user", Content: prompt},
	}

	reqBody, err := json.Marshal(Request{
		Model:       model,
		Messages:    messages,
		Temperature: 0.1,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}
