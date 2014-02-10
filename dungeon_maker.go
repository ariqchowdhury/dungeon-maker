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
	boundary_half_dimension int
	cell_quad_tree CellQuadTree
}

func (d *Dungeon) SetBoundaryHalfDimension(dim int) {
	d.boundary_half_dimension = dim
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
func (d *Dungeon) PlaceCells() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	std_dev := float64(d.boundary_half_dimension / 4)
	mean := float64(d.boundary_half_dimension)

	for i := range d.cells {
		x := int(r.NormFloat64() * std_dev + mean)
		y := int(r.NormFloat64() * std_dev + mean)
		d.cells[i].x, d.cells[i].y = x, y
	}
}

func (d *Dungeon) PlaceCellsQuadTree() {
	var initial_bb = BoundingBox{d.boundary_half_dimension, d.boundary_half_dimension, d.boundary_half_dimension}
	d.cell_quad_tree = CellQuadTree{bounding_box: initial_bb}
	d.cell_quad_tree.init(4)

	for _, cell := range d.cells {
		if d.cell_quad_tree.insert(cell) == false{
			fmt.Println("Insert Failed")
		}
	}

}

func (d *Dungeon) PrintCellsQuadTree() {
	fmt.Println("CellQuadTree:")
	d.cell_quad_tree.print()
}

// pushes cells apart so that they don't overlap
func (d *Dungeon) SeperateCells() {
	var cells_to_move []CellPair

	// Compare 1 cell at a time with all other cells
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