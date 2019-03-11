package api

type (
	AgentSink struct {
		Name          string `json:"name"`
		TokenFilePath string `json:"token_file_path"`
		DHType        string `json:"dh_type"`
		DHPath        string `json:"dh_path"`
	}
)
