package aoc2019

import (
	"fmt"
	"io/ioutil"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type day17TileType rune

const (
	day17TileTypeNewline    day17TileType = 10
	day17TileTypeRobotDown  day17TileType = 'v'
	day17TileTypeRobotLeft  day17TileType = '<'
	day17TileTypeRobotLost  day17TileType = 'X'
	day17TileTypeRobotRight day17TileType = '>'
	day17TileTypeRobotUp    day17TileType = '^'
	day17TileTypeScaffold   day17TileType = 35 // #
	day17TileTypeSpace      day17TileType = 46 // .
)

type day17Tile struct {
	X, Y int64
	Type day17TileType
}

type day17Grid map[string]*day17Tile

func (d day17Grid) bounds() (minX, minY, maxX, maxY int64) {
	minX, minY = math.MaxInt64, math.MaxInt64

	for _, t := range d {
		if t.X < minX {
			minX = t.X
		}
		if t.Y < minY {
			minY = t.Y
		}
		if t.X > maxX {
			maxX = t.X
		}
		if t.Y > maxY {
			maxY = t.Y
		}
	}

	return
}

func (d day17Grid) findCurrentPosition() (int64, int64, day17TileType) {
	for _, t := range d {
		if t.Type == day17TileTypeRobotDown || t.Type == day17TileTypeRobotLeft || t.Type == day17TileTypeRobotRight || t.Type == day17TileTypeRobotUp {
			return t.X, t.Y, t.Type
		}
	}
	return -1, -1, day17TileTypeRobotLost
}

func (d day17Grid) findPath() []int64 {
	var (
		count           int64
		directives      []int64
		x, y, direction = d.findCurrentPosition()
	)

	move := func() {
		count++

		switch direction {
		case day17TileTypeRobotDown:
			y += 1
		case day17TileTypeRobotLeft:
			x -= 1
		case day17TileTypeRobotRight:
			x += 1
		case day17TileTypeRobotUp:
			y -= 1
		}
	}

	opportunities := func() (forward, left, right bool) {
		var tileFwd, tileL, tileR *day17Tile

		switch direction {
		case day17TileTypeRobotDown:
			tileFwd, tileL, tileR = d[d.mapKey(x, y+1)], d[d.mapKey(x+1, y)], d[d.mapKey(x-1, y)]
		case day17TileTypeRobotLeft:
			tileFwd, tileL, tileR = d[d.mapKey(x-1, y)], d[d.mapKey(x, y+1)], d[d.mapKey(x, y-1)]
		case day17TileTypeRobotRight:
			tileFwd, tileL, tileR = d[d.mapKey(x+1, y)], d[d.mapKey(x, y-1)], d[d.mapKey(x, y+1)]
		case day17TileTypeRobotUp:
			tileFwd, tileL, tileR = d[d.mapKey(x, y-1)], d[d.mapKey(x-1, y)], d[d.mapKey(x+1, y)]
		}

		forward = tileFwd != nil && tileFwd.Type == day17TileTypeScaffold
		left = tileL != nil && tileL.Type == day17TileTypeScaffold
		right = tileR != nil && tileR.Type == day17TileTypeScaffold

		return
	}

	turn := func(left bool) {
		var (
			dl = []day17TileType{day17TileTypeRobotUp, day17TileTypeRobotLeft, day17TileTypeRobotDown, day17TileTypeRobotRight}
			dr = []day17TileType{day17TileTypeRobotUp, day17TileTypeRobotRight, day17TileTypeRobotDown, day17TileTypeRobotLeft}
			nd []day17TileType
		)

		if count > 0 {
			directives = append(directives, count)
			count = 0
		}

		if left {
			nd = dl
			directives = append(directives, int64('L'))
		} else {
			nd = dr
			directives = append(directives, int64('R'))
		}

		var cd int
		for i, t := range nd {
			if t == direction {
				cd = i
				break
			}
		}
		if cd == len(nd)-1 {
			cd = -1
		}

		direction = nd[cd+1]
	}

	// Do the movement
	for {
		fw, l, r := opportunities()

		if fw {
			move()
			continue
		}

		if l {
			turn(true)
			continue
		}

		if r {
			turn(false)
			continue
		}

		// Nothing possible, must be the end
		directives = append(directives, count)
		break
	}

	return directives
}

func (d day17Grid) isScaffoldIntersection(x, y int64) bool {
	t := d[d.mapKey(x, y)]

	if t.Type == day17TileTypeSpace {
		// Space cannot be a scaffold intersection
		return false
	}

	var count int64
	for _, at := range []*day17Tile{
		d[d.mapKey(x-1, y)],
		d[d.mapKey(x+1, y)],
		d[d.mapKey(x, y-1)],
		d[d.mapKey(x, y+1)],
	} {
		if at != nil && at.Type != day17TileTypeSpace {
			count++
		}
	}

	// One adjacent scaffold is a line
	// Two adjacent scaffolds is a curve or a straight line
	// Three or four adjacent scaffolds cannot be a line or curve, must be intersection
	return count >= 3
}

func (d day17Grid) mapKey(x, y int64) string { return fmt.Sprintf("%d:%d", x, y) }

func (d day17Grid) print() {
	var minX, minY, maxX, maxY = d.bounds()

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			fmt.Printf("%s", string(d[d.mapKey(x, y)].Type))
		}
		fmt.Println()
	}
}

