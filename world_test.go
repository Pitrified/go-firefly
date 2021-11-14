package firefly

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChangeCell(t *testing.T) {
	w := NewWorld(10, 10, 100, 1_000_000, 25_000, 50_000, 50, 500_000, 900_000, 1_1000_000)
	f := NewFirefly(0, 0, 0, 0, 1000000, w)

	c := f.c
	assert.Contains(t, c.Fireflies, f.Id)

	// change cell
	nc := w.Cells[1][1]
	w.chChangeCell <- &ChangeCellReq{f, f.c, nc}
	<-w.chChangeCellDone
	assert.NotContains(t, c.Fireflies, f.Id)
}

func TestMove(t *testing.T) {
	w := NewWorld(10, 10, 100, 1_000_000, 25_000, 50_000, 50, 500_000, 900_000, 1_1000_000)

	// near the top right corner, pointing right
	f := NewFirefly(99.5, 99.5, 0, 0, 1000000, w)
	assert.Contains(t, w.Cells[0][0].Fireflies, f.Id)
	// move to the right
	w.Move()
	assert.Contains(t, w.Cells[1][0].Fireflies, f.Id)
	// move to the top
	f.O = 90
	w.Move()
	assert.Contains(t, w.Cells[1][1].Fireflies, f.Id)
	// move to the left
	f.O = 180
	w.Move()
	assert.Contains(t, w.Cells[0][1].Fireflies, f.Id)
	// move to the bottom
	f.O = 270
	w.Move()
	assert.Contains(t, w.Cells[0][0].Fireflies, f.Id)
}

func TestHatch(t *testing.T) {
	w := NewWorld(10, 10, 100, 1_000_000, 25_000, 50_000, 50, 500_000, 900_000, 1_1000_000)
	nF := 10
	w.HatchFireflies(nF)

	// count how many were created
	tot := 0
	for i := 0; i < w.CellWNum; i++ {
		for ii := 0; ii < w.CellHNum; ii++ {
			tot += len(w.Cells[i][ii].Fireflies)
		}
	}
	assert.Equal(t, nF, tot, fmt.Sprintf("Failed %+v, created %+v", nF, tot))

}

func TestMoveWrap(t *testing.T) {
	w := NewWorld(10, 10, 100, 1_000_000, 25_000, 50_000, 50, 500_000, 900_000, 1_1000_000)

	cases := []struct {
		cx, cy, dcx, dcy, nx, ny int
	}{
		{0, 0, 0, 0, 0, 0},
		{0, 0, -1, -1, 9, 9},
		{0, 0, -101, -101, 9, 9},
		{0, 0, 10, 10, 0, 0},
		{0, 0, 100, 100, 0, 0},
	}
	for _, c := range cases {
		gotX, gotY := w.MoveWrapCell(c.cx, c.cy, c.dcx, c.dcy)
		assert.Equal(t, gotX, c.nx, fmt.Sprintf("Failed case %+v, got %+v", c, gotX))
		assert.Equal(t, gotY, c.ny, fmt.Sprintf("Failed case %+v, got %+v", c, gotY))
	}
}

func TestValidatePos(t *testing.T) {
	w := NewWorld(10, 10, 100, 1_000_000, 25_000, 50_000, 50, 500_000, 900_000, 1_1000_000)

	cases := []struct {
		x, y   float32
		nx, ny float32
	}{
		{0, 0, 0, 0},
		{1010, 1010, 10, 10},
		{-10, -10, 990, 990},
	}
	for _, c := range cases {
		f := NewFirefly(c.x, c.y, 0, 0, 1000000, w)
		gotX, gotY := w.validatePos(f.X, f.Y)
		assert.InDelta(t, gotX, c.nx, 1e-6, fmt.Sprintf("Failed case %+v, got %+v", c, gotX))
		assert.InDelta(t, gotY, c.ny, 1e-6, fmt.Sprintf("Failed case %+v, got %+v", c, gotY))
	}
}

