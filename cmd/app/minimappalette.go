package main

import "github.com/andrlzpt/dim/internal/ascii"

var (
	minimapBackground = ascii.RGB(0x090B0A) // Almost-black charcoal.
	minimapBorder     = ascii.RGB(0x45483F) // Dark ash.
	minimapRooms      = ascii.RGB(0x777B70) // Ash gray.
	minimapPlayer     = ascii.RGB(0xD8CFAD) // Aged bone.
	minimapObjective  = ascii.RGB(0xC49A57) // Tarnished brass.
	minimapReached    = ascii.RGB(0x809078) // Faded sage.
)

func minimapGlyphColor(glyph rune) ascii.Color {
	switch glyph {
	case '@', '>', 'v', '<', '^':
		return minimapPlayer
	case '*', 'X':
		return minimapObjective
	case 'V':
		return minimapReached
	default:
		return minimapRooms
	}
}