func day17ReadGrid(code []int64) (day17Grid, error) {
	var (
		grid = make(day17Grid)
		out  = make(chan int64)
		x, y int64
	)

	go executeIntcode(code, nil, out)

	for o := range out {
		switch day17TileType(o) {

		case day17TileTypeNewline:
			y += 1
			x = 0

		case day17TileTypeRobotDown, day17TileTypeRobotLeft, day17TileTypeRobotRight, day17TileTypeRobotUp:
			fallthrough // Not yet used

		case day17TileTypeRobotLost:
			fallthrough // Not yet used

		case day17TileTypeScaffold, day17TileTypeSpace:
			grid[grid.mapKey(x, y)] = &day17Tile{X: x, Y: y, Type: day17TileType(o)}
			x += 1

		default:
			return nil, errors.Errorf("Invalid character %d", o)

		}
	}

	return grid, nil
}

func solveDay17Part1(inFile string) (int64, error) {
	rawCode, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read intcode")
	}

	code, err := parseIntcode(strings.TrimSpace(string(rawCode)))
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse intcode")
	}

	grid, err := day17ReadGrid(code)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read grid")
	}

	var (
		apSum                  int64
		minX, minY, maxX, maxY = grid.bounds()
	)

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {

			if grid.isScaffoldIntersection(x, y) {
				apSum += x * y
			}

		}
	}

	return apSum, nil
}

func solveDay17Part2(inFile string) (int64, error) {
	rawCode, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read intcode")
	}

	code, err := parseIntcode(strings.TrimSpace(string(rawCode)))
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse intcode")
	}

	grid, err := day17ReadGrid(code)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read grid")
	}

	// Get the instructions to follow the whole path
	path := grid.findPath()

	isMovementFunction := func(n int64) bool {
		return n == int64('A') || n == int64('B') || n == int64('C')
	}

	feedSlice := func(in chan int64, s []int64) {
		for i, n := range s {
			if isMovementFunction(n) || n == int64('L') || n == int64('R') || n == int64('n') {
				in <- n
			} else {
				for _, c := range strconv.FormatInt(n, 10) {
					in <- int64(c)
				}
			}

			if i == len(s)-1 {
				in <- 10 // Newline to terminate
			} else {
				in <- 44 // Comma to delimit chars
			}
		}
	}

	findSubmatches := func(haystack, needle []int64) []int64 {
		var matches []int64

		for i := 0; i < len(haystack)-len(needle)+1; i++ {
			if reflect.DeepEqual(haystack[i:i+len(needle)], needle) {
				matches = append(matches, int64(i))
				i += len(needle) - 1
			}
		}

		return matches
	}

	var (
		pattern [][]int64
		main    = make([]int64, len(path))
	)
	copy(main, path)

	for len(findSubmatches(main, []int64{int64('L')}))+len(findSubmatches(main, []int64{int64('R')})) > 0 {
		var start, length int64 = 0, 2

		for isMovementFunction(main[start]) {
			start++
		}

		for len(findSubmatches(main, main[start:start+length])) > 1 && !isMovementFunction(main[start+length-1]) {
			length += 2
		}

		length -= 2 // Revert last addition as it caused trouble

		var patternID = int64('A') + int64(len(pattern))
		pattern = append(pattern, main[start:start+length])

		var (
			pos            int64
			submatchStarts = findSubmatches(main, pattern[len(pattern)-1])
			tmpMain        []int64
		)

		for _, smStart := range submatchStarts {
			tmpMain = append(tmpMain, main[pos:smStart]...)
			tmpMain = append(tmpMain, patternID)
			pos = smStart + length
		}
		tmpMain = append(tmpMain, main[pos:]...)
		main = tmpMain
	}

	if len(pattern) != 3 {
		return 0, errors.Errorf("Required more than 3 pattern: %d", len(pattern))
	}

	// Feed program with new data
	var (
		in  = make(chan int64, 1000) // I could calculate the length but I don't care
		out = make(chan int64)
		wg  = new(sync.WaitGroup)
	)

	// Feed main movement routine
	feedSlice(in, main)
	// Feed movement routines
	feedSlice(in, pattern[0]) // A
	feedSlice(in, pattern[1]) // B
	feedSlice(in, pattern[2]) // C
	// Answer "continuous video feed" question
	feedSlice(in, []int64{int64('n')})

	code, err = parseIntcode(strings.TrimSpace(string(rawCode)))
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse intcode")
	}

	// Force the vacuum robot to wake up by changing the value in your
	// ASCII program at address 0 from 1 to 2.
	code[0] = 2

	// Execute the program and throw away all but last output, we know
	// how the grid looks
	wg.Add(1)
	go executeIntcode(code, in, out)

	var result int64
	go func() {
		for o := range out {
			result = o
		}
		wg.Done()
	}()

	wg.Wait()

	return result, nil
}
