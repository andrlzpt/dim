package renderer

import (
	"github.com/andrlzpt/dim/internal/ascii"
	"github.com/go-gl/gl/v3.2-compatibility/gl"
)

const textVertexFloats = 7 // X, Y, U, V, R, G, B.

// These resources belong to the current single-window rendering context.
var textVertexBuffer uint32
var textVertices []float32

func beginTextBatch() {
	initializeFontAtlas()

	if textVertexBuffer == 0 {
		gl.GenBuffers(1, &textVertexBuffer)
	}

	textVertices = textVertices[:0]
}

func appendTextVertex(
	x, y, u, v float32,
	color ascii.Color,
	config Config,
) {
	textVertices = append(
		textVertices,
		pixelToClipX(x, config),
		pixelToClipY(y, config),
		u,
		v,
		color.R,
		color.G,
		color.B,
	)
}

func drawTextBatch() {
	if len(textVertices) == 0 {
		return
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.ClientActiveTexture(gl.TEXTURE0)
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, fontAtlasTexture)
	gl.TexEnvi(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.MODULATE)

	// Discard transparent texels so glyph backgrounds don't cover other cells.
	gl.Enable(gl.ALPHA_TEST)
	gl.AlphaFunc(gl.GREATER, 0.5)

	gl.BindBuffer(gl.ARRAY_BUFFER, textVertexBuffer)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(textVertices)*4,
		gl.Ptr(textVertices),
		gl.STREAM_DRAW,
	)

	const stride = textVertexFloats * 4

	gl.EnableClientState(gl.VERTEX_ARRAY)
	gl.EnableClientState(gl.TEXTURE_COORD_ARRAY)
	gl.EnableClientState(gl.COLOR_ARRAY)

	gl.VertexPointer(2, gl.FLOAT, stride, gl.PtrOffset(0))
	gl.TexCoordPointer(2, gl.FLOAT, stride, gl.PtrOffset(2*4))
	gl.ColorPointer(3, gl.FLOAT, stride, gl.PtrOffset(4*4))

	gl.DrawArrays(
		gl.QUADS,
		0,
		int32(len(textVertices)/textVertexFloats),
	)

	gl.DisableClientState(gl.COLOR_ARRAY)
	gl.DisableClientState(gl.TEXTURE_COORD_ARRAY)
	gl.DisableClientState(gl.VERTEX_ARRAY)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.Disable(gl.ALPHA_TEST)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.Disable(gl.TEXTURE_2D)
}
