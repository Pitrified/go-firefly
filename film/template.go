package main

import (
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
)

// CIE-L*C*h° (HCL): This is generally the most useful one; CIE-L*a*b* space in
// polar coordinates, i.e. a better HSV. H° is in [0..360], C* almost in [0..1]
// and L* as in CIE-L*a*b*.

// Blending and clamping
// https://github.com/lucasb-eyer/go-colorful/issues/14#issuecomment-324205385

// elements:
// Head
// Body
// bOdy glowing
// Wings
// wIngs glowing
// bAckground
// baCkground glowing

// blit
// draw.Draw(img, image.Rect(i*blockw,160,(i+1)*blockw,200), &image.Uniform{c1.BlendHcl(c2, float64(i)/float64(blocks-1)).Clamped()}, image.Point{}, draw.Src)

var elColor = map[byte]colorful.Color{
	'M': colorful.Hcl(60, 1, 1),
	'H': colorful.Hcl(60, 0.7, 0.7),
}

type RangeColorHCL struct {
	H, C, Lh, Ll float64
}

var elemColor = map[byte]RangeColorHCL{
	'H': {60, 0.7, 0.7, 0.7}, // Head
	'B': {60, 0.7, 0.7, 0.7}, // Body
	'O': {60, 0.7, 0.7, 0.7}, // bOdy glowing
	'W': {60, 0.7, 0.7, 0.7}, // Wings
	'I': {60, 0.7, 0.7, 0.7}, // wIngs glowing
	'A': {60, 0.7, 0.7, 0.7}, // bAckground
	'C': {60, 0.7, 0.7, 0.7}, // baCkground glowing
}

// TODO
// create a grid with the colors blended
// method GetBlent(t) that returns the blended color?

func tryColorful() {
	c := colorful.Hcl(80, 0.3, 0.5)
	fmt.Printf("RGB values: %v, %v, %v\n", c.R, c.G, c.B)

	fmt.Printf("elColor = %+v\n", elColor)

	fmt.Printf("elColor = %+v\n", elColor['M'].BlendHcl(elColor['M'], 0))
	fmt.Printf("elColor = %+v\n", elColor['M'].BlendHcl(elColor['M'], 0).Clamped())

}
