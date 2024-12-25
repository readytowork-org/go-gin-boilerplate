package utils

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"mime/multipart"

	"github.com/chai2010/webp"
)

// DecodeImage image
func DecodeImage(file multipart.File, fileType string) (image.Image, error) {
	var img image.Image
	var err error

	// Seek back to beginning of file for CreateThumbnail
	if _, err := file.Seek(0, 0); err != nil {
		return nil, err
	}

	switch fileType {
	case "image/png":
		img, err = png.Decode(file)
	case "image/jpeg", "image/jpg":
		img, err = jpeg.Decode(file)
	case "image/gif":
		img, err = gif.Decode(file)
	case "image/webp":
		img, err = webp.Decode(file)
	default:
		return nil, err
	}
	return img, nil
}
