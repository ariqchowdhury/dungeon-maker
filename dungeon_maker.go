package dungeon

import (
	"fmt"
	"math/rand"
	"time"
	"bufio"
	"io"
	"container/list"
	"math"
)

type Dungeon struct {
	cells []Cell
	boundary_half_dimension int
	cell_quad_tree CellQuadTree
	target_range_bb_half_width int
	rooms []Cell
	corridors []Cell
}

func (d *Dungeon) SetBoundaryHalfDimension(dim int) {
	d.boundary_half_dimension = dim
}

// Creates a num_cells number of random sized cells
// Cell sizes use normal distribution, so use std_dev and mean
// to control size
func (d *Dungeon) CreateCells(num_cells int, std_dev, mean float64) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	c := make([]Cell, num_cells)

	for i := range c {
		radius := r.NormFloat64() * std_dev
		radius = math.Abs(radius)
		radius *= 20.0
		c[i].radius = int(radius)
		c[i].id = i
	}

	d.cells = c

	// Set the target range bounding box here, once we know likely max radius
	// TODO: if we are using a skewed distribution of radii such that there are more
	// small than a few large, then this target range is over kill and will slow down performance
	d.target_range_bb_half_width = 35
}

// Place x,y coordinates of cells randomly, with normal distribution.
// Use std_dev and mean to control size of overall map
func (d *Dungeon) PlaceCells() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	std_dev := float64(d.boundary_half_dimension / 5)
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
		d.cell_quad_tree.insert(cell)
	}
}

func (d *Dungeon) PrintCellsQuadTree() {
	fmt.Println("CellQuadTree:")
	d.cell_quad_tree.print()
}

// pushes cells apart so that they don't overlap
func (d *Dungeon) SeperateCells() {	
	// do some number of iterations of checking, and then moving away
	var all_seperated bool = false
	var max_itr int = 0

	for !all_seperated && max_itr < 200 {
		all_seperated = true
		for i, cell := range d.cells {
			// Check if a Cell has fellow cells in our target range bounding box
			// target range is a square with half width of 2*max radius (estimate)
			target_range := BoundingBox{cell.x, cell.y, d.target_range_bb_half_width}
			l := d.cell_quad_tree.check_range(target_range)

			cells_that_collide := list.New()
			// for the list of potential targets, check for collision
			for iter := l.Front(); iter != nil; iter = iter.Next() {
				if does_intersect(iter.Value.(Cell), cell) {
					cells_that_collide.PushBack(iter.Value.(Cell))
				}
			}
			// Remove the first item in list because that is the cell itself
			if cells_that_collide.Len() > 0 {
				cells_that_collide.Remove(cells_that_collide.Front())
			}

			if cells_that_collide.Len() != 0 {
				all_seperated = all_seperated && false
			}

			// for the list of collisions, increase position of current cell;
			// do this by measuring distance between points, and add a vector in a
			// direction proportional to how close they are -- if they collide, but are far
			// away move small 
			for iter := cells_that_collide.Front(); iter != nil; iter = iter.Next() {
				delta_x, delta_y := cell_distance_xy_components(iter.Value.(Cell), cell)

				d.cells[i].x += (delta_x / 2)
				d.cells[i].y += (delta_y / 2)
			}
		}
		d.PlaceCellsQuadTree()
		max_itr++
	}

	// There may be some problematic remaining cells that intersect others; remove
	// them
	// for i, cell := range d.cells {
	// 	target_range := BoundingBox{cell.x, cell.y, d.target_range_bb_half_width}
	// 	l := d.cell_quad_tree.check_range(target_range)

	// 	cells_that_collide := list.New()
	// 	// for the list of potential targets, check for collision
	// 	for iter := l.Front(); iter != nil; iter = iter.Next() {
	// 		if does_intersect(iter.Value.(Cell), cell) {
	// 			cells_that_collide.PushBack(iter.Value.(Cell))
	// 		}
	// 	}
	// 	// Remove the first item in list because that is the cell itself
	// 	if cells_that_collide.Len() > 0 {
	// 		cells_that_collide.Remove(cells_that_collide.Front())
	// 	}

	// 	if 
	// }
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