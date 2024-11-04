package responsemodel_auth

type CommonResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"after execution,omitempty"`
	Error      interface{} `json:"error,omitempty"`
}
