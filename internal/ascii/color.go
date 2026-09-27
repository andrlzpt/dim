package ascii

// Color contains RGB components between zero and one.
type Color struct {
	R, G, B float32
}

// RGB converts a hexadecimal color such as 0xD8CFAD.
func RGB(hex uint32) Color {
	return Color{
		R: float32((hex>>16)&0xFF) / 255,
		G: float32((hex>>8)&0xFF) / 255,
		B: float32(hex&0xFF) / 255,
	}
}
