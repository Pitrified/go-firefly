package main

import (
	"flag"
	"fmt"
)

func film(
	cellSize, cw, ch,
	nudgeRadius,
	filmDuration int,
	drawCircle bool,
) {

	// decay = 1.0 / 600_000.0
	// rand.Seed(time.Now().UnixNano())

	// a.w = firefly.NewWorld(
	// 	a.wCellW, a.wCellH, float32(a.wCellSize),
	// 	1_000_000, a.clockTickLen,
	// 	a.nudgeAmount, a.nudgeRadius,
	// 	a.blinkCooldown,
	// 	a.periodMin, a.periodMax,
	// )
	// a.w.HatchFireflies(a.nF)

}

func main() {
	fmt.Println("Start filming.")

	// world params
	cw := flag.Int("cw", 16, "Width of the world in cells.")
	ch := flag.Int("ch", 9, "Height of the world in cells.")
	cellSize := flag.Int("cs", 120, "Size of each cell.")
	nudgeRadius := flag.Int("nr", 12, "Max distance between interacting fireflies.")

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

	// TryColorful()
	GenBlitMap()
}
