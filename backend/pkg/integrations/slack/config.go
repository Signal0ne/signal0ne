package slack

type Config struct {
	WorkspaceID string `json:"workspaceId"`
	Host        string `json:"host"`
	Port        string `json:"port"`
}
