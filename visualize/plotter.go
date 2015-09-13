package visualize

import (
	"math"
)

// Generating a slice of points from a set of information.
//
//   freq     - The frequency of the sine wave.
//   initialY - The initial y-axis value of the sine wave.
//   fidelity - The sample rate to emulate.
//
//   pointsLength - The number of points to generate.
//   width        - The width of the render window.
func GenerateSinePoints(freq, initialPhase, fidelity float32, pointsLength, width int) []Point {
	spacing := float32(2.0 / 640.0)
	points := make([]Point, pointsLength)

	y := initialPhase
	x := float32(-1.0)
	for i := 0; i < pointsLength; i++ {
		points[i] = Point{x, float32(math.Sin(float64(y))) * 0.5}

		x += spacing
		y += 2 * math.Pi * (freq / fidelity)
	}

	return points
}

// Only for testing purposes.
func DefaultGenerateSinePoints(freq, initialPhase float32) []Point {
	return GenerateSinePoints(freq, initialPhase, 44100, 640, 640)
}
