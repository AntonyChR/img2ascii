package img2ascii

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"

	"io"
)

type Pixel struct {
	R int
	G int
	B int
	A int
}

func (p *Pixel) Average() float32 {
	return float32((p.R + p.G + p.B) / 3)
}

func getPixelsFromImage(file io.Reader, c Config) ([][]Pixel, error) {

	switch c.Extension {
	case "png":
		image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	case "jpg":
		image.RegisterFormat("jpg", "jpg", jpeg.Decode, jpeg.DecodeConfig)
	case "jpeg":
		image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	default:
		return nil, errors.New("invalid format (png, jpg, jpeg)")
	}

	pixels, err := getPixels(file)

	if err != nil {
		return nil, err
	}

	return pixels, nil
}

func getPixels(file io.Reader) ([][]Pixel, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return [][]Pixel{}, err
	}

	bounds := img.Bounds()

	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, nil

}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func Resize(rgbaMatrix [][]Pixel, n int) [][]Pixel {
	if n == 0 {
		return rgbaMatrix
	}

	return Resize(scale(rgbaMatrix), n-1)
}

func scale(rgbaMatrix [][]Pixel) [][]Pixel {
	height, width := len(rgbaMatrix), len(rgbaMatrix[0])

	var newMatrix [][]Pixel

	for i := 0; i < height; i += 2 {
		var tempRow []Pixel
		for j := 0; j < width; j += 2 {
			tempRow = append(tempRow, rgbaMatrix[i][j])
		}
		newMatrix = append(newMatrix, tempRow)
	}

	return newMatrix
}
