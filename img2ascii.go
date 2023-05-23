package img2ascii

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Config struct {
	DMap      string // Density map
	Reverse   bool   // Reverse density map
	Extension string
	NewWidth  int
	NewHeight int
}

func (c *Config) GetDensityMapFromTextFile(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	densityMap := string(file)

	if densityMap == "" {
		return errors.New("density map is empty")
	}

	c.DMap = densityMap
	return nil
}

var DEFAULT_DENSITY_MAP = " .:;-~=+*#%@"

func Generate(img io.Reader, c Config) ([][]string, error) {

	if c.DMap == "" {
		c.DMap = DEFAULT_DENSITY_MAP
	}

	if c.Reverse {
		c.DMap = reverseString(c.DMap)
	}

	rgbaPixels, _ := getPixelsFromImage(img, c)
	ascii, _ := ToAscii(ResizeImage(rgbaPixels, c.NewHeight, c.NewWidth), c.DMap)
	return ascii, nil
}

func GenerateTextFile(path string, ascii [][]string) error {
	file, err := os.Create(path)

	if err != nil {
		return err
	}

	for _, row := range ascii {
		newRow := fmt.Sprintf("%s\n", strings.Join(row, ""))
		file.WriteString(newRow)
	}

	return nil
}

func reverseString(s string) string {
	rns := []rune(s)
	var newString []rune
	for i := len(s) - 1; i >= 0; i-- {
		newString = append(newString, rns[i])
	}
	return string(newString)
}
