package main

import (
	"image"
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

type RangeColorHCL struct {
	H, C, Lh, Ll float64
	cHCLh, cHCLl colorful.Color
}

func NewRangeColorHCL(H, C, Lh, Ll float64) *RangeColorHCL {
	r := &RangeColorHCL{
		H:  H,
		C:  C,
		Lh: Lh,
		Ll: Ll,
		// cHCLh: colorful.Hcl(H, C, Lh),
		// cHCLl: colorful.Hcl(H, C, Ll),
		cHCLh: colorful.HSLuv(H, C, Lh),
		cHCLl: colorful.HSLuv(H, C, Ll),
	}
	return r
}

// Get the blent color.
func (r *RangeColorHCL) GetBlent(t float64) colorful.Color {
	// return r.cHCLh.BlendHcl(r.cHCLl, t)
	return r.cHCLl.BlendLuvLCh(r.cHCLh, t)
}

var elemColor = map[byte]*RangeColorHCL{
	'h': NewRangeColorHCL(20, 0.5, 0.7, 0.2),   // head
	'b': NewRangeColorHCL(36, 0.5, 0.01, 0.01), // body
	'B': NewRangeColorHCL(55, 0.9, 0.7, 0.2),   // Body glowing
	'w': NewRangeColorHCL(240, 0.7, 0.2, 0.2),  // wings
	'W': NewRangeColorHCL(240, 0.7, 0.7, 0.2),  // Wings glowing
	'a': NewRangeColorHCL(0, 0.0, 0.2, 0.2),    // background
	'A': NewRangeColorHCL(0, 0.0, 0.4, 0.2),    // bAckground glowing
	'1': NewRangeColorHCL(55, 0.9, 0.7, 0.2),   // just glow a bit
	'2': NewRangeColorHCL(55, 0.9, 0.6, 0.18),  // just glow a bit
	'3': NewRangeColorHCL(55, 0.9, 0.5, 0.16),  // just glow a bit
	'4': NewRangeColorHCL(55, 0.9, 0.4, 0.14),  // just glow a bit
	'5': NewRangeColorHCL(55, 0.9, 0.3, 0.12),  // just glow a bit
	'6': NewRangeColorHCL(55, 0.9, 0.2, 0.1),   // just glow a bit
	'7': NewRangeColorHCL(55, 0.9, 0.1, 0.08),  // just glow a bit
}

// firefly templates
var TemplateFirefly3 = [][][]byte{
	{
		{'a', 'h', 'a'},
		// {'w', 'b', 'w'},
		{'W', 'B', 'W'},
		{'W', 'B', 'W'},
	},
	{
		// {'a', 'W', 'h'},
		// {'W', 'B', 'W'},
		// {'B', 'W', 'a'},
		{'W', 'W', 'h'},
		{'a', 'B', 'W'},
		{'B', 'a', 'W'},
	},
}

var TemplateSpherical5 = [][][]byte{
	{
		{'5', '4', '4', '4', '5'},
		{'4', '3', '2', '3', '4'},
		{'4', '2', '1', '2', '4'},
		{'4', '3', '2', '3', '4'},
		{'5', '4', '4', '4', '5'},
	},
}

var TemplateFirefly5 = [][][]byte{
	{
		{'a', 'a', 'h', 'a', 'a'},
		{'W', 'W', 'b', 'W', 'W'},
		{'W', 'W', 'B', 'W', 'W'},
		{'W', 'A', 'B', 'A', 'W'},
		{'a', 'A', 'B', 'A', 'a'},
	},
	{
		{'W', 'W', 'W', 'a', 'h'},
		{'a', 'W', 'W', 'b', 'a'},
		{'a', 'A', 'B', 'W', 'W'},
		{'A', 'B', 'A', 'W', 'W'},
		{'B', 'A', 'a', 'a', 'W'},
	},
}

// orientation and lightness level
func findBlitPos(o int16, l, templateSize, rotNum int) (int, int) {

	sectorSize := int16(90 / rotNum)

	tO := o + sectorSize/2
	if tO > 360 {
		tO -= 360
	}
	rotI := tO / sectorSize

	return l * templateSize, int(rotI) * templateSize
}

func remapOri(o int16) int16 {
	// screen | fire
	// 0   0  | 90
	// 45  1  | 45
	// 90  2  | 0
	// 135 3  | -45  315
	// 180 4  | -90  270
	// 225 5  | -135 225
	// 270 6  | -180 180
	// 315 7  | -225 135

	remappedOri := -o + 90
	// fmt.Printf("o = %+v remappedOri = %+v\n", o, remappedOri)

	for remappedOri < 0 {
		remappedOri += 360
	}
	for remappedOri > 360 {
		remappedOri -= 360
	}

	return remappedOri

}

// Generate an image with all the needed fireflies to use.
// horizontal change the luminosity
// vertical change the rotation
func genBlitMap(lLevels int, whichTemplate string) *image.RGBA {

	var templateFirefly [][][]byte
	switch whichTemplate {
	case "F3":
		templateFirefly = TemplateFirefly3
	case "F5":
		templateFirefly = TemplateFirefly5
	case "L5":
		templateFirefly = TemplateSpherical5
	}

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
