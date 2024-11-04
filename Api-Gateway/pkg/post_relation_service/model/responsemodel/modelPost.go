package responsemodel_post

type AddPostResp struct {
	Caption string `json:"caption,omitempty"`
	UserId  string `json:"userid,omitempty"`

	Media string `json:"media,omitempty"`
}

type EditPostResp struct {
	Caption string `json:"caption"`
	PostId  string `json:"postid"`
	UserId  string `json:"userid" `
}

type CommonResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"after execution,omitempty"`
	Error      interface{} `json:"error,omitempty"`
}
