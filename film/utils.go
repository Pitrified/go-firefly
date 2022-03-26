package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"

	"golang.org/x/image/draw"
)

func SavePNG(name string, img image.Image) {
	toimg, err := os.Create(name)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	defer toimg.Close()
	png.Encode(toimg, img)
}

// rescale image
// https://gist.github.com/logrusorgru/570d64fd6a051e0441014387b89286ca
func UpscaleImg(img image.Image, pixSize int) *image.RGBA {
	largeSize := image.Rect(0, 0, img.Bounds().Dx()*pixSize, img.Bounds().Dy()*pixSize)
	dst := image.NewRGBA(largeSize)
	draw.NearestNeighbor.Scale(dst, largeSize, img, img.Bounds(), draw.Src, nil)
	return dst
}

// brightness goes from 1 to 0, according to decay constant
// tau = 1 / decay
// https://www.wolframalpha.com/input/?i=e+**+%28+-x+*+1%2F50+%29+for+0+%3C+x+%3C+250
func Brightness(x int, decay float64) float64 {
	if x < 0 {
		return 0
	}
	return math.Exp(-float64(x) * decay)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
