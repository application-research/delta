package models

type StatsCheckResponse struct {
	Content struct {
		ID      int64  `json:"id"`
		Name    string `json:"name"`
		Cid     string `json:"cid,omitempty"`
		Status  string `json:"status"`
		Message string `json:"message,omitempty"`
	} `json:"content"`
}
