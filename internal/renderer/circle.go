package renderer

import (
	"math"

	"github.com/go-gl/gl/v3.2-compatibility/gl"
)

// RenderCircle draws an opaque black disk with a gray outline.
func RenderCircle(centerX, centerY, radius float32, config Config) {
	const segments = 64

	vertex := func(angle float64) {
		x := centerX + radius*float32(math.Cos(angle))
		y := centerY + radius*float32(math.Sin(angle))
		gl.Vertex2f(pixelToClipX(x, config), pixelToClipY(y, config))
	}

	gl.Color3f(0, 0, 0)
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Vertex2f(
		pixelToClipX(centerX, config),
		pixelToClipY(centerY, config),
	)
	for i := 0; i <= segments; i++ {
		vertex(float64(i) * 2 * math.Pi / segments)
	}
	gl.End()

	gl.Color3f(0.4, 0.4, 0.4)
	gl.Begin(gl.LINE_LOOP)
	for i := 0; i < segments; i++ {
		vertex(float64(i) * 2 * math.Pi / segments)
	}
	gl.End()
}
