package dungeon

import (
	"math"
)	

type Cell struct {
	id int
	x, y, radius int
}

func is_point_in_cell(c Cell, x, y int) bool {
	return (c.x - x)*(c.x-x) + (c.y - y)*(c.y - y) < c.radius*c.radius
}


func cell_distance_xy_components(c1, c2 Cell) (int, int) {
	return (c2.x-c1.x), (c2.y-c1.y)
}

func cell_distance(c1, c2 Cell) float64 {
	return math.Sqrt(float64((c1.x-c2.x)*(c1.x-c2.x) + (c1.y-c2.y)*(c1.y*c2.y)))
}

func does_intersect(c1, c2 Cell) bool {
	var b1 bool = (c1.radius-c2.radius)*(c1.radius-c2.radius) <= 
	(c1.x-c2.x)*(c1.x-c2.x) + (c1.y-c2.y)*(c1.y-c2.y)

	var b2 bool = (c1.radius+c2.radius)*(c1.radius+c2.radius) >= 
	(c1.x-c2.x)*(c1.x-c2.x) + (c1.y-c2.y)*(c1.y-c2.y)

	return b1 && b2
}


