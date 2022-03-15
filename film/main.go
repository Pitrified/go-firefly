package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"os"
	"path/filepath"

	"github.com/Pitrified/go-firefly"
)

func film(
	cellSize, cw, ch,
	nudgeRadius,
	nF,
	filmDuration int,
	drawCircle bool,
) {

	// various simulation params
	clockTickLen := 25_000
	blinkCooldown := 500_000
	nudgeAmount := 20_000
	periodMin := 900_000
	periodMax := 1_100_000

	// decay = 1.0 / 600_000.0
	// rand.Seed(time.Now().UnixNano())

	// film parameters
	fps := 25
	scale := 3
	frameSize := image.Rect(0, 0, cw*cellSize*scale, ch*cellSize*scale)

	// path
	// outputFolder := fmt.Sprintf("film_%v", time.Now().Unix())
	outputFolder := fmt.Sprintf("film_%v", 0)
	fmt.Printf("outputFolder = %+v\n", outputFolder)
	err := os.Mkdir(outputFolder, 0755)
	check(err)

	// setup the blit map
	// fireflyBlit := genBlitMap()

	// background color
	backCol := elemColor['A'].GetBlent(1)

	// start world
	w := firefly.NewWorld(
		cw, ch, float32(cellSize),
		1_000_000, clockTickLen,
		nudgeAmount, float32(nudgeRadius),
		blinkCooldown,
		periodMin, periodMax,
	)
	w.HatchFireflies(nF)

	for frameI := 0; frameI < filmDuration*fps; frameI++ {

		// ########## //
		//   render   //
		// ########## //

		img := image.NewRGBA(frameSize)

		// fill background
		draw.Draw(
			img, img.Bounds(),
			&image.Uniform{backCol},
			image.Point{0, 0},
			draw.Src,
		)

		// save the frame
		frameName := fmt.Sprintf("frame_%06d.png", frameI)
		fmt.Printf("frameName = %+v\n", frameName)
		framePath := filepath.Join(outputFolder, frameName)
		SavePNG(framePath, img)

		// ########## //
		//  simulate  //
		// ########## //

		if frameI == 0 {
			break
		}
	}

}

func main() {
	fmt.Println("Start filming.")

	// world params
	cw := flag.Int("cw", 16, "Width of the world in cells.")
	ch := flag.Int("ch", 9, "Height of the world in cells.")
	cellSize := flag.Int("cs", 80, "Size of each cell.")
	nudgeRadius := flag.Int("nr", 22, "Max distance between interacting fireflies.")
	nF := flag.Int("nf", 1000, "Number of fireflies to simulate.")

	// film params
	filmDuration := flag.Int("fd", 10, "Lenght of the output in seconds.")
	drawCircle := flag.Bool("dc", false, "Draw a circle to show the nudge radius value.")

	flag.Parse()

	fmt.Println("cs    :", *cellSize)
	fmt.Println("cw ch :", *cw, *ch)
	fmt.Println("nr    :", *nudgeRadius)
	fmt.Println("nf    :", *nF)
	fmt.Println("fd    :", *filmDuration)
	fmt.Println("dc    :", *drawCircle)

	film(
		*cellSize,
		*cw, *ch,
		*nudgeRadius,
		*nF,
		*filmDuration,
		*drawCircle,
	)
}
