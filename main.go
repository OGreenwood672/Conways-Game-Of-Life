package main

import (
	"fmt"
	"time"
	"github.com/veandco/go-sdl2/sdl"  //https://godoc.org/github.com/veandco/go-sdl2/sdl
)

const (
	WIDTH = 800
	HEIGHT = 800
	BOXSIZE = 50
)

var MOUSESTATE uint32

type Cell struct {
	x int
	y int
	state bool
	FutureState bool
}

func (c Cell) GetNeighbours(grid [BOXSIZE][BOXSIZE]Cell) int {
	var neighbours int
	options := [8][2]int{{0, 1}, {1, 0}, {1, 1}, {-1, 1}, {1, -1}, {-1, -1}, {-1, 0}, {0, -1}}
	for i:=0; i != len(options); i++ {
		y := c.y + options[i][1]
		x := c.x + options[i][0]
		if 0 <= y && y < len(grid) && 0 <= x && x < len(grid) {
			cell := grid[y][x]
			if cell.state {
				neighbours++
			}
		}
	}
	return neighbours
}

func setUp() [BOXSIZE][BOXSIZE]Cell {

	defer fmt.Println("**Set Up Complete**")

	fmt.Println("****WELCOME****")
	fmt.Println("**INSTRUCTIONS**\nPress <space> to begin simulation\nPress <backspace> to halt simulation")

	grid := [BOXSIZE][BOXSIZE]Cell{}
	for y, row := range grid {
		for x, _ := range row {
			grid[y][x].x = x
			grid[y][x].y = y
		}
	}
	return grid
}


func draw(renderer *sdl.Renderer, grid [BOXSIZE][BOXSIZE]Cell) {

	for y, row := range grid {
		for x, cell := range row {

			var err error

			if cell.state {
				err = renderer.SetDrawColor(255, 255, 255, 255)
			} else {
				err = renderer.SetDrawColor(0, 0, 0, 255)
			}

			if err != nil {
				fmt.Println(err)
			}

			rect := sdl.Rect{int32(x * BOXSIZE), int32(y * BOXSIZE), BOXSIZE, BOXSIZE}
			
			renderer.FillRect(&rect)
			renderer.DrawRect(&rect)
		}
	}
}

func changeCell(grid *[BOXSIZE][BOXSIZE]Cell) {
	x, y, state := sdl.GetMouseState()
	x = int32(x / BOXSIZE)
	y = int32(y / BOXSIZE)
	if state - MOUSESTATE == 1 {
		grid[y][x].state = !grid[y][x].state
	}
	MOUSESTATE = state
}

func update(grid *[BOXSIZE][BOXSIZE]Cell) bool {

	var change bool

	for y, row := range *grid {
		for x, cell := range row {
			neighbours := cell.GetNeighbours(*grid)
			if (neighbours == 2 || neighbours == 3) && cell.state {
				grid[y][x].FutureState = true
			} else if !cell.state && neighbours == 3 {
				grid[y][x].FutureState = true
				change = true
			} else {
				if grid[y][x].state {
					change = true
				}
				grid[y][x].FutureState = false
			}
		}
	}
	for y, row := range grid {
		for x, cell := range row {
			grid[y][x].state = cell.FutureState
		}
	}

	return change
}

func checkPress(scancode int) bool {
	keys := sdl.GetKeyboardState()
	if keys[scancode] == 1 { return true }
	return false
}

func finished(grid [BOXSIZE][BOXSIZE]Cell, change bool) bool {

	if checkPress(sdl.SCANCODE_BACKSPACE) {
		fmt.Println("****Civilization Halted****")
		return false
	}

	if !change {
		fmt.Println("****Optimal civilization reached****")
		return false
	}

	for _, row := range grid {
		for _, cell := range row {
			if cell.state {
				return true
			}
		}
	}
	fmt.Println("****You have gone extinct!!****")
	
	return false
}

func main() {
	
	grid := setUp()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("Unable to initilize sdl:", err)
		return
	}

	window, err := sdl.CreateWindow(
		"Conway's Game of Life",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WIDTH, HEIGHT,
		sdl.WINDOW_OPENGL,
	)
	if err != nil {
		fmt.Println("Unable to initilize window:", err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("Unable to initilize renderer:", err)
		return
	}
	defer renderer.Destroy()

	running := false
	var change bool

	ticker := time.NewTicker(time.Second / 5)

	for {

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		renderer.Clear()

		draw(renderer, grid)

		if !running {
			changeCell(&grid)
			running = checkPress(sdl.SCANCODE_SPACE)
		} else {
			change = update(&grid)
			running = finished(grid, change)
			<-ticker.C
		}

		renderer.Present()
	}

}