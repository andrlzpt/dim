package renderer

import (
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type glyphPixel struct {
	X, Y float32
}

// Accessed only by the rendering thread
var glyphCache = make(map[rune][]glyphPixel)

func cachedGlyphPixels(r rune) []glyphPixel {
	if pixels, found := glyphCache[r]; found {
		return pixels
	}

	face := basicfont.Face7x13
	baseline := fixed.P(0, face.Ascent)
	bounds, mask, maskPoint, _, ok := face.Glyph(baseline, r)
	if !ok {
		pixels := cachedGlyphPixels('?')
		glyphCache[r] = pixels
		return pixels
	}
	var pixels []glyphPixel

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			maskX := maskPoint.X + x - bounds.Min.X
			maskY := maskPoint.Y + y - bounds.Min.Y

			_, _, _, alpha := mask.At(maskX, maskY).RGBA()
			if alpha == 0 {
				continue
			}

			pixels = append(pixels, glyphPixel{
				X: float32(x),
				Y: float32(y),
			})
		}
	}

	glyphCache[r] = pixels
	return pixels
}
