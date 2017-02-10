package parser

import "math"

func getSliderEndPoint(sliderType string, sliderLength float64, points []Point) Point {
	if len(points) < 2 {
		return Point{} // Wtf slider with less than 2 points?
	}
	if len(points) == 2 {
		return pointOnLine(points[0], points[1], sliderLength)
	}
	switch sliderType {
	default:
		return Point{}
	case "linear":
		return pointOnLine(points[0], points[1], sliderLength)
	case "catmull":
		return Point{} // unsupported
	case "bezier":
		pts := make([]Point, len(points))
		copy(pts, points)
		var (
			b           bezier
			prev, point *Point
		)
		for i, l := 0, len(pts); i < l; i++ {
			point = &pts[i]
			if prev == nil {
				prev = point
				continue
			}
			if *point == *prev {
				b = newBezier(pts[0:i])
				pts = pts[i:]
				sliderLength -= b.Pxlength
				i = 0
				l = len(pts)
			}
			prev = point
		}
		b = newBezier(pts)
		px := pointAtDistance(&b, sliderLength)
		if px == nil {
			return Point{}
		}
		return *px
	case "pass-through":
		if len(points) > 3 {
			return getSliderEndPoint("bezier", sliderLength, points)
		}
		var (
			p1 = points[0]
			p2 = points[1]
			p3 = points[2]
		)
		circumCircle := getCircumCircle(p1, p2, p3)
		radians := sliderLength / circumCircle.R
		if isLeft(p1, p2, p3) {
			radians *= -1
		}
		return rotate(circumCircle.C, p1, radians)
	}
}

func pointOnLine(p1, p2 Point, length float64) Point {
	f := distancePoints(p1, p2)
	n := f - length
	return Point{
		(n*p1.X + length*p2.X) / f,
		(n*p1.Y + length*p2.Y) / f,
	}
}

// Get coordinates of a point in a circle, given the center, a startpoint and a distance in radians
func rotate(c Point, p Point, r float64) Point {
	var (
		cos = math.Cos(r)
		sin = math.Sin(r)
	)
	return Point{
		(cos * (p.X - c.X)) - (sin * (p.Y - c.Y)) + c.X,
		(sin * (p.X - c.X)) + (cos * (p.Y - c.Y)) + c.Y,
	}
}

// Checks if C is on left side of [AB]
func isLeft(a, b, c Point) bool {
	return ((b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X)) < 0
}

type circle struct {
	C Point   // Circle center
	R float64 // Circle radius
}

// Squares a number
func sq(x float64) float64 { return x * x }

// Gets circle circumventing 3 points
func getCircumCircle(a, b, c Point) circle {
	D := 2 * (a.X*(b.Y-c.Y) + b.X*(c.Y-a.Y) + c.X*(a.Y-b.Y))
	var (
		Ux = ((sq(a.X)+sq(a.Y))*(b.Y-c.Y) + (sq(b.X)+sq(b.Y))*(c.Y-a.Y) + (sq(c.X)+sq(c.Y))*(a.Y-b.Y)) / D
		Uy = ((sq(a.X)+sq(a.Y))*(c.X-b.X) + (sq(b.X)+sq(b.Y))*(a.X-c.X) + (sq(c.X)+sq(c.Y))*(b.X-a.X)) / D
	)
	return circle{Point{Ux, Uy}, math.Sqrt(sq(Ux-a.X) + sq(Uy-a.Y))}
}
