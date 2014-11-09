//
// resize.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Nov  8 2014.
//

package rw

import (
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

const (
	jpgOutput = 1
	pngOutput = 2
	gifOutput = 3
)

type Resize struct {
	Height        uint                         `json:"height"`
	Width         uint                         `json:"width"`
	Interpolation string                       `json:"interpolation,omitempty"`
	Output        string                       `json:"output"`
	Name          string                       `json:"-"`
	interpolation resize.InterpolationFunction `json:"-"`
	output        uint                         `json:"-"`
}

func (re *Resize) GetName() string {
	return re.Name
}

func (re *Resize) Init() error {
	if re.Height == 0 && re.Width == 0 {
		return RwError(re, "height and width cannot be both equal to 0")
	}
	switch re.Interpolation {
	case "NearestNeighbor":
		re.interpolation = resize.NearestNeighbor
	case "Bilinear":
		re.interpolation = resize.Bilinear
	case "Bicubic":
		re.interpolation = resize.Bicubic
	case "MitchellNetravali":
		re.interpolation = resize.MitchellNetravali
	case "Lanczos2":
		re.interpolation = resize.Lanczos2
	case "", "Lanczos3":
		re.interpolation = resize.Lanczos3
	default:
		return RwError(re, fmt.Sprintf("unknow interpolation function '%s'", re.Interpolation))
	}

	switch re.Output {
	case "", "jpg":
		re.output = jpgOutput
	case "png":
		re.output = pngOutput
	case "gif":
		re.output = gifOutput
	default:
		return RwError(re, fmt.Sprintf("unsuported output format '%s'", re.Output))
	}
	return nil
}

func (re *Resize) Encode(r io.Reader, w io.Writer, d *Data) error {
	i, _, err := image.Decode(r)
	if err != nil {
		return errors.New(fmt.Sprintf("%s: %s", re.GetName, err))
	}
	newImage := resize.Resize(re.Width, re.Height, i, re.interpolation)
	switch re.output {
	case jpgOutput:
		err = jpeg.Encode(w, newImage, nil)
	case pngOutput:
		err = png.Encode(w, newImage)
	case gifOutput:
		err = gif.Encode(w, newImage, nil)
	}
	if err != nil {
		return RwError(re, err.Error())
	}
	return nil
}
