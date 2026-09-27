package renderer

import "github.com/andrlzpt/dim/internal/ascii"

// RenderText draws a character buffer in one batch.
func RenderText(buffer *ascii.Buffer, config Config) {
	RenderTextAt(buffer, 0, 0, config)
}

// RenderText draws a character buffer in one batch.
// Call BeginFrame first.
func RenderTextAt(buffer *ascii.Buffer, left, top float32, config Config) {
	beginTextBatch()

	transform := func(x, y float32) (float32, float32) {
		return left + x, top + y
	}

	for y := 0; y < buffer.Height; y++ {
		for x := 0; x < buffer.Width; x++ {
			cell := buffer.Cells[y*buffer.Width+x]
			if cell.Rune == 0 || cell.Rune == ' ' || cell.Brightness <= 0 {
				continue
			}

			appendTextGlyph(
				x,
				y,
				cell.Rune,
				cell.Color,
				cell.Brightness,
				config,
				transform,
			)
		}
	}

	drawTextBatch()
}

// RenderTextLine draws a centered label at pixel coordinates.
// quarterTurns: 0 horizontal, 1 clockwise, -1 counterclockwise.
func RenderTextLine(
	label string,
	centerX, centerY float32,
	quarterTurns int,
	config Config,
) {
	runes := []rune(label)
	width := float32(len(runes) * config.CellWidth)
	height := float32(config.CellHeight)

	transform := func(x, y float32) (float32, float32) {
		x -= width / 2
		y -= height / 2

		switch quarterTurns {
		case 1:
			x, y = -y, x
		case -1:
			x, y = y, -x
		}

		return centerX + x, centerY + y
	}

	beginTextBatch()

	for x, r := range runes {
		if r != ' ' {
			appendTextGlyph(x, 0, r, ascii.RGB(0xFFFFFF), 1, config, transform)
		}
	}

	drawTextBatch()
}

// appendTextGlyph adds one textured rectangle without issuing GL calls.
func appendTextGlyph(
	cellX, cellY int,
	r rune,
	color ascii.Color,
	brightness float32,
	config Config,
	transform func(float32, float32) (float32, float32),
) {
	uv, found := fontAtlasGlyphs[r]
	if !found {
		uv = fontAtlasGlyphs['?']
	}

	glyphWidth := config.GlyphWidth
	if glyphWidth <= 0 {
		glyphWidth = config.CellWidth
	}

	glyphHeight := config.GlyphHeight
	if glyphHeight <= 0 {
		glyphHeight = config.CellHeight
	}

	// Apply brightness once per glyph, preserving its color.
	color.R *= brightness
	color.G *= brightness
	color.B *= brightness

	left := float32(cellX * config.CellWidth)
	top := float32(cellY * config.CellHeight)
	right := left + float32(glyphWidth)
	bottom := top + float32(glyphHeight)

	vertex := func(x, y, u, v float32) {
		if transform != nil {
			x, y = transform(x, y)
		}
		appendTextVertex(x, y, u, v, color, config)
	}

	vertex(left, top, uv.Left, uv.Top)
	vertex(right, top, uv.Right, uv.Top)
	vertex(right, bottom, uv.Right, uv.Bottom)
	vertex(left, bottom, uv.Left, uv.Bottom)
}
