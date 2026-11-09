package handler

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

func TestEncodeImageRoundTrip(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 4, 3))
	img.Set(1, 1, color.RGBA{R: 255, A: 255})

	for _, ext := range []string{".png", ".jpg", ".jpeg", ".gif"} {
		t.Run(ext, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "image"+ext)
			if err := encodeImage(path, img); err != nil {
				t.Fatalf("encodeImage: %v", err)
			}
			got, err := decodeImage(path, nil)
			if err != nil {
				t.Fatalf("decodeImage: %v", err)
			}
			if got.Bounds() != img.Bounds() {
				t.Fatalf("bounds: got %v, want %v", got.Bounds(), img.Bounds())
			}
		})
	}
}

func TestEncodeImageUnsupportedFormat(t *testing.T) {
	path := filepath.Join(t.TempDir(), "image.bmp")
	if err := encodeImage(path, image.NewRGBA(image.Rect(0, 0, 1, 1))); err == nil {
		t.Fatal("expected an error for .bmp")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected no file to be created, got %v", err)
	}
}
