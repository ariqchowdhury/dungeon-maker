package dungeon

import (
	"fmt"
	"container/list"
)

func int_abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type BoundingBox struct {
	x, y int
	half_width int
}

func (bb BoundingBox) contains_point(x, y int) bool {
	var check_x bool = int_abs(x - bb.x) < bb.half_width
	var check_y bool = int_abs(y - bb.y) < bb.half_width

	return check_x && check_y
}

func box_does_intersect(b1, b2 BoundingBox) bool {
	var check_x bool = int_abs(b1.x - b2.x) < (b1.half_width + b2.half_width)
	var check_y bool = int_abs(b1.y - b2.y) < (b1.half_width + b2.half_width)

	return check_x && check_y
}

type CellQuadTree struct {
	bounding_box BoundingBox

	cells *list.List
	num_cells int

	nw *CellQuadTree
	ne *CellQuadTree
	sw *CellQuadTree
	se *CellQuadTree
}

func (root *CellQuadTree) init (num_cells int) {
	root.cells = list.New()
	root.num_cells = num_cells
}

// Insert a Cell into the tree. The position will depend on 
// how 'close' it is to Cells in the root
func (root *CellQuadTree) insert (c Cell) bool {

	if !root.bounding_box.contains_point(c.x, c.y) {
		return false
	}

	if root.cells.Len() < root.num_cells {
		root.cells.PushBack(c)
		return true
	}

	// If no room, we subdivide and add point to whichever node will accept it
	if root.nw == nil {
		root.subdivide()
	}

	if (root.nw.insert(c)) {
		return true
	}
	if (root.ne.insert(c)) {
		return true
	}
	if (root.sw.insert(c)) {
		return true
	}
	if (root.se.insert(c)) {
		return true
	}

	return false
}

// This function adds 4 child quad trees with appropriate bounding boxes to
// root quad tree
func (root *CellQuadTree) subdivide () {
	var new_half_width int = root.bounding_box.half_width / 2

	var nw_center_x, nw_center_y int = root.bounding_box.x - new_half_width, root.bounding_box.y + new_half_width
	var ne_center_x, ne_center_y int = root.bounding_box.x + new_half_width, root.bounding_box.y + new_half_width
	var sw_center_x, sw_center_y int = root.bounding_box.x - new_half_width, root.bounding_box.y - new_half_width
	var se_center_x, se_center_y int = root.bounding_box.x + new_half_width, root.bounding_box.y - new_half_width

	nw_bb := BoundingBox{nw_center_x, nw_center_y, new_half_width}
	ne_bb := BoundingBox{ne_center_x, ne_center_y, new_half_width}
	sw_bb := BoundingBox{sw_center_x, sw_center_y, new_half_width}
	se_bb := BoundingBox{se_center_x, se_center_y, new_half_width}

	var nw = &CellQuadTree{bounding_box: nw_bb}
	var ne = &CellQuadTree{bounding_box: ne_bb}
	var sw = &CellQuadTree{bounding_box: sw_bb}
	var se = &CellQuadTree{bounding_box: se_bb}

	nw.init(4)
	ne.init(4)
	sw.init(4)
	se.init(4)

	root.nw = nw
	root.ne = ne
	root.sw = sw
	root.se = se

}

func (root *CellQuadTree) check_range (target_range BoundingBox) *list.List {
	// List of results
	cells_in_range := list.New()

	// if the range doesn't intersect root's bounding box, abort
	if !box_does_intersect(root.bounding_box, target_range) {
		return cells_in_range
	}

	// Check this level of quad for cells in target range
	for itr := root.cells.Front(); itr != nil; itr = itr.Next() {
		if target_range.contains_point(itr.Value.(Cell).x, itr.Value.(Cell).y) {
			cells_in_range.PushBack(itr.Value)
		}
	}

	if root.nw == nil {
		return cells_in_range
	}

	cells_in_range.PushBackList(root.nw.check_range(target_range))
	cells_in_range.PushBackList(root.ne.check_range(target_range))
	cells_in_range.PushBackList(root.sw.check_range(target_range))
	cells_in_range.PushBackList(root.se.check_range(target_range))

	return cells_in_range
}

func (root *CellQuadTree) print () {
	
	for itr := root.cells.Front(); itr != nil; itr = itr.Next() { 
		fmt.Println(itr.Value)
	}
	fmt.Println("\n")

	if root.nw != nil {
		fmt.Println("nw:", "\n")
		root.nw.print()
		fmt.Println("ne:", "\n")
		root.ne.print()
		fmt.Println("sw:", "\n")
		root.sw.print()
		fmt.Println("se:", "\n")
		root.se.print()
	}

}