package mediafileformatchecker

import (
	"errors"
	"mime/multipart"
	"net/http"
)

func Mediafileformatchecker(file multipart.File) (*string, error) {
	allowedTypes := map[string]struct{}{
		"image/jpeg":      {},
		"image/png":       {},
		"image/gif":       {},
		"video/mp4":       {},
		"video/quicktime": {},
	}

	// Read the first 512 bytes to determine the content type
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return nil, err
	}

	// Reset the file position after reading
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	// Get the content type based on the file content
	contentType := http.DetectContentType(buffer)

	// Check if the content type is allowed
	if _, ok := allowedTypes[contentType]; !ok {
		return nil, errors.New("unsupported file type,should be a jpeg,png,gif,mp4 or quicktime")
	}

	return &contentType, nil
}

func ProfileImageFileFormatChecker(file multipart.File) (*string, error) {
	allowedTypes := map[string]struct{}{
		"image/jpeg": {},
		"image/png":  {},
	}

	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return nil, err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	contentType := http.DetectContentType(buffer)

	if _, ok := allowedTypes[contentType]; !ok {
		return nil, errors.New("unsupported file type, should be jpeg or png")
	}

	return &contentType, nil
}
