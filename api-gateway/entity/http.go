package entity

// HttpResp is used as the standard response of all endpoints
type HttpResp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}
