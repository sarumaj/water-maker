package handler

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/httpsOmkar/graphics-go/graphics"
	data "github.com/sarumaj/water-maker/pkg/data"
)

type File struct {
	FullPath, Dir, Base, Ext string
}

func (file *File) SetWatermark() {
	imgb, err := os.Open(file.FullPath)
	if err != nil {
		panic(err)
	}
	defer imgb.Close()

	var img image.Image
	switch file.Ext {
	case "png":
		img, err = png.Decode(imgb)
	case "jpeg", "jpg":
		img, err = jpeg.Decode(imgb)
	case "gif":
		img, err = gif.Decode(imgb)
	default:
		return
	}
	if err != nil {
		panic(err)
	}

	wmb, err := data.Fs.Open("images/watermark.png")
	if err != nil {
		panic(err)
	}
	defer wmb.Close()

	watermark, err := png.Decode(wmb)
	if err != nil {
		panic(err)
	}

	if s := img.Bounds().Size(); s.Y > s.X {
		dstImage := image.NewRGBA(image.Rect(0, 0, watermark.Bounds().Dy(), watermark.Bounds().Dx()))
		graphics.Rotate(dstImage, watermark, &graphics.RotateOptions{Angle: 3.0 * math.Pi / 2.0})
		watermark = dstImage
	}
	dstImage := image.NewRGBA(img.Bounds())
	graphics.Scale(dstImage, watermark)
	watermark = dstImage

	m := image.NewRGBA(img.Bounds())
	draw.Draw(m, img.Bounds(), img, image.Point{}, draw.Src)
	draw.Draw(m, watermark.Bounds(), watermark, image.Point{}, draw.Over)

	file.Base = strings.Replace(file.Base, "."+file.Ext, "_watermarked."+file.Ext, 1)
	file.FullPath = filepath.Join(file.Dir, file.Base)
	imgw, err := os.Create(file.FullPath)
	if err != nil {
		panic(err)
	}
	defer imgw.Close()

	switch file.Ext {
	case "jpeg", "jpg":
		jpeg.Encode(imgw, m, &jpeg.Options{Quality: jpeg.DefaultQuality})
	case "png":
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		encoder.Encode(imgw, m)
	case "gif":
		gif.Encode(imgw, m, &gif.Options{NumColors: 256, Quantizer: nil, Drawer: nil})
	default:
		panic(fmt.Errorf("unsupported image format"))
	}
}
