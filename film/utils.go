package main

import (
	"fmt"
	"image"
	"image/png"
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
