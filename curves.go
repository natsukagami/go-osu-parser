package parser

import (
	"math"
	"sort"
)

/**
 * Taken from Osu-Web with some fixes
 * https://github.com/pictuga/osu-web
 */

func isPointInCircle(point, center Point, radius float64) bool {
	return distancePoints(point, center) <= radius
}

func distancePoints(p1, p2 Point) float64 {
	var (
		x = p1.X - p2.X
		y = p1.Y - p2.Y
	)
	return math.Sqrt(x*x + y*y)
}

func distanceFromPoints(array []Point) float64 {
	distance := 0.0
	for i := 1; i < len(array)-1; i++ {
		distance += distancePoints(array[i-1], array[i])
	}
	return distance
}

func angleFromPoints(p1, p2 Point) float64 {
	return math.Atan((p2.Y - p1.Y) / (p2.X - p1.X))
}

func cartFromPol(r, teta float64) (float64, float64) {
	return r * math.Cos(teta), r * math.Sin(teta)
}

func pointAtDistanceArray(arr []Point, dist float64) (Point, float64, int) {
	if len(arr) < 2 {
		return Point{}, 0, 0
	}
	d := 0.0
	var i int
	for i = 0; i+1 < len(arr); i++ {
		nd := distancePoints(arr[i], arr[i+1])
		if d+nd >= dist {
			break
		}
		d += nd
	}
	if i == len(arr)-1 {
		return arr[i], angleFromPoints(arr[i-1], arr[i]), i - 1
	}
	angle := angleFromPoints(arr[i], arr[i+1])
	x, y := cartFromPol(dist-d, angle)
	if arr[i].X > arr[i+1].X {
		return Point{arr[i].X - x, arr[i].Y - y}, angle, i
	}
	return Point{arr[i].X + x, arr[i].Y + y}, angle, i
}

// todo: Will convert to int if ans needed isn't large?
func pCn(p, n int) float64 {
	if p < 0 || p > n {
		return 0
	}
	if n-p < p {
		p = n - p
	}
	out := 1.0
	for i := 1.0; i < float64(p)+1; i += 1.0 {
		out = out * (float64(n-p) + i) / i
	}
	return out
}

type floatSort []float64

func (f floatSort) Len() int           { return len(f) }
func (f floatSort) Less(i, j int) bool { return f[i] < f[j] }
func (f floatSort) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

func arrayValues(mp map[float64]Point) []Point {
	arr := make([]Point, len(mp))
	key := make(floatSort, 0)
	for k := range mp {
		key = append(key, k)
	}
	sort.Sort(key)
	for id, val := range key {
		arr[id] = mp[val]
	}
	return arr
}

func calc(op float64, x, y Point) Point {
	return Point{x.X + y.X*op, x.Y + y.Y*op}
}

/*************************************************************/

type pointSystem interface {
	Points() []Point
	Order() int
	CalcPoints()
	Pos() map[float64]Point
}

func pointAtDistance(b pointSystem, dist float64) *Point {
	var p Point
	switch b.Order() {
	case 0:
		return nil
	case 1:
		p = b.Points()[0]
	default:
		b.CalcPoints()
		p, _, _ = pointAtDistanceArray(arrayValues(b.Pos()), dist)
	}
	return &p
}

type bezier struct {
	points   []Point
	order    int
	Step     float64
	pos      map[float64]Point
	Pxlength float64
}

func (b *bezier) Points() []Point        { return b.points }
func (b *bezier) Order() int             { return b.order }
func (b *bezier) Pos() map[float64]Point { return b.pos }

func newBezier(points []Point) bezier {
	b := bezier{
		points: points,
		order:  len(points),
		Step:   0.0025 / float64(len(points)),
		pos:    make(map[float64]Point)}
	b.CalcPoints()
	return b
}

func (b bezier) at(t float64) Point {
	if ans, ok := b.pos[t]; ok {
		return ans
	}
	p := Point{}
	n := b.order - 1
	for i := 0; i <= n; i++ {
		p.X += pCn(i, n) * math.Pow(1.0-t, float64(n-i)) * math.Pow(t, float64(i)) * b.points[i].X
		p.Y += pCn(i, n) * math.Pow(1.0-t, float64(n-i)) * math.Pow(t, float64(i)) * b.points[i].Y
	}
	b.pos[t] = p
	return p
}

func (b *bezier) CalcPoints() {
	if len(b.pos) > 0 {
		return
	}
	b.Pxlength = 0.0
	var (
		prev    = b.at(0)
		current Point
	)
	for i := 0.0; i < 1+b.Step; i += b.Step {
		current = b.at(i)
		b.Pxlength += distancePoints(prev, current)
		prev = current
	}
}

/*************************************************************/

type catmull struct {
	points   []Point
	order    int
	Step     float64
	pos      map[float64]Point
	Pxlength float64
}

func (b *catmull) Points() []Point        { return b.points }
func (b *catmull) Order() int             { return b.order }
func (b *catmull) Pos() map[float64]Point { return b.pos }

func (b catmull) at(x int, t float64) Point {
	var (
		v1 = b.points[x]
		v2 = b.points[x]
	)
	if x >= 1 {
		v1 = b.points[x-1]
	}
	var v3 = calc(1, v2, calc(-1, v2, v1))
	if x+1 < len(b.points) {
		v3 = b.points[x+1]
	}
	var v4 = calc(1, v3, calc(-1, v3, v2))
	if x+2 < len(b.points) {
		v4 = b.points[x+2]
	}
	ret := Point{
		(-v1.X+3*v2.X-3*v3.X+v4.X)*(t*t*t) + (2*v1.X-5*v2.X+4*v3.X-v4.X)*t*t + (-v1.X+v3.X)*t + 2*v2.X,
		(-v1.Y+3*v2.Y-3*v3.Y+v4.Y)*(t*t*t) + (2*v1.Y-5*v2.Y+4*v3.Y-v4.Y)*t*t + (-v1.Y+v3.Y)*t + 2*v2.Y,
	}
	return ret
}

func (b *catmull) CalcPoints() {
	if len(b.pos) > 0 {
		return
	}
	for i := 0; i < b.order-1; i++ {
		for t := 0.0; t < 1+b.Step; t += b.Step {
			b.pos[float64(len(b.pos))] = b.at(i, t)
		}
	}
}

func newCatmull(points []Point) catmull {
	c := catmull{
		points: points,
		order:  len(points),
		Step:   0.025,
		pos:    make(map[float64]Point),
	}
	c.CalcPoints()
	return c
}

func init() {
	var (
		_ pointSystem = (*catmull)(nil)
		_ pointSystem = (*bezier)(nil)
	)
}
