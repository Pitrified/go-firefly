package main

import (
	"image"
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

// Hue in [0..360], Saturation and Luminance in [0..1].
type RangeColorHSL struct {
	H, S, Lh, Ll float64
	cHSLh, cHSLl colorful.Color
}

func NewRangeColorHSL(H, S, Lh, Ll float64) *RangeColorHSL {
	r := &RangeColorHSL{
		H:     H,
		S:     S,
		Lh:    Lh,
		Ll:    Ll,
		cHSLh: colorful.HSLuv(H, S, Lh),
		cHSLl: colorful.HSLuv(H, S, Ll),
	}
	return r
}

// Get the blent color.
// https://github.com/lucasb-eyer/go-colorful/issues/14#issuecomment-324205385
func (r *RangeColorHSL) GetBlent(t float64) colorful.Color {
	return r.cHSLl.BlendLuvLCh(r.cHSLh, t)
}

var elemColor = map[byte]*RangeColorHSL{
	'H': NewRangeColorHSL(20, 0.5, 0.5, 0.3),   // Head
	'B': NewRangeColorHSL(36, 0.5, 0.01, 0.01), // Body
	'O': NewRangeColorHSL(55, 0.9, 0.7, 0.2),   // bOdy glowing
	'W': NewRangeColorHSL(240, 0.7, 0.2, 0.2),  // Wings
	'I': NewRangeColorHSL(240, 0.7, 0.7, 0.2),  // wIngs glowing
	'A': NewRangeColorHSL(0, 0.0, 0.2, 0.2),    // bAckground
	'C': NewRangeColorHSL(0, 0.0, 0.4, 0.2),    // baCkground glowing
}

// Generate an image with all the needed fireflies to use.
// horizontal change the luminosity
// vertical change the rotation
func genBlitMap() *image.RGBA {

	// firefly templates
	templateFirefly := [][][]byte{
		{
			{'A', 'H', 'A'},
			// {'W', 'B', 'W'},
			{'I', 'O', 'I'},
			{'I', 'O', 'I'},
		},
		{
			// {'A', 'I', 'H'},
			// {'I', 'O', 'I'},
			// {'O', 'I', 'A'},
			{'I', 'I', 'H'},
			{'A', 'O', 'I'},
			{'O', 'A', 'I'},
		},
	}

	// number of lightness levels (-1 as it is inclusive)
	lLevels := 100

	numTemplates := len(templateFirefly)
	fSize := len(templateFirefly[0])

	fireflyRotImg := image.NewRGBA(image.Rect(0, 0, fSize*(lLevels+1), fSize*numTemplates*4))

	// helper to read rotation directions
	rotName := []string{
		"up",
		"right",
		"down",
		"left",
	}

	// iterate over the lightness level
	for il := 0; il <= lLevels; il++ {

		// compute the lightness level in [0,1]
		l := float64(il) * 1.0 / float64(lLevels)
		// how much to shift the template right
		lSh := fSize * il

		// iterate over the template position
		for y := 0; y < fSize; y++ {
			for x := 0; x < fSize; x++ {

				// iterate over the different templates
				for it := 0; it < numTemplates; it++ {

					// get the color to use
					key := templateFirefly[it][y][x] // this is swapped, the first row is y=0
					blend := elemColor[key].GetBlent(l)
					r, g, b := blend.Clamped().RGB255()

					// how much to shift the template down
					tSh := fSize * (numTemplates - 1) * it

					// iterate over the different rotations
					for ir := 0; ir < len(rotName); ir++ {
						// get the rotated coordinates
						xr, yr := getRotatedCoords(x, y, fSize, rotName[ir])
						// how much to shift the template right
						rSh := fSize * numTemplates * ir
						fireflyRotImg.SetRGBA(xr+lSh, yr+tSh+rSh, color.RGBA{r, g, b, 255})
					}
				}
			}
		}
	}

	dst := UpscaleImg(fireflyRotImg, 20)
	SavePNG("testBlitFirefly.png", dst)

	return fireflyRotImg
}

func getRotatedCoords(x, y, size int, rot string) (int, int) {
	switch rot {
	case "up":
		return x, y
	case "right":
		return -y + size - 1, x
	case "down":
		return -x + size - 1, -y + size - 1
	case "left":
		return y, -x + size - 1
	}
	// this should never happen lol
	return 0, 0
}
