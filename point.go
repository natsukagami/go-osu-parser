package parser

import (
	"fmt"
	"strconv"
)

// Point represents a point on the screen
type Point struct {
	X float64
	Y float64
}

// MarshalJSON parses Point as a 2-member array in JSON
func (p Point) MarshalJSON() (b []byte, err error) {
	b = []byte(fmt.Sprintf("[%v,%v]", p.X, p.Y))
	return
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
