package main

import (
	"flag"
	"fmt"

	"github.com/Pitrified/go-firefly"
)

func film(
	cellSize, cw, ch,
	nudgeRadius,
	filmDuration int,
	drawCircle bool,
) {

	nF := 1000

	clockTickLen := 25_000
	blinkCooldown := 500_000
	nudgeAmount := 20_000

	periodMin := 900_000
	periodMax := 1_100_000

	// setup the blit map
	genBlitMap()

	// decay = 1.0 / 600_000.0
	// rand.Seed(time.Now().UnixNano())

	w := firefly.NewWorld(
		cw, ch, float32(cellSize),
		1_000_000, clockTickLen,
		nudgeAmount, float32(nudgeRadius),
		blinkCooldown,
		periodMin, periodMax,
	)
	w.HatchFireflies(nF)

}

func main() {
	fmt.Println("Start filming.")

	// world params
	cw := flag.Int("cw", 16, "Width of the world in cells.")
	ch := flag.Int("ch", 9, "Height of the world in cells.")
	cellSize := flag.Int("cs", 80, "Size of each cell.")
	nudgeRadius := flag.Int("nr", 22, "Max distance between interacting fireflies.")

	// film params
	filmDuration := flag.Int("fd", 10, "Lenght of the output in seconds.")
	drawCircle := flag.Bool("dc", false, "Draw a circle to show the nudge radius value.")

	flag.Parse()

	fmt.Println("cs    :", *cellSize)
	fmt.Println("cw ch :", *cw, *ch)
	fmt.Println("nr    :", *nudgeRadius)
	fmt.Println("fd    :", *filmDuration)
	fmt.Println("dc    :", *drawCircle)

	film(
		*cellSize,
		*cw, *ch,
		*nudgeRadius,
		*filmDuration,
		*drawCircle,
	)
}
