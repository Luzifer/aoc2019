package aoc2019

import (
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type day3LineSegment [4]int

func (d day3LineSegment) Draw(img *image.RGBA, col color.Color) {
	var ptMod func(x, y int) (int, int, bool)

	switch {
	case d.IsVertical() && d[1] < d[3]:
		ptMod = func(x, y int) (int, int, bool) { return x, y + 1, y+1 <= d[3] }
	case d.IsVertical() && d[1] > d[3]:
		ptMod = func(x, y int) (int, int, bool) { return x, y - 1, y-1 >= d[3] }
	case !d.IsVertical() && d[0] < d[2]:
		ptMod = func(x, y int) (int, int, bool) { return x + 1, y, x+1 <= d[2] }
	case !d.IsVertical() && d[0] > d[2]:
		ptMod = func(x, y int) (int, int, bool) { return x - 1, y, x-1 >= d[2] }
	}

	var pX, pY, ok = d[0], d[1], true

	for ok {
		img.Set(pX, pY, col)
		pX, pY, ok = ptMod(pX, pY)
	}
}

func (d day3LineSegment) GetIntersection(in day3LineSegment) ([2]int, bool) {
	if d.IsVertical() == in.IsVertical() {
		// Both lines have the same direction, there is no intersection
		// NOTE: This might yield false negative when lines are overlapping?
		return [2]int{0, 0}, false
	}

	var (
		vert, horiz day3LineSegment
	)

	if d.IsVertical() {
		vert, horiz = d, in
	} else {
		vert, horiz = in, d
	}

	// Normalize lines for intersection finding
	if vert[1] > vert[3] {
		vert[1], vert[3] = vert[3], vert[1]
	}

	if horiz[0] > horiz[2] {
		horiz[0], horiz[2] = horiz[2], horiz[0]
	}

	if vert[0] < horiz[0] || vert[0] > horiz[2] || horiz[1] < vert[1] || horiz[1] > vert[3] {
		// Lines do not have an intersection
		return [2]int{0, 0}, false
	}

	return [2]int{vert[0], horiz[1]}, true
}

func (d day3LineSegment) HasPoint(x, y int) bool {
	switch {
	case d.IsVertical() && x != d[0]: // vertical line, test Y not on line
		fallthrough
	case !d.IsVertical() && y != d[1]: // horizontal line, text X not on line
		fallthrough
	case d.IsVertical() && d[1] < d[3] && (y < d[1] || y > d[3]): // vertical low to high, test Y not in range
		fallthrough
	case d.IsVertical() && d[1] > d[3] && (y > d[1] || y < d[3]): // vertical high to low, test Y not in range
		fallthrough
	case !d.IsVertical() && d[0] < d[2] && (x < d[0] || x > d[2]): // horizontal left to right, test X not in range
		fallthrough
	case !d.IsVertical() && d[0] > d[2] && (x > d[0] || x < d[2]): // horizontal right to left, test X not in range

		return false
	}

	return true
}

func (d day3LineSegment) IsVertical() bool {
	return d[0] == d[2]
}

func (d day3LineSegment) Steps() int {
	if d.IsVertical() {
		return int(math.Abs(float64(d[1] - d[3])))
	}
	return int(math.Abs(float64(d[0] - d[2])))
}

func (d day3LineSegment) StepsToPoint(x, y int) int {
	if d.IsVertical() {
		return int(math.Abs(float64(d[1] - y)))
	}
	return int(math.Abs(float64(d[0] - x)))
}

type day3Line []day3LineSegment

func (d day3Line) GetIntersections(in day3Line) [][2]int {
	var inter [][2]int

	for _, dseg := range d {
		for _, iseg := range in {
			if is, ok := dseg.GetIntersection(iseg); ok {
				inter = append(inter, is)
			}
		}
	}

	return inter
}

func debugDay3DrawImage(l1, l2 day3Line) error {
	dest := image.NewRGBA(image.Rect(-500, -500, 500, 500))

	// Draw l1 in red
	for _, ls := range l1 {
		ls.Draw(dest, color.RGBA{0xff, 0x0, 0x0, 0xff})
	}

	// Draw l2 in blue
	for _, ls := range l2 {
		ls.Draw(dest, color.RGBA{0x0, 0x0, 0xff, 0xff})
	}

	// Draw detected intersections in purple
	for _, i := range l1.GetIntersections(l2) {
		dest.SetRGBA(i[0], i[1], color.RGBA{0xff, 0x0, 0xff, 0xff})
	}

	// Draw "home-point" in green
	dest.SetRGBA(0, 0, color.RGBA{0x0, 0xff, 0x0, 0xff})

	f, err := os.Create("day03_debug.png")
	if err != nil {
		return errors.Wrap(err, "Unable to open output file")
	}
	defer f.Close()

	return errors.Wrap(png.Encode(f, dest), "Unable to write image")
}

func getDay3MinIntersectionDistance(l1, l2 day3Line, originX, originY int) int {
	var (
		inter = l1.GetIntersections(l2)
		min   = math.MaxInt64
	)

	for _, i := range inter {
		dist := manhattenDistance(originX, originY, i[0], i[1])
		if dist > 0 && dist < min {
			min = dist
		}
	}

	return min
}

func getDay3MinIntersectionSteps(l1, l2 day3Line) int {
	var minSteps int = math.MaxInt64

	for _, is := range l1.GetIntersections(l2) {
		var combinedsteps int

		for _, l := range []day3Line{l1, l2} {
			for _, ls := range l {
				if ls.HasPoint(is[0], is[1]) {
					combinedsteps += ls.StepsToPoint(is[0], is[1])
					break
				}
				combinedsteps += ls.Steps()
			}
		}

		if combinedsteps > 0 && combinedsteps < minSteps {
			minSteps = combinedsteps
		}
	}

	return minSteps
}

func parseDay3LineDefinition(definition string, startX, startY int) (day3Line, error) {
	var (
		directions = strings.Split(strings.TrimSpace(definition), ",")
		pX, pY     = startX, startY
		out        day3Line
	)

	for _, d := range directions {
		l, err := strconv.Atoi(d[1:])
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to parse direction %q", d)
		}

		var tX, tY int
		switch d[0] {
		case 'D':
			tX, tY = pX, pY-l

		case 'L':
			tX, tY = pX-l, pY

		case 'R':
			tX, tY = pX+l, pY

		case 'U':
			tX, tY = pX, pY+l

		default:
			return nil, errors.Wrapf(err, "Invalid direction %q given", d)
		}

		out = append(out, [4]int{pX, pY, tX, tY})
		pX, pY = tX, tY
	}

	return out, nil
}

func solveDay3Part1(inFile string) (int, error) {
	raw, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read input")
	}

	defs := strings.Split(string(raw), "\n")
	l1, err := parseDay3LineDefinition(defs[0], 0, 0)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse L1")
	}

	l2, err := parseDay3LineDefinition(defs[1], 0, 0)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse L1")
	}

	return getDay3MinIntersectionDistance(l1, l2, 0, 0), nil
}

func solveDay3Part2(inFile string) (int, error) {
	raw, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read input")
	}

	defs := strings.Split(string(raw), "\n")
	l1, err := parseDay3LineDefinition(defs[0], 0, 0)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse L1")
	}

	l2, err := parseDay3LineDefinition(defs[1], 0, 0)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse L1")
	}

	return getDay3MinIntersectionSteps(l1, l2), nil
}
