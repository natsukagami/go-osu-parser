package parser

import "strconv"

// Point represents a point on the screen
type Point struct {
	X float64
	Y float64
}

func parsePoint(x, y string) (p Point, err error) {
	if p.X, err = strconv.ParseFloat(x, 64); err != nil {
		return
	}
	if p.Y, err = strconv.ParseFloat(y, 64); err != nil {
		return
	}
	return
}
