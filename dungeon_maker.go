package dungeon

import (
	"fmt"
	"math/rand"
	"time"
)

type Dungeon struct {
	grid Grid
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
func (d *Dungeon) seperate_cells() {
	
}

func (d *Dungeon) PrintCells() {
	for i := range d.cells {
		fmt.Println(d.cells[i])
	}
}

// Create a grid and place the Dungeon's list of cells randomly about
func MakeDungeon() Cell {

	cell := Cell{2, 4, 6}

	return cell
}

// Create the connections between cells on the grid
func MakeCorriders() {

}