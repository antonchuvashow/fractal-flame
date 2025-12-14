package infrastructure

import (
	"image"
	"image/png"
	"os"
)

type PngWriter struct{}

func NewPngWriter() *PngWriter {
	return &PngWriter{}
}

func (w *PngWriter) Write(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
