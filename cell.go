package dungeon

import (
	"math"
)	

type Cell struct {
	id int
	x, y, radius float64
}

func is_point_in_cell(c Cell, x, y float64) bool {
	return (c.x - x)*(c.x-x) + (c.y - y)*(c.y - y) < c.radius*c.radius
}


func cell_distance_xy_components(c1, c2 Cell) (float64, float64) {
	return math.Abs(c2.x-c1.x), math.Abs(c2.y-c1.y)
}

func cell_distance(c1, c2 Cell) float64 {

	dx := (c1.x-c2.x)
	dy := (c1.y-c2.y)
	sum_squares := dx*dx + dy*dy
	sqrt := math.Sqrt(sum_squares)

	return sqrt
}

// Checks if circles intersect or if one circle is completely in the other
// So this is more of a collision check than an intersect check
func does_collide(c1, c2 Cell) bool {
	var b1 bool = (c1.radius+c2.radius)*(c1.radius+c2.radius) > 
	(c1.x-c2.x)*(c1.x-c2.x) + (c1.y-c2.y)*(c1.y-c2.y)

	return b1
}


