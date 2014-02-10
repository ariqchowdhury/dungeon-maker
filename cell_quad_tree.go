package dungeon

import "fmt"

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
	return int_abs(x - bb.x) < bb.half_width && int_abs(y - bb.y) < bb.half_width
}

func box_does_intersect(b1, b2 BoundingBox) bool {
	return int_abs(b1.x - b2.x) < (b1.half_width + b2.half_width) && int_abs(b1.y - b2.y) < (b1.half_width + b2.half_width)
}

type CellQuadTree struct {
	bounding_box BoundingBox

	cells []Cell
	num_cells int

	nw *CellQuadTree
	ne *CellQuadTree
	sw *CellQuadTree
	se *CellQuadTree
}

func (root *CellQuadTree) init (num_cells int) {
	c := make([]Cell, 0, num_cells)
	root.cells = c
	root.num_cells = num_cells
}

// Insert a Cell into the tree. The position will depend on 
// how 'close' it is to Cells in the root
func (root *CellQuadTree) insert (c Cell) bool {

	if root.bounding_box.contains_point(c.x, c.y) {
		return false
	}

	if len(root.cells) < root.num_cells {
		root.cells = append(root.cells, c)
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

func (root *CellQuadTree) check_range (c Cell) {

}

func (root *CellQuadTree) print () {
	fmt.Println(root.cells, "\n")
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