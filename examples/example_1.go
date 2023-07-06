package examples

import (
	"fmt"
	"math"
)

// Point represents a 2D point
type Point struct {
	X, Y int
}

func NewPoint(x, y int) *Point {
	return &Point{X: x, Y: y}
}

func (p *Point) Distance(q *Point) float64 {
	dx := p.X - q.X
	dy := p.Y - q.Y
	return math.Sqrt(float64(dx*dx + dy*dy))
}

func main() {
	p := NewPoint(0, 0)
	q := NewPoint(3, 4)
	distance := p.Distance(q)
	fmt.Println("Distance:", distance)
}
