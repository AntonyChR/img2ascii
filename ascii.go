package img2ascii

import (
	"math"
)

// [rmin-rmax] -> [tmin-tmax]}
//
// "value" must be in range [rmin-rmax]
func normailize(value float32, rmin float32, rmax float32, tmin float32, tmax int) float64 {
	return float64(((value-rmin)/(rmax-rmin))*(float32(tmax)-tmin) + tmin)
}

type AsciiMatrix [][]string

func ToAscii(rgbaPixels [][]Pixel, densityMap string) (AsciiMatrix, error) {
	var acsiiMatrix AsciiMatrix

	for _, row := range rgbaPixels {
		var tempRow []string
		for _, pixel := range row {
			i := int(math.Floor(normailize(pixel.Average(), 0, 255, 0, len(densityMap)-2)))

			char := string([]rune(densityMap)[i])
			tempRow = append(tempRow, char)

		}

		acsiiMatrix = append(acsiiMatrix, tempRow)
	}

	return acsiiMatrix, nil
}
