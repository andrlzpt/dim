package renderer

import (
	"github.com/go-gl/gl/v3.2-compatibility/gl"
	"golang.org/x/image/font/basicfont"
)

type glyphUV struct {
	Left, Top, Right, Bottom float32
}

var fontAtlasTexture uint32
var fontAtlasGlyphs map[rune]glyphUV

// initializeFontAtlas runs after BeginFrame has initialized OpenGL.
func initializeFontAtlas() {
	if fontAtlasTexture != 0 {
		return
	}

	face := basicfont.Face7x13
	var characters []rune
	for _, interval := range face.Ranges {
		for r := interval.Low; r < interval.High; r++ {
			characters = append(characters, r)
		}
	}

	const columns = 16
	const padding = 1

	slotWidth := face.Advance + padding*2
	slotHeight := face.Height + padding*2
	rows := (len(characters) + columns - 1) / columns

	width := columns * slotWidth
	height := rows * slotHeight
	pixels := make([]byte, width*height*4)
	fontAtlasGlyphs = make(map[rune]glyphUV, len(characters))

	for index, r := range characters {
		left := (index%columns)*slotWidth + padding
		top := (index/columns)*slotHeight + padding

		for _, pixel := range cachedGlyphPixels(r) {
			x := left + int(pixel.X)
			y := top + int(pixel.Y)
			offset := (y*width + x) * 4

			pixels[offset] = 255
			pixels[offset+1] = 255
			pixels[offset+2] = 255
			pixels[offset+3] = 255
		}

		fontAtlasGlyphs[r] = glyphUV{
			Left:   float32(left) / float32(width),
			Top:    float32(top) / float32(height),
			Right:  float32(left+face.Advance) / float32(width),
			Bottom: float32(top+face.Height) / float32(height),
		}
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.GenTextures(1, &fontAtlasTexture)
	gl.BindTexture(gl.TEXTURE_2D, fontAtlasTexture)

	// Nearest filtering keeps the bitmap glyphs crisp.
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA8,
		int32(width),
		int32(height),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(pixels),
	)

	gl.BindTexture(gl.TEXTURE_2D, 0)
}
