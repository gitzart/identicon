// Package identicon implements identicon creation.
//
// This specific file does not contain the identicon implementions.
// It only provides some helper functions to smooth the identicon creation.
package identicon

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"

	"github.com/gitzart/identicon/avatar"
)

// Identicon represents an avatar.
type Identicon interface {
	Create() (image.Image, error)
}

// NewAvatar returns a new pointer avatar.
func NewAvatar(text string, size, padding int) *avatar.Avatar {
	return &avatar.Avatar{
		Text:    text,
		Size:    size,
		Padding: padding,
	}
}

// Must panics on identicon creation error.
// Must is a convenience wrapper that panics when passed a non-nil error value
func Must(m image.Image, err error) image.Image {
	if err != nil {
		panic(err)
	}
	return m
}

// EncodeJPEG performs the JPEG formatting.
func EncodeJPEG(m image.Image, quality int) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, m, &jpeg.Options{Quality: quality})
	return buf.Bytes(), err
}

// EncodePNG performs the PNG formatting with the BestCompression setting.
func EncodePNG(m image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	err := encoder.Encode(buf, m)
	return buf.Bytes(), err
}

// EncodeBase64 performs base64 formatting for HTML base64 image tag.
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// SaveToFile writes the image to a file.
func SaveToFile(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0644)
}
