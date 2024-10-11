package github

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/internal/utils"
	"signal0ne/pkg/integrations/helpers"
	"strings"
)

type GithubIntegration struct {
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:"config" bson:"config"`
}

var functions = map[string]models.WorkflowFunctionDefinition{
	"get_content": models.WorkflowFunctionDefinition{
		Function:   getContent,
		Input:      GetFileContentInput{},
		OutputTags: []string{"metadata", "logs"},
	},
	"get_commit_diff": models.WorkflowFunctionDefinition{
		Function:   getCommitDiff,
		Input:      GetCommitDiff{},
		OutputTags: []string{"metadata", "code"},
	},
	"inspect_github_actions": models.WorkflowFunctionDefinition{
		Function:   inspectGithubActions,
		Input:      InspectGithubActionsInput{},
		OutputTags: []string{"metadata"},
	},
}

func (integration GithubIntegration) Execute(
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

	results = tools.ExecutionResultWrapper(intermediateResults, output, function.OutputTags)

	return results, nil
}

func (integration GithubIntegration) Initialize() map[string]string {
	// Implement your config initialization here
	return nil
}

func (integration GithubIntegration) Validate() error {
	// Implement your config validation here
	return nil
}

func (integration GithubIntegration) ValidateStep(
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

type GetCommitDiff struct {
	Commit   string `json:"commit" bson:"commit"`
	RepoName string `json:"repo_name" bson:"repo_name"`
}

type GetFileContentInput struct {
	ContentUrl string `json:"content_url" bson:"content_url"`
	Path       string `json:"path" bson:"path"`
	Type       string `json:"type" bson:"type"`
}

type InspectGithubActionsInput struct {
	ActionName string `json:"action_name" bson:"action_name"`
	BranchName string `json:"branch_name" bson:"branch_name"`
	RepoName   string `json:"repo_name" bson:"repo_name"`
}

func getCommitDiff(input any, integration any) ([]any, error) {
	var parsedInput GetCommitDiff
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "get_commit_diff")
	if err != nil {
		return output, err
	}

	assertedIntegration := integration.(GithubIntegration)

	commitId := parsedInput.Commit

	assertedIntegration.Config.Url = strings.TrimSuffix(assertedIntegration.Config.Url, "/")
	decomposedUri := strings.Split(assertedIntegration.Config.Url, "/")
	org := decomposedUri[len(decomposedUri)-1]

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", org, parsedInput.RepoName, commitId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return output, err
	}

	req.Header.Set("Accept", "application/vnd.github.diff")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Authorization",
		fmt.Sprintf("Bearer %s", assertedIntegration.Config.ApiKey),
	)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return output, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return output, fmt.Errorf("failed to get diff: %s", resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return output, err
	}

	parsedResponseBody := string(responseBody)

	output = append(output, map[string]any{
		"diff": parsedResponseBody,
	})

	return output, nil
}

func getContent(input any, integration any) ([]any, error) {
	var parsedInput GetFileContentInput
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "get_content")
	if err != nil {
		return output, err
	}

	assertedIntegration := integration.(GithubIntegration)

	var url string
	if parsedInput.Path != "" {
		url = fmt.Sprintf("%s/%s", parsedInput.ContentUrl, parsedInput.Path)
	} else {
		url = parsedInput.ContentUrl
	}

	var content string
	switch parsedInput.Type {
	case "file":
		content, err = getAndDecodeFileContent(url, assertedIntegration.ApiKey)
		if err != nil {
			return output, err
		}
	case "logs":
		content, err = getLogsContent(url, assertedIntegration.ApiKey)
		if err != nil {
			return output, err
		}
	}

	output = append(output, map[string]any{
		"content": content,
	})

	return output, nil
}

func inspectGithubActions(input any, integration any) ([]any, error) {
	var parsedInput InspectGithubActionsInput
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "inspect_github_actions")
	if err != nil {
		return output, err
	}

	assertedIntegration := integration.(GithubIntegration)

	assertedIntegration.Config.Url = strings.TrimSuffix(assertedIntegration.Config.Url, "/")
	decomposedUri := strings.Split(assertedIntegration.Config.Url, "/")
	org := decomposedUri[len(decomposedUri)-1]

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs", org, parsedInput.RepoName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return output, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Authorization",
		fmt.Sprintf("Bearer %s", assertedIntegration.Config.ApiKey),
	)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return output, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return output, fmt.Errorf("failed to get GitHub actions: %s", resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return output, err
	}

	var actionsResponse map[string]any
	err = json.Unmarshal(responseBody, &actionsResponse)
	if err != nil {
		return output, err
	}

	actionRuns, ok := actionsResponse["workflow_runs"].([]any)
	if !ok {
		return output, fmt.Errorf("cannot parse GitHub actions")
	}

	latestAction, ok := actionRuns[0].(map[string]any)
	if !ok {
		return output, fmt.Errorf("cannot parse GitHub actions")
	}

	output = append(output, map[string]any{
		"status":     latestAction["conclusion"],
		"commit":     latestAction["head_sha"],
		"created_at": latestAction["created_at"],
		"url":        latestAction["html_url"],
	})

	return output, nil
}

func getAndDecodeFileContent(url string, token string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Authorization",
		fmt.Sprintf("Bearer %s", token),
	)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get file content: %s", resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	parsedResponseBody := make(map[string]any)
	err = json.Unmarshal(responseBody, &parsedResponseBody)
	if err != nil {
		return "", err
	}

	base64EncodedContent, exists := parsedResponseBody["content"].(string)
	if !exists {
		return "", fmt.Errorf("cannot parse file content")
	}

	content, err := base64.StdEncoding.DecodeString(base64EncodedContent)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func getLogsContent(url string, token string) (string, error) {
	var content []string
	var relevantLogs string

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return relevantLogs, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Authorization",
		fmt.Sprintf("Bearer %s", token),
	)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return relevantLogs, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return relevantLogs, fmt.Errorf("failed to get file content: %s", resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return relevantLogs, err
	}

	zipFilePath := fmt.Sprintf("./%s-logs.zip", utils.GenerateRandomString())
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return relevantLogs, err
	}
	defer zipFile.Close()

	_, err = zipFile.Write(responseBody)
	if err != nil {
		return relevantLogs, err
	}

	// Unzip the file
	zipReader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return relevantLogs, err
	}
	defer zipReader.Close()

	var buildFiles []string

	// Iterate through the files in the zip archive
	for _, file := range zipReader.File {
		// Check if the file is in the "build" directory
		if strings.HasPrefix(file.Name, "build/") && !file.FileInfo().IsDir() {
			buildFiles = append(buildFiles, file.Name)
		}
	}

	if len(buildFiles) == 0 {
		return relevantLogs, fmt.Errorf("no build files found in the logs")
	}

	for _, fileName := range buildFiles {
		file, err := zipReader.Open(fileName)
		if err != nil {
			return relevantLogs, err
		}
		defer file.Close()

		fileContent, err := io.ReadAll(file)
		if err != nil {
			return relevantLogs, err
		}

		content = append(content, string(fileContent))
	}

	relevantLogs = getRelevantLogs(content)

	return relevantLogs, nil
}

func getRelevantLogs(logs []string) string {
	//TBD:  Relevant keywords parametrized instead of hardcoded "error"
	var relevantLogs string

	for _, log := range logs {
		if strings.Contains(log, "error") {
			relevantLogs = log
			break
		}
	}

	return relevantLogs
}
