package requestmodel_post

import "mime/multipart"

type AddPostData struct {
	Caption string                  `form:"caption" validate:"lte=60"`
	Media   []*multipart.FileHeader `form:"media" validate:"required"`

	UserId string `validate:"required"`
}

type EditPost struct {
	Caption string `form:"caption" validate:"lte=60"`
	PostId  string `json:"postid" validate:"required,number"`

	UserId string `validate:"required"`
}
