package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"os"
	"path/filepath"
	"sync"

	"github.com/Pitrified/go-firefly"
	"github.com/lucasb-eyer/go-colorful"
)

type Filmer struct {

	// input
	cellSize, cw, ch int
	nudgeRadius      int
	nF               int
	filmDuration     int
	drawCircle       bool

	// utils
	blitTemplate  *image.RGBA
	backCol       colorful.Color
	w             *firefly.World
	fps           int
	scale         int
	frameSize     image.Rectangle
	outputFolder  string
	renderWG      sync.WaitGroup
	decay         float64
	lLevels       int
	templateSize  int
	rotNum        int
	whichTemplate string

	// misc
	clockTickLen  int
	blinkCooldown int
	nudgeAmount   int
	periodMin     int
	periodMax     int
}

func NewFilmer(
	cellSize, cw, ch,
	nudgeRadius,
	nF,
	filmDuration int,
	drawCircle bool,
) *Filmer {

	f := &Filmer{}

	f.cellSize = cellSize
	f.cw = cw
	f.ch = ch
	f.nudgeRadius = nudgeRadius
	f.nF = nF
	f.filmDuration = filmDuration
	f.drawCircle = drawCircle

	return f
}

func (f *Filmer) film() {

	// various simulation params
	f.clockTickLen = 25_000
	f.blinkCooldown = 500_000
	f.nudgeAmount = 20_000
	f.periodMin = 900_000
	f.periodMax = 1_100_000

	// f.whichTemplate = "F3"
	f.whichTemplate = "F5"
	// f.whichTemplate = "L5"

	// rand.Seed(time.Now().UnixNano())
	f.decay = 1.0 / 600_000.0

	// film parameters
	f.fps = 25
	// TODO this might also be linked to which template you are using
	f.scale = 1
	f.frameSize = image.Rect(0, 0, f.cw*f.cellSize*f.scale, f.ch*f.cellSize*f.scale)

	// path
	// outputFolder := fmt.Sprintf("film_%v", time.Now().Unix())
	f.outputFolder = fmt.Sprintf("film_%v", 2)
	fmt.Printf("outputFolder = %+v\n", f.outputFolder)
	err := os.RemoveAll(f.outputFolder)
	check(err)
	err = os.Mkdir(f.outputFolder, 0755)
	check(err)

	// number of lightness levels (-1 as it is inclusive)
	f.lLevels = 100
	// setup the blit map
	f.blitTemplate = genBlitMap(f.lLevels, f.whichTemplate)
	// TODO this is dependent on which template you are using
	switch f.whichTemplate {
	case "F3":
		f.templateSize = 3
		f.rotNum = 2
	case "L5":
		f.templateSize = 5
		f.rotNum = 1
	case "F5":
		f.templateSize = 5
		f.rotNum = 2
	}

	// background color
	f.backCol = elemColor['a'].GetBlent(1)

	// start world
	f.w = firefly.NewWorld(
		f.cw, f.ch, float32(f.cellSize),
		1_000_000, f.clockTickLen,
		f.nudgeAmount, float32(f.nudgeRadius),
		f.blinkCooldown,
		f.periodMin, f.periodMax,
	)
	f.w.HatchFireflies(f.nF)
	// firefly.NewFirefly(100, 100, 0, 0, 1000000, f.w)
	// firefly.NewFirefly(100, 110, 45, 1, 1000000, f.w)
	// firefly.NewFirefly(90, 110, 90, 2, 1000000, f.w)
	// firefly.NewFirefly(80, 110, 135, 3, 1000000, f.w)
	// firefly.NewFirefly(80, 100, 180, 4, 1000000, f.w)
	// firefly.NewFirefly(80, 90, 225, 5, 1000000, f.w)
	// firefly.NewFirefly(90, 90, 270, 6, 1000000, f.w)
	// firefly.NewFirefly(100, 90, 315, 7, 1000000, f.w)

	for frameI := 0; frameI < f.filmDuration*f.fps; frameI++ {

		// ########## //
		//   render   //
		// ########## //

		fmt.Printf("render frameI = %+v\n", frameI)
		f.renderFrame(frameI)

		// ########## //
		//  simulate  //
		// ########## //

		fmt.Printf("simulate frameI = %+v\n", frameI)
		// f.w.DoStep <- 'M'
		f.w.Move()
		f.w.ClockTick()

		// if frameI == 100 {
		// 	break
		// }
	}

	// to turn the frames into a video:
	// ffmpeg -framerate 25 -i frame_%06d.png -c:v libx264 -r 25 -pix_fmt yuv420p out.mp4
	// https://trac.ffmpeg.org/wiki/Slideshow
	// https://stackoverflow.com/questions/24961127/how-to-create-a-video-from-images-with-ffmpeg
	// https://hamelot.io/visualization/using-ffmpeg-to-convert-a-set-of-images-into-a-video/
}

func (f *Filmer) renderFrame(frameI int) {

	img := image.NewRGBA(f.frameSize)

	// fill background
	draw.Draw(
		img, img.Bounds(),
		&image.Uniform{f.backCol},
		image.Point{0, 0},
		draw.Src,
	)

	// draw each cell
	for i := 0; i < f.w.CellWNum; i++ {
		for ii := 0; ii < f.w.CellHNum; ii++ {
			f.renderWG.Add(1)
			go f.renderCell(f.w.Cells[i][ii], img)
		}
	}
	f.renderWG.Wait()

	// save the frame
	frameName := fmt.Sprintf("frame_%06d.png", frameI)
	fmt.Printf("frameName = %+v\n", frameName)
	framePath := filepath.Join(f.outputFolder, frameName)
	// img = UpscaleImg(img, 5)
	SavePNG(framePath, img)
}

func (F *Filmer) renderCell(c *firefly.Cell, m *image.RGBA) {

	for _, f := range c.Fireflies {
		// blit the right firefly in the right place

		// get the lightness level
		since := F.w.Clock - (f.NextBlink - f.Period)
		br := Brightness(since, F.decay)
		lLev := int(br * float64(F.lLevels))

		// go from firefly to template reference system
		remappedOri := remapOri(f.O)
		// find the corner of the rect in the source (template) image
		bX, bY := findBlitPos(remappedOri, lLev, F.templateSize, F.rotNum)
		// rectangle in the source image
		sr := image.Rect(bX, bY, bX+F.templateSize, bY+F.templateSize)
		// corner of the rect in the dest image
		dp := image.Pt(int(f.X*float32(F.scale)), int(f.Y*float32(F.scale)))
		// rectangle in the dest image
		dr := image.Rectangle{dp, dp.Add(sr.Size())}
		draw.Draw(m, dr, F.blitTemplate, sr.Min, draw.Src)

	}

	F.renderWG.Done()
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

	f := NewFilmer(
		*cellSize, *cw, *ch,
		*nudgeRadius,
		*nF,
		*filmDuration,
		*drawCircle,
	)

	f.film()
}
