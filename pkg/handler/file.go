package handler

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"

	"github.com/httpsOmkar/graphics-go/graphics"
	data "github.com/sarumaj/water-maker/pkg/data"
)

type File struct {
	FullPath, Dir, Base, Ext string
}

func (file *File) SetWatermark() {
	img, err := decodeImage(file.FullPath, nil)
	if err != nil {
		panic(err)
	}

	watermark, err := decodeWatermark()
	if err != nil {
		panic(err)
	}

	file.Base = fmt.Sprintf("%s_watermarked.%s", file.Base, file.Ext)
	file.FullPath = filepath.Join(file.Dir, file.Base)
	if err := encodeImage(file.FullPath, drawWatermark(img, watermark)); err != nil {
		panic(err)
	}
}

func decodeImage(path string, fsys fs.FS) (img image.Image, err error) {
	var imgb io.ReadCloser
	if fsys == nil {
		imgb, err = os.Open(path)
	} else {
		imgb, err = fsys.Open(path)
	}
	if err != nil {
		return nil, err
	}
	defer imgb.Close()

	switch filepath.Ext(path) {
	case ".png":
		img, err = png.Decode(imgb)
	case ".jpeg", ".jpg":
		img, err = jpeg.Decode(imgb)
	case ".gif":
		img, err = gif.Decode(imgb)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", path)
	}
	if err != nil {
		return nil, err
	}

	return img, nil
}

func decodeWatermark() (wm image.Image, err error) {
	wmfile := os.Getenv("WATERMARK_FILE")
	if wmfile != "" {
		return decodeImage(wmfile, nil)
	}
	return decodeImage("images/watermark.png", data.Fs)
}

func drawWatermark(img, watermark image.Image) image.Image {
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
	return m
}

func encodeImage(path string, img image.Image) error {
	imgw, err := os.Create(path)
	if err != nil {
		return err
	}
	defer imgw.Close()

	switch filepath.Ext(path) {
	case "jpeg", "jpg":
		jpeg.Encode(imgw, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
	case "png":
		(&png.Encoder{CompressionLevel: png.BestCompression}).Encode(imgw, img)
	case "gif":
		gif.Encode(imgw, img, &gif.Options{NumColors: 256, Quantizer: nil, Drawer: nil})
	default:
		return fmt.Errorf("unsupported image format")
	}

	return nil
}
