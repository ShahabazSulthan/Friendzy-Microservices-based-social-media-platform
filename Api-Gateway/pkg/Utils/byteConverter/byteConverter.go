package byteconverter

import (
	"io"
	"mime/multipart"
)

func MultipartFileheaderToBytes(file *multipart.File) ([]byte, error) {
	content, err := io.ReadAll(*file)
	if err != nil {
		return nil, err
	}

	return content, nil
}
