package firefly

import (
	"fmt"
	"math/rand"
	"sync"
)

// World represents the whole environment.
type World struct {
	Cells     [][]*Cell // Cells in the world.
	CellWNum  int       // Width of the world in cells.
	CellHNum  int       // Height of the world in cells.
	CellSize  float32   // Size of the cells in pixels.
	SizeW     float32   // Width of the world in pixels.
	SizeH     float32   // Height of the world in pixels.
	sizeHalfW float32   // Half the width of the world in pixels.
	sizeHalfH float32   // Half the height of the world in pixels.

	Clock         int            // Internal time of the simulation, in us.
	ClockTickLen  int            // Update per tick.
	wgClockTick   sync.WaitGroup // WG to sync the blinking.
	NudgeAmount   int            // How much to nudge the firefly deadlines.
	NudgeRadius   float32        // Max distance between communicating fireflies.
	borderDist    float32        // Distance from a border to require a blinkQueue to the neighbor.
	BlinkCooldown int            // Cooldown after blinking while the Firefly is not nudgeable.
	PeriodMin     int            // Minimum length of the fireflies' period.
	PeriodMax     int            // Maximum length of the fireflies' period.

	chChangeCell     chan *ChangeCellReq   // A firefly needs to enter/leave the cell.
	chChangeCellDone chan bool             // The cell change is done.
	chChangeCells    chan []*ChangeCellReq // Channel for many fireflies to enter/leave the cell.

	DoStep   chan byte      // Channel to request a step of the env.
	DoneStep chan bool      // Channel to signal the end of a step of the env.
	wgMove   sync.WaitGroup // WG to sync the fireflies movement.
}

// NewWorld creates a new World.
func NewWorld(
	cw, ch int,
	cellSize float32,
	clockStart, clockTickLen int,
	nudgeAmount int,
	nudgeRadius float32,
	blinkCooldown int,
	periodMin, periodMax int,
) *World {

	cacheCosSin()

	w := &World{}

	// dimensions params
	w.CellSize = cellSize
	w.CellWNum = cw
	w.CellHNum = ch
	w.SizeW = float32(cw) * cellSize
	w.SizeH = float32(ch) * cellSize
	w.sizeHalfW = w.SizeW / 2
	w.sizeHalfH = w.SizeH / 2

	// nudging params
	w.Clock = clockStart
	w.ClockTickLen = clockTickLen
	w.NudgeAmount = nudgeAmount
	w.NudgeRadius = nudgeRadius
	w.borderDist = w.NudgeRadius / 2
	w.BlinkCooldown = blinkCooldown
	w.PeriodMin = periodMin
	w.PeriodMax = periodMax
	// w.Clock = 1_000_000     // start at 1 second
	// w.ClockTickLen = 25_000 // 25 ms
	// w.ClockTickLen = 1_000 // 1 ms
	// w.NudgeAmount = 50_000 // 50 ms
	// w.NudgeAmount = 100_000 // 50 ms
	// w.NudgeRadius = 20
	// w.NudgeRadius = 50
	// w.NudgeRadius = 50
	// w.NudgeRadius = 100
	// w.BlinkCooldown = 500_000 // 200 ms

	// channels
	w.chChangeCell = make(chan *ChangeCellReq, 100)
	w.chChangeCellDone = make(chan bool)
	w.chChangeCells = make(chan []*ChangeCellReq)
	w.DoStep = make(chan byte)
	w.DoneStep = make(chan bool)

	// create the cells
	c := make([][]*Cell, cw)
	for i := 0; i < cw; i++ {
		c[i] = make([]*Cell, ch)
		for ii := 0; ii < ch; ii++ {
			c[i][ii] = NewCell(w, i, ii)
		}
	}
	w.Cells = c

	// start listening
	go w.Listen()

	return w
}

// HatchFireflies creates a swarm of fireflies.
func (w *World) HatchFireflies(n int) {
	w.HatchFirefliesFromID(n, 0)
}

// HatchFireflies creates a swarm of fireflies, with IDs starting from idStart.
func (w *World) HatchFirefliesFromID(n, idStart int) {
	for i := idStart; i < n+idStart; i++ {
		// random pos/ori/period
		x := rand.Float32() * w.SizeW
		y := rand.Float32() * w.SizeH
		o := int16(rand.Float64() * 360)
		p := RandRangeInt(w.PeriodMin, w.PeriodMax)
		NewFirefly(x, y, o, i, p, w)
	}
}

// Listen to all the channels to react.
func (w *World) Listen() {
	for {
		select {

		// change cell of a firefly
		case r := <-w.chChangeCell:
			w.ChangeCell(r)
			w.chChangeCellDone <- true

		// tick forward the env
		case <-w.DoStep:
			w.Step()
			w.DoneStep <- true
		}
	}
}

// Perform a step of the simulation: move the fireflies and advance the clock.
func (w *World) Step() {
	w.Move()
	w.ClockTick()
}

