package tools

import (
	"encoding/json"
	"log"
	"time"
)

type LogEntry struct {
	Timestamp   string `json:"timestamp"`
	Integration string `json:"integration"`
	Function    string `json:"function"`
	Input       any    `json:"input"`
	Output      any    `json:"output"`
	Result      string `json:"result"`
	Message     string `json:"message"`
}

func LogStep(
	logger *log.Logger,
	function string,
	integration string,
	message string,
	result string,
	input any,
	output any,
) {

	if logger == nil {
		return
	}

	entry := LogEntry{
		Timestamp:   time.Now().Format(time.RFC3339),
		Integration: integration,
		Function:    function,
		Result:      result,
		Message:     message,
		Input:       input,
		Output:      output,
	}
	entryJSON, err := json.Marshal(entry)
	if err != nil {
		logger.Printf("Failed to marshal log entry: %v", err)
		return
	}
	logger.Println("--------------------------------------------------")
	logger.Println(string(entryJSON))
	logger.Println("--------------------------------------------------")
}