func TestSendBlinkTo(t *testing.T) {
	w := NewWorld(10, 10, 100, 1_000_000, 25_000, 50_000, 50, 500_000, 900_000, 1_1000_000)

	// near the right top corner
	f := NewFirefly(99.5, 99.5, 0, 0, 1000000, w)
	w.SendBlinkTo(f, w.Cells[0][0], 'R')
	assert.Equal(t, 1, len(w.Cells[1][0].blinkQueue),
		"The cell to the right should have received the Firefly on the blinkQueue.")
	w.SendBlinkTo(f, w.Cells[0][0], 'T')
	assert.Equal(t, 1, len(w.Cells[0][1].blinkQueue),
		"The cell to the top should have received the Firefly on the blinkQueue.")

	// near the left bottom corner
	g := NewFirefly(0.5, 0.5, 0, 0, 1000000, w)
	w.SendBlinkTo(g, w.Cells[0][0], 'L')
	assert.Equal(t, 1, len(w.Cells[9][0].blinkQueue),
		"The cell to the left should have received the Firefly on the blinkQueue.")
	w.SendBlinkTo(g, w.Cells[0][0], 'B')
	assert.Equal(t, 1, len(w.Cells[0][9].blinkQueue),
		"The cell to the bottom should have received the Firefly on the blinkQueue.")
}

func TestSendBlinkToIdle(t *testing.T) {
	w := NewWorld(3, 3, 100, 1_000_000, 25_000, 50_000, 50, 500_000, 900_000, 1_1000_000)

	// f1 will blink immediately (in cell 2)
	f1 := NewFirefly(201, 150, 0, 0, 1000000, w)
	f1.SetNextBlink(w.Clock - 1)
	// f2 will blink when nudged by f1 (in cell 1)
	f2 := NewFirefly(199, 151, 0, 1, 1000000, w)
	f2.SetNextBlink(w.Clock + 1)

	w.DoStep <- 'S'
	<-w.DoneStep

	assert.Equal(t, false, f1.nudgeable,
		"Firefly 1 should have blinked.")
	assert.Equal(t, false, f2.nudgeable,
		"Firefly 2 should have blinked.")

}

func TestClockTick(t *testing.T) {
	w := NewWorld(4, 4, 50, 1_000_000, 25_000, 50_000, 50, 500_000, 900_000, 1_1000_000)
	w.HatchFireflies(10)
	s := time.Now()
	w.ClockTick()
	fmt.Printf("time.Since(s) = %+v\n", time.Since(s))
}

// Test the computed Manhattan distances on a toro.
func TestManhattanDist(t *testing.T) {
	w := NewWorld(10, 10, 100, 1_000_000, 25_000, 50_000, 50, 500_000, 900_000, 1_1000_000)
	cases := []struct {
		f, g *Firefly
		want float32
	}{
		{
			NewFirefly(99.5, 99.5, 0, 0, 1000000, w),
			NewFirefly(99.5, 99.5, 0, 0, 1000000, w),
			0,
		},
		{
			NewFirefly(50, 50, 0, 0, 1000000, w),
			NewFirefly(50, 950, 0, 0, 1000000, w),
			100,
		},
		{
			NewFirefly(50, 850, 0, 0, 1000000, w),
			NewFirefly(50, 950, 0, 0, 1000000, w),
			100,
		},
		{
			NewFirefly(50, 50, 0, 0, 1000000, w),
			NewFirefly(950, 50, 0, 0, 1000000, w),
			100,
		},
		{
			NewFirefly(50, 50, 0, 0, 1000000, w),
			NewFirefly(950, 950, 0, 0, 1000000, w),
			200,
		},
		{
			NewFirefly(50, 50, 0, 0, 1000000, w),
			NewFirefly(150, 150, 0, 0, 1000000, w),
			200,
		},
	}
	for _, c := range cases {
		got := w.ManhattanDist(c.f, c.g)
		assert.InDelta(t, got, c.want, 1e-6, fmt.Sprintf("Failed case %+v, got %+v", c, got))
	}
}

// Check that the fields/verbs used when printing are valid.
func TestStringWorld(t *testing.T) {
	w := NewWorld(10, 10, 100, 1_000_000, 25_000, 50_000, 50, 500_000, 900_000, 1_1000_000)
	w.String()
}
