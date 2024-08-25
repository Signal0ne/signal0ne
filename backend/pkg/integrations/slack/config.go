package slack

type Config struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	WorkspaceID string `json:"workspaceId"`
}
