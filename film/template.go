package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

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
// to upscale
// draw.Draw(img,
// image.Rect(i*blockw, 160, (i+1)*blockw, 200),
// &image.Uniform{c1.BlendHcl(c2, float64(i)/float64(blocks-1)).Clamped()},
// image.Point{}, draw.Src)

type RangeColorHCL struct {
	H, C, Lh, Ll float64
	cHCLh, cHCLl colorful.Color
}

func NewRangeColorHCL(H, C, Lh, Ll float64) *RangeColorHCL {
	r := &RangeColorHCL{
		H:     H,
		C:     C,
		Lh:    Lh,
		Ll:    Ll,
		cHCLh: colorful.Hcl(H, C, Lh),
		cHCLl: colorful.Hcl(H, C, Ll),
	}
	return r
}

// Get the blent color.
func (r RangeColorHCL) GetBlent(t float64) colorful.Color {
	return r.cHCLh.BlendHcl(r.cHCLl, t)
}

var elemColor = map[byte]*RangeColorHCL{
	'H': NewRangeColorHCL(60, 0.7, 0.7, 0.7), // Head
	'B': NewRangeColorHCL(60, 0.7, 0.7, 0.7), // Body
	'O': NewRangeColorHCL(60, 0.7, 0.7, 0.2), // bOdy glowing
	'W': NewRangeColorHCL(60, 0.7, 0.7, 0.7), // Wings
	'I': NewRangeColorHCL(60, 0.7, 0.7, 0.2), // wIngs glowing
	'A': NewRangeColorHCL(60, 0.7, 0.7, 0.7), // bAckground
	'C': NewRangeColorHCL(60, 0.7, 0.7, 0.2), // baCkground glowing
}

var templateFirefly = [][]byte{
	{'A', 'H', 'A'},
	{'I', 'O', 'I'},
	{'W', 'O', 'W'},
}

func tryColorful() {
	c := colorful.Hcl(80, 0.3, 0.5)
	fmt.Printf("RGB values: %v, %v, %v\n", c.R, c.G, c.B)

	blocks := 11
	blockw := 40
	img := image.NewRGBA(image.Rect(0, 0, blocks*blockw, 200))

	fmt.Printf("elemColor['O'] = %+v\n", elemColor['O'])
	for i := 0; i <= 10; i++ {
		t := float64(i) * 0.1
		col := elemColor['O'].GetBlent(t)
		fmt.Printf("elemColor['O'].GetBlent(%v) = %+v\n",
			t,
			col,
		)

		draw.Draw(img,
			image.Rect(i*blockw, 160, (i+1)*blockw, 200),
			&image.Uniform{col},
			image.Point{}, draw.Src)
	}
	SavePNG("testBlend.png", img)

	fmt.Printf("templateFirefly = %+v\n", templateFirefly)

}

func SavePNG(name string, img image.Image) {
	toimg, err := os.Create("colorblend.png")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	defer toimg.Close()

	png.Encode(toimg, img)
}
