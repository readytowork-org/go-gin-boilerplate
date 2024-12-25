package utils

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/chai2010/webp"
)

// EncodeImage image
func EncodeImage(file image.Image, fileType string) (*bytes.Buffer, error) {
	var img bytes.Buffer
	var err error

	switch fileType {
	case "image/png":
		err = png.Encode(&img, file)
	case "image/jpeg", "image/jpg":
		err = jpeg.Encode(&img, file, nil)
	case "image/gif":
		err = gif.Encode(&img, file, nil)
	case "image/webp":
		err = webp.Encode(&img, file, nil)
	default:
		return nil, err
	}
	return &img, nil
}
