// Package avatar implements a 5x5 block avatar image creation.
//
// Example Usage
//
// 		import "github.com/gitzart/identicon/avatar"
//
// 		a := new(avatar.Avatar)
// 		a.Text = "sometext"
// 		m, _ := a.Create()
//
// The output is "image.Image" type.
//
// Further image processing can be done by using the helper functions
// provided in "github.com/gitzart/identicon" package.
package avatar

import (
	"crypto/sha1"
	"errors"
	"image"
	"image/color"
	"strings"
)

const nblock = 5

// DefaultBG is the default image background color.
var DefaultBG = color.NRGBA{0xed, 0xed, 0xed, 0xff}

// Avatar defines the properties to make an avatar image.
type Avatar struct {
	// Case insensitive text
	Text string

	// Non-zero positive image size
	Size int

	// The value is in percentage. Ranges from 0 to 10.
	// Off limits values will be clipped.
	Padding int

	// Customizable avatar colors
	BGColor color.NRGBA
	Color   color.NRGBA

	palette palette
}

type palette [nblock][nblock]bool

// Create performs the algorithm to make an avatar image.
func (a *Avatar) Create() (image.Image, error) {
	if err := a.init(); err != nil {
		return nil, err
	}

	r := image.Rect(0, 0, a.Size, a.Size)
	m := image.NewNRGBA(r)

	fillRect(m, r, a.BGColor) // set background

	// Set avatar
	avatarRect := alignCenter(a.Size, a.Padding)
	blockSize := avatarRect.Dx() / nblock
	vt := avatarRect.Min.Y
	for i := 0; i < nblock; i++ {
		hr := avatarRect.Min.X
		for j := 0; j < nblock; j++ {
			if a.palette[i][j] {
				b := image.Rect(hr, vt, hr+blockSize, vt+blockSize)
				fillRect(m, b, a.Color)
			}
			hr += blockSize
		}
		vt += blockSize
	}

	return m, nil
}

// init verifies the avatar properties and initalizes them.
func (a *Avatar) init() error {
	if a.Size < 1 {
		return errors.New("invalid Avatar.Size")
	}
	if a.Padding < 0 {
		a.Padding = 0
	} else if a.Padding > 10 {
		a.Padding = 10
	}
	a.Padding *= a.Size / 100

	s := sha1.Sum([]byte(strings.ToLower(a.Text)))
	sum := s[:]
	var zc color.NRGBA // zero value color

	if a.BGColor == zc {
		a.BGColor = DefaultBG
	}
	if a.Color == zc {
		a.Color = color.NRGBA{sum[0], sum[1], sum[2], 0xff}
	}
	a.palette = mixColor(sum[3:])
	return nil
}

// alignCenter finds the center position of the image.
func alignCenter(size, padding int) image.Rectangle {
	rem := (size - padding<<1) % nblock
	min := padding + rem>>1
	max := size - min
	return image.Rect(min, min, max, max)
}

// fillRect colors a certain area of the image.
func fillRect(m *image.NRGBA, r image.Rectangle, c color.NRGBA) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			m.SetNRGBA(x, y, c)
		}
	}
}

// mixColor makes a symmetrical color palette but with boolean values.
func mixColor(recipe []byte) (p palette) {
	l := len(p)
	mid := l>>1 + 1
	z := 0
	for i := 0; i < l; i++ {
		for j := 0; j < mid; j++ {
			p[i][j] = recipe[z]%2 == 0
			// Mirror the bool to the other end.
			p[i][l-(j+1)] = p[i][j]
			z++
		}
	}
	return p
}
