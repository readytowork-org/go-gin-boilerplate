package utils

import (
	"bytes"
	"mime/multipart"

	"github.com/nfnt/resize"
)

// CreateThumbnail image
func CreateThumbnail(file multipart.File, fileType string, width, height uint) (*bytes.Buffer, error) {
	thumbnailImg, err := DecodeImage(file, fileType)
	if err != nil {
		return nil, err
	}
	thumbnail := resize.Resize(width, height, thumbnailImg, resize.Lanczos3)
	thumbnailFile, err := EncodeImage(thumbnail, fileType)
	if err != nil {
		return nil, err
	}

	return thumbnailFile, nil
}
