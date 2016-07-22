package los

import "math"

func MakeLine(xStart int, yStart int, xEnd int, yEnd int) [][]int {
	path := [][]int{}

	deltaX := math.Floor(math.Abs(float64(xEnd - xStart)))
	deltaY := math.Floor(math.Abs(float64(yEnd - yStart)))

	var xStep float64
	var yStep float64

	if xEnd >= xStart {
		xStep = 1
	} else {
		xStep = -1
	}

	if yEnd >= yStart {
		yStep = 1
	} else {
		yStep = -1
	}

	e := deltaX - deltaY
	xGrid := float64(xStart)
	yGrid := float64(yStart)

	for !(xGrid == float64(xEnd) && yGrid == float64(yEnd)) {
		ret := []int{int(xGrid), int(yGrid)}
		path = append(path, ret)
		twoError := 2 * e
		if twoError > (-1 * deltaY) {
			e -= deltaY
			xGrid += xStep
		}
		if twoError < deltaX {
			e += deltaX
			yGrid += yStep
		}
	}

	return path
}
