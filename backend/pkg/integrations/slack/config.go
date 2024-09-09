package slack

type Config struct {
	Url         string `json:"url" bson:"url"`
	WorkspaceID string `json:"workspaceId" bson:"workspaceId"`
}
