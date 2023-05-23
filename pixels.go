package img2ascii

import (
	"errors"
	"fmt"
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

func ResizeImage(rgbaMatrix [][]Pixel, newHeight, newWidth int) [][]Pixel {
	if newHeight+newWidth == 0 {
		return rgbaMatrix
	}

	height := len(rgbaMatrix)
	width := len(rgbaMatrix[0])

	if newWidth == 0 {
		newWidth = int(float64(newHeight*width) / float64(height))
	}

	if newHeight == 0 {
		newHeight = int(float64(newWidth*height) / float64(width))
	}

	fmt.Printf("[h: %v, w: %v] -> [h: %v, w: %v] \n", height, width, newHeight, newWidth)

	// Calculate the scale factors for resizing
	widthScale := float64(width) / float64(newWidth)
	heightScale := float64(height) / float64(newHeight)

	// Create a new matrix to store the resized pixels
	resizedPixels := make([][]Pixel, newHeight)
	for i := 0; i < newHeight; i++ {
		resizedPixels[i] = make([]Pixel, newWidth)
	}

	// Resize the image
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Calculate the corresponding position in the original image
			origX := int(float64(x) * widthScale)
			origY := int(float64(y) * heightScale)

			// Get the corresponding pixel in the original image
			pixel := rgbaMatrix[origY][origX]

			// Assign the resized pixel to the new matrix
			resizedPixels[y][x] = pixel
		}
	}

	return resizedPixels
}