// Perform a movement of the fireflies.
func (w *World) Move() {

	// move all the fireflies
	for i := 0; i < w.CellWNum; i++ {
		for ii := 0; ii < w.CellHNum; ii++ {
			w.wgMove.Add(1)
			w.Cells[i][ii].chMove <- 'M'
		}
	}

	// wait for the wg to be done
	// so that all the fireflies are done moving
	// and no cell is still iterating on c.Fireflies
	w.wgMove.Wait()

	// perform all the cell change
	for i := 0; i < w.CellWNum; i++ {
		for ii := 0; ii < w.CellHNum; ii++ {
			reqs := <-w.chChangeCells
			for _, r := range reqs {
				w.ChangeCell(r)
			}
		}
	}

}

// Perform a clock tick and blink the fireflies.
func (w *World) ClockTick() {
	w.Clock += w.ClockTickLen

	// reset all the cells to working
	for i := 0; i < w.CellWNum; i++ {
		for ii := 0; ii < w.CellHNum; ii++ {
			w.Cells[i][ii].idle = false
		}
	}

	// blink the fireflies in each cell
	// wait for all the cells to be done simultaneously
	for i := 0; i < w.CellWNum; i++ {
		for ii := 0; ii < w.CellHNum; ii++ {
			w.wgClockTick.Add(1)
			w.Cells[i][ii].chBlink <- 'B'
		}
	}
	w.wgClockTick.Wait()

	// send a signal to all cells to quit blinking
	for i := 0; i < w.CellWNum; i++ {
		for ii := 0; ii < w.CellHNum; ii++ {
			w.Cells[i][ii].blinkDone <- true
		}
	}
}

// ChangeCell moves a firefly from a cell to another.
func (w *World) ChangeCell(r *ChangeCellReq) {
	// update the cells
	if r.from != nil {
		r.from.Leave(r.f)
	}
	r.to.Enter(r.f)
	// update the info inside the firefly
	r.f.c = r.to
}

// EnterCell moves a firefly into a cell.
//
// Will block until the change has been completed: Listen must be already be active.
func (w *World) EnterCell(f *Firefly, c *Cell) {
	w.chChangeCell <- &ChangeCellReq{f, nil, c}
	<-w.chChangeCellDone
}

// Move by (dcx, dcy) around the cells' toro, from cell (cx, cy).
func (w *World) MoveWrapCell(cx, cy, dcx, dcy int) (int, int) {
	cx += dcx
	for cx < 0 {
		cx += w.CellWNum
	}
	for cx >= w.CellWNum {
		cx -= w.CellWNum
	}
	cy += dcy
	for cy < 0 {
		cy += w.CellHNum
	}
	for cy >= w.CellHNum {
		cy -= w.CellHNum
	}
	return cx, cy
}

// Send a blink to the requested neighbor.
func (w *World) SendBlinkTo(f *Firefly, c *Cell, dir byte) {

	dx, dy := 0, 0
	switch dir {
	case 'L':
		dx = -1
	case 'R':
		dx = 1
	case 'B':
		dy = -1
	case 'T':
		dy = 1
	}

	// find the neighboring cell on the toro
	ncx, ncy := f.w.MoveWrapCell(f.c.Cx, f.c.Cy, dx, dy)
	nc := w.Cells[ncx][ncy]

	// check if ncx,ncy was idling
	// if so, set idle to false and Add(1) on the WaitGroup counter
	nc.idleLock.Lock()
	nc.blinkQueue <- f
	if nc.idle {
		nc.w.wgClockTick.Add(1)
		nc.idle = false
	}
	nc.idleLock.Unlock()
}

// Compute the Manhattan distance on a torus between two fireflies.
func (w *World) ManhattanDist(f, g *Firefly) float32 {

	// if the two are further apart than the SizeHalf
	// the shorter distance is by going around the toro
	ax := AbsFloat32(f.X - g.X)
	if ax > w.sizeHalfW {
		ax = w.SizeW - ax
	}
	ay := AbsFloat32(f.Y - g.Y)
	if ay > w.sizeHalfH {
		ay = w.SizeH - ay
	}

	return ax + ay
}

// Ensure that the coordinates provided are a valid world position.
func (w *World) validatePos(x, y float32) (float32, float32) {
	for x < 0 {
		x += w.SizeW
	}
	for x >= w.SizeW {
		x -= w.SizeW
	}
	for y < 0 {
		y += w.SizeH
	}
	for y >= w.SizeH {
		y -= w.SizeH
	}
	return x, y
}

// String implements fmt.Stringer.
func (w *World) String() string {
	s := fmt.Sprintf("W: %dx%d (%.2f) %.2fx%.2f",
		w.CellWNum, w.CellHNum,
		w.CellSize,
		w.SizeW, w.SizeH,
	)
	for i := 0; i < w.CellWNum; i++ {
		for ii := 0; ii < w.CellHNum; ii++ {
			// Add the state of the cell to the World repr.
			s += fmt.Sprintf("\nC: %v",
				w.Cells[i][ii],
			)
		}
	}
	return s
}
