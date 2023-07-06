package examples

import (
	"fmt"
	"math"
)

// Point represents a 2D point
type Point struct {
	X, Y int
}

// NewPoint creates a new Point object with the given x and y coordinates.
// Example usage: p := NewPoint(3, 5)
func NewPoint(x, y int) *Point {
	return &Point{X: x, Y: y}
}

// Distance calculates the euclidean distance between two points in 2D space.
// Example:
// p := &Point{X: 1.0, Y: 2.0}
// q := &Point{X: 3.0, Y: 4.0}
// dist := p.Distance(q)
func (p *Point) Distance(q *Point) float64 {
	dx := p.X - q.X
	dy := p.Y - q.Y
	return math.Sqrt(float64(dx*dx + dy*dy))
}

// This function calculates the distance between two points using the Pythagorean theorem.
// Example:
// p := NewPoint(0, 0)
// q := NewPoint(3, 4)
// distance := p.Distance(q)
// fmt.Println("Distance:", distance)
func main() {
	p := NewPoint(0, 0)
	q := NewPoint(3, 4)
	distance := p.Distance(q)
	fmt.Println("Distance:", distance)
}
