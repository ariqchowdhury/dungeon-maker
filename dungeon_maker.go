package dungeon

import (
	"fmt"
	"math/rand"
	"time"
	"bufio"
	"io"
)

type Dungeon struct {
	cells []Cell
}

// Creates a num_cells number of random sized cells
// Cell sizes use normal distribution, so use std_dev and mean
// to control size
func (d *Dungeon) CreateCells(num_cells int, std_dev, mean float64) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	radius := int(r.NormFloat64() * std_dev + mean)
	c := make([]Cell, num_cells)

	for i := range c {
		radius = int(r.NormFloat64() * std_dev + mean)
		c[i].radius = radius
		c[i].id = i
	}

	d.cells = c
}

// Place x,y coordinates of cells randomly, with normal distribution.
// Use std_dev and mean to control size of overall map
func (d *Dungeon) PlaceCells(std_dev, mean float64) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range d.cells {
		x := int(r.NormFloat64() * std_dev + mean)
		y := int(r.NormFloat64() * std_dev + mean)
		d.cells[i].x, d.cells[i].y = x, y
	}
}

// pushes cells apart so that they don't overlap
func (d *Dungeon) SeperateCells() {
	var cells_to_move []CellPair

	// find cells that are probably close enough to intersect cell
	// based on likely radius (max would be 2 * max radius)

	for i, c := range d.cells {
		for _, cc := range d.cells[i+1:len(d.cells)] {
			if does_intersect(c, cc) {
				intersector := CellPair{c,cc}
				cells_to_move = append(cells_to_move, intersector)
			}
		}
	}

	fmt.Println(cells_to_move)
}

func (d *Dungeon) WriteCells(w io.Writer) {
	writer := bufio.NewWriter(w)

	for i := range d.cells {
		fmt.Fprint(writer, d.cells[i])
		fmt.Fprint(writer, "\n")
	}
	writer.Flush()
}

// Create a grid and place the Dungeon's list of cells randomly about
func MakeDungeon() {

}

// Create the connections between cells on the grid
func MakeCorriders() {

}