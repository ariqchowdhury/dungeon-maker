package dungeon

import (
	"fmt"
	"math/rand"
	"time"
	"bufio"
	"io"
	"container/list"
	"os"
	"strings"
	"strconv"
)

type Dungeon struct {
	cells *list.List
	boundary_half_dimension int
	cell_quad_tree CellQuadTree
	target_range_bb_half_width int
	rooms *list.List
	corridors *list.List
}

func (d *Dungeon) SetBoundaryHalfDimension(dim int) {
	d.boundary_half_dimension = dim
}

func (d *Dungeon) CreateCellsFromFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Could not open ", filename)
	}

	d.cells = list.New()

	reader := bufio.NewReader(file)
	line, e := reader.ReadString('\n')

	for e == nil {
		if e != nil {
			fmt.Println(e)
		}
		words := strings.Split(line, " ")
		sid := words[0]
		sx := words[1]
		sy := words[2]
		sradius := words[3]
		sradius = strings.Replace(sradius, "\n", "", -1)

		id, es := strconv.ParseInt(sid, 10, 64)
		x, es := strconv.ParseInt(sx, 10, 64)
		y, es := strconv.ParseInt(sy, 10, 64)
		radius, es := strconv.ParseInt(sradius, 10, 64)

		if es != nil {
			fmt.Println(es)
		}

		cell := &Cell{int(id), int(x), int(y), int(radius)}
		d.cells.PushBack(cell)
		line, e = reader.ReadString('\n')
	}
	d.target_range_bb_half_width = 62
}

// Creates a num_cells number of random sized cells
// Cell sizes use normal distribution, so use std_dev and mean
// to control size
func (d *Dungeon) CreateCells(num_cells int, std_dev, mean float64) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	dim_std_dev := float64(d.boundary_half_dimension / 5)
	dim_mean := float64(d.boundary_half_dimension)

	d.cells = list.New()

	for i := 0; i < num_cells; i++ {
		radius := r.NormFloat64() * std_dev + mean
		x := int(r.NormFloat64() * dim_std_dev + dim_mean)
		y := int(r.NormFloat64() * dim_std_dev + dim_mean)

		cell := &Cell{i, x, y, int(radius)}
		d.cells.PushBack(cell)
	}

	// Set the target range bounding box here, once we know likely max radius
	// TODO: if we are using a skewed distribution of radii such that there are more
	// small than a few large, then this target range is over kill and will slow down performance
	d.target_range_bb_half_width = int((std_dev + mean + 1) * 2.0)
}

func (d *Dungeon) PlaceCellsQuadTree() {
	var initial_bb = BoundingBox{d.boundary_half_dimension, d.boundary_half_dimension, d.boundary_half_dimension}
	d.cell_quad_tree = CellQuadTree{bounding_box: initial_bb}
	d.cell_quad_tree.init(4)

	for itr := d.cells.Front(); itr != nil; itr = itr.Next() {
		d.cell_quad_tree.insert(itr.Value.(*Cell))
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

	for !all_seperated && max_itr < 100 {
		all_seperated = true
		for i := d.cells.Front(); i != nil; i = i.Next() {
			// Check if a Cell has fellow cells in our target range bounding box
			// target range is a square with half width of 2*max radius (estimate)
			target_range := BoundingBox{i.Value.(*Cell).x, i.Value.(*Cell).y, d.target_range_bb_half_width}
			l := d.cell_quad_tree.check_range(target_range)

			cells_that_collide := list.New()
			// for the list of potential targets, check for collision
			for iter := l.Front(); iter != nil; iter = iter.Next() {
				if does_intersect(*iter.Value.(*Cell), *i.Value.(*Cell)) {
					cells_that_collide.PushBack(iter.Value.(*Cell))
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

			for neighbour := cells_that_collide.Front(); neighbour != nil; neighbour = neighbour.Next() {
				// delta_x, delta_y := cell_distance_xy_components(iter.Value.(Cell), *i.Value.(*Cell))

				cell_ptr := neighbour.Value.(*Cell)
				var dx, dy int

				if cell_ptr.x > i.Value.(*Cell).x {
					dx = 1
				} else {
					dx = -1
				}

				if cell_ptr.y > i.Value.(*Cell).y {
					dy = 1
				} else {
					dy = -1
				}

				cell_ptr.x += dx
				cell_ptr.y += dy
				// var delta_y, delta_x int


			}
		}
		d.PlaceCellsQuadTree()
		max_itr++
	}

	for i := d.cells.Front(); i != nil; i = i.Next() {
		cell_ptr := i.Value.(*Cell)
		if cell_ptr.x < 0 {
			cell_ptr.x = 0
		}
		if cell_ptr.y < 0 {
			cell_ptr.y = 0
		}
	}
}

func (d *Dungeon) WriteCells(w io.Writer) {
	writer := bufio.NewWriter(w)

	for itr := d.cells.Front(); itr != nil; itr = itr.Next() {
		fmt.Fprint(writer, itr.Value)
		fmt.Fprint(writer, "\n")
	}

	writer.Flush()
}

// Create a grid and place the Dungeon's list of cells randomly about
func(d *Dungeon) MakeDungeon() {
	for itr := d.cells.Front(); itr != nil; itr = itr.Next() {
		fmt.Println(*itr.Value.(*Cell))
	}

	for itr := d.cell_quad_tree.cells.Front(); itr != nil; itr = itr.Next() {
		cell_ptr := itr.Value.(*Cell)

		cell_ptr.x = 999
		cell_ptr.y = 999
	}

	for itr := d.cells.Front(); itr != nil; itr = itr.Next() {
		fmt.Println(*itr.Value.(*Cell))
	}

}

// Create the connections between cells on the grid
func MakeCorriders() {

}