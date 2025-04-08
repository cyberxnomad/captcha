package captcha

import (
	"image"
	"image/color"
	"math/rand/v2"
)

func randomString(length int, set string) string {
	b := make([]byte, length)

	for i := range b {
		b[i] = set[rand.IntN(len(set))]
	}

	return string(b)
}

func pixelBounds(src *image.RGBA) image.Rectangle {
	bounds := src.Bounds()

	minX, minY := bounds.Max.X, bounds.Max.Y
	maxX, maxY := bounds.Min.X, bounds.Min.Y
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if src.RGBAAt(x, y).A > 0 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	return image.Rect(minX, minY, maxX, maxY)
}

func randomNearColor(base color.Color) color.RGBA {
	r, g, b, a := base.RGBA()
	rOffset, gOffset, bOffset := rand.IntN(50)-25, rand.IntN(50)-25, rand.IntN(50)-25
	return color.RGBA{
		R: clamp(uint8(r>>8)+uint8(rOffset), 0, 255),
		G: clamp(uint8(g>>8)+uint8(gOffset), 0, 255),
		B: clamp(uint8(b>>8)+uint8(bOffset), 0, 255),
		A: uint8(a >> 8),
	}
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func clamp(v, min, max uint8) uint8 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
