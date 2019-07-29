// Quick and dirty command line identicon generator.
package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"

	"github.com/gitzart/identicon"
)

var appVersion = "1.0.0"

var (
	help, version     bool
	size, padding     int
	bgColor, colorVar string
	filename, path    string
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	flag.Parse()

	// End with help or version
	if help {
		usage()
		os.Exit(0)
	}
	if version {
		fmt.Println("version " + appVersion)
		os.Exit(0)
	}

	text := flag.Arg(0)
	if text == "" {
		fmt.Println("text is required. See help (-h)")
		os.Exit(1)
	}

	// Parse colors
	bg, err := parseColor(bgColor)
	failCheck(err)
	c, err := parseColor(colorVar)
	failCheck(err)

	a := identicon.NewAvatar(text, size, padding)
	a.BGColor = bg
	a.Color = c

	// Change directory
	if path != "" {
		failCheck(os.Chdir(path))
	}

	m, err := a.Create()
	failCheck(err)

	data, err := identicon.EncodePNG(m)
	failCheck(err)

	if filename == "" {
		filename = strings.Replace(text, " ", "_", -1) + ".png"
	}

	failCheck(identicon.SaveToFile(filename, data))
	fmt.Println("done")
}

func failCheck(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseColor(c string) (color.NRGBA, error) {
	var zc color.NRGBA // zero value color
	if c == "" {
		return zc, nil
	}
	if len(c) != 6 {
		return zc, fmt.Errorf("color must have 6 digits")
	}
	var ca []byte
	for i := 0; i < 6; i += 2 {
		v, err := strconv.ParseUint(c[i:i+2], 16, 8)
		if err != nil {
			return zc, err
		}
		ca = append(ca, byte(v))
	}
	return color.NRGBA{ca[0], ca[1], ca[2], 255}, nil
}

func usage() {
	s := `Usage:
  avatar [options] <text>

Options:
  -h	app usage
  -v	app version

  -s	non-zero positive image size
  -p	image padding (0-10)
  -b	image background color (ffffff)
  -c	avatar color (ffffff)
  -n	image file name
  -d	path to write the image to
`
	fmt.Fprintf(os.Stderr, s)
}

func init() {
	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&version, "v", false, "")
	flag.IntVar(&size, "s", 320, "")
	flag.IntVar(&padding, "p", 10, "")
	flag.StringVar(&bgColor, "b", "", "")
	flag.StringVar(&colorVar, "c", "", "")
	flag.StringVar(&filename, "n", "", "")
	flag.StringVar(&path, "d", "", "")

	flag.Usage = func() {
		fmt.Println("See help (-h)")
	}
}
