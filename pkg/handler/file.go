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
	watermarked, err := drawWatermark(img, watermark)
	if err != nil {
		panic(err)
	}
	if err := encodeImage(file.FullPath, watermarked); err != nil {
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
	defer func() { _ = imgb.Close() }()

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

func drawWatermark(img, watermark image.Image) (image.Image, error) {
	if s := img.Bounds().Size(); s.Y > s.X {
		dstImage := image.NewRGBA(image.Rect(0, 0, watermark.Bounds().Dy(), watermark.Bounds().Dx()))
		if err := graphics.Rotate(dstImage, watermark, &graphics.RotateOptions{Angle: 3.0 * math.Pi / 2.0}); err != nil {
			return nil, err
		}
		watermark = dstImage
	}
	dstImage := image.NewRGBA(img.Bounds())
	if err := graphics.Scale(dstImage, watermark); err != nil {
		return nil, err
	}
	watermark = dstImage

	m := image.NewRGBA(img.Bounds())
	draw.Draw(m, img.Bounds(), img, image.Point{}, draw.Src)
	draw.Draw(m, watermark.Bounds(), watermark, image.Point{}, draw.Over)
	return m, nil
}

func encodeImage(path string, img image.Image) (err error) {
	var encode func(io.Writer, image.Image) error
	switch filepath.Ext(path) {
	case ".jpeg", ".jpg":
		encode = func(w io.Writer, m image.Image) error {
			return jpeg.Encode(w, m, &jpeg.Options{Quality: jpeg.DefaultQuality})
		}
	case ".png":
		encode = (&png.Encoder{CompressionLevel: png.BestCompression}).Encode
	case ".gif":
		encode = func(w io.Writer, m image.Image) error {
			return gif.Encode(w, m, &gif.Options{NumColors: 256, Quantizer: nil, Drawer: nil})
		}
	default:
		return fmt.Errorf("unsupported image format: %s", path)
	}

	imgw, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := imgw.Close(); err == nil {
			err = cerr
		}
	}()

	return encode(imgw, img)
}
