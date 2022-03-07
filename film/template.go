package main

import (
	"fmt"
	"image"
	"image/color"
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
func (r *RangeColorHCL) GetBlent(t float64) colorful.Color {
	return r.cHCLh.BlendHcl(r.cHCLl, t)
}

var elemColor = map[byte]*RangeColorHCL{
	'H': NewRangeColorHCL(40, 0.7, 0.7, 0.7),   // Head
	'B': NewRangeColorHCL(36, 0.5, 0.01, 0.01), // Body
	'O': NewRangeColorHCL(55, 0.9, 0.7, 0.2),   // bOdy glowing
	'W': NewRangeColorHCL(90, 0.7, 0.7, 0.7),   // Wings
	'I': NewRangeColorHCL(120, 0.7, 0.7, 0.2),  // wIngs glowing
	'A': NewRangeColorHCL(60, 0.7, 0.7, 0.7),   // bAckground
	'C': NewRangeColorHCL(180, 0.7, 0.7, 0.2),  // baCkground glowing
	// 'a': NewRangeColorHCL(180, 1, 0.7, 0.2),   // baCkground glowing
	// 'b': NewRangeColorHCL(180, 0.7, 0.7, 0.2), // baCkground glowing
	// 'c': NewRangeColorHCL(180, 0.5, 0.7, 0.2), // baCkground glowing
	// 'd': NewRangeColorHCL(180, 0.3, 0.7, 0.2), // baCkground glowing
	// 'e': NewRangeColorHCL(180, 0.2, 0.7, 0.2), // baCkground glowing
	// 'f': NewRangeColorHCL(180, 0.1, 0.7, 0.2), // baCkground glowing
}

var templateFirefly = [][]byte{
	{'A', 'H', 'A'},
	{'I', 'B', 'I'},
	{'W', 'O', 'W'},
}

var templateFireflyLarge = [][]byte{
	{'A', 'A', 'H', 'A', 'A'},
	{'W', 'W', 'B', 'W', 'W'},
	{'W', 'W', 'O', 'W', 'W'},
	{'W', 'C', 'O', 'C', 'W'},
	{'A', 'C', 'O', 'C', 'A'},
}

func tryColorful() {

	// a random color
	c := colorful.Hcl(80, 0.3, 0.5)
	fmt.Printf("RGB values: %v, %v, %v\n", c.R, c.G, c.B)

	// a random RangeColorHCL
	fmt.Printf("elemColor['O'] = %+v\n", elemColor['O'])

	keys := []byte{
		'O', 'I', 'C',
		'H', 'B', 'W', 'A',
		// 'a', 'b', 'c', 'd', 'e', 'f',
	}

	blocks := 11
	blockw := 40
	img := image.NewRGBA(image.Rect(0, 0, blocks*blockw, blockw*len(keys)))

	for i := 0; i <= 10; i++ {
		t := float64(i) * 0.1

		for ii := 0; ii < len(keys); ii++ {

			key := keys[ii]
			col := elemColor[key].GetBlent(t)
			fmt.Printf("elemColor[%v].GetBlent(%v) = %+v\n",
				key, t, col,
			)
			draw.Draw(img,
				image.Rect(i*blockw, ii*blockw, (i+1)*blockw, (ii+1)*blockw),
				&image.Uniform{col},
				image.Point{}, draw.Src)
		}
	}
	SavePNG("testBlend.png", img)

	fmt.Printf("templateFirefly = %+v\n", templateFirefly)
	fmt.Printf("templateFireflyLarge = %+v\n", templateFireflyLarge)

	img = image.NewRGBA(image.Rect(0, 0, 100, 100))
	H := 60.0
	for i := 0; i <= 100; i++ {
		for ii := 0; ii < 100; ii++ {
			c := float64(i) / 100
			l := float64(ii) / 100
			col := colorful.Hcl(H, c, l)
			img.SetRGBA(i, ii, color.RGBA{uint8(col.R * 255), uint8(col.G * 255), uint8(col.B * 255), 255})
		}
	}
	SavePNG("testHCL.png", img)

	img = image.NewRGBA(image.Rect(0, 0, 100, 100))
	H = 60.0
	for i := 0; i <= 100; i++ {
		for ii := 0; ii < 100; ii++ {
			s := float64(i) / 100
			v := float64(ii) / 100
			col := colorful.Hsv(H, s, v)
			img.SetRGBA(i, ii, color.RGBA{uint8(col.R * 255), uint8(col.G * 255), uint8(col.B * 255), 255})
		}
	}
	SavePNG("testHSV.png", img)

	img = image.NewRGBA(image.Rect(0, 0, 100, 100))
	H = 60.0
	for i := 0; i <= 100; i++ {
		for ii := 0; ii < 100; ii++ {
			s := float64(i) / 100
			v := float64(ii) / 100
			col := colorful.HSLuv(H, s, v)
			img.SetRGBA(i, ii, color.RGBA{uint8(col.R * 255), uint8(col.G * 255), uint8(col.B * 255), 255})
		}
	}
	SavePNG("testHSLuv.png", img)
}

func SavePNG(name string, img image.Image) {
	toimg, err := os.Create(name)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	defer toimg.Close()

	png.Encode(toimg, img)
}
