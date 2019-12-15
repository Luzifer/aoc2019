package aoc2019

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"strings"

	"github.com/pkg/errors"
)

type day15TileType int64

const (
	day15TileTypeWall day15TileType = iota
	day15TileTypeFloor
	day15TileTypeOxygen
	day15TileTypeUnknown
)

type day15Tile struct {
	X, Y int64
	Type day15TileType

	distFromStart int64
}

func (d day15Tile) key() string { return fmt.Sprintf("%d:%d", d.X, d.Y) }

type day15Grid map[string]*day15Tile

func (d day15Grid) annotateDistance(x, y, dist int64) {
	var t = d.getTile(x, y)
	if t == nil {
		panic("Access to non-existent tile")
	}

	if t.Type != day15TileTypeFloor && t.Type != day15TileTypeOxygen {
		// Do not annotate distance on walls
		return
	}

	if t.distFromStart <= dist {
		// Distance already set, no need to set again
		return
	}

	// Set distance
	t.distFromStart = dist

	// Annotate next fields
	d.annotateDistance(x-1, y, dist+1)
	d.annotateDistance(x+1, y, dist+1)
	d.annotateDistance(x, y-1, dist+1)
	d.annotateDistance(x, y+1, dist+1)
}

func (d day15Grid) bounds() (minX, minY, maxX, maxY int64) {
	minX = math.MaxInt64
	minY = math.MaxInt64

	for _, t := range d {
		if t.X < minX {
			minX = t.X
		}
		if t.X > maxX {
			maxX = t.X
		}
		if t.Y < minY {
			minY = t.Y
		}
		if t.Y > maxY {
			maxY = t.Y
		}
	}

	return
}

func (d day15Grid) getTile(x, y int64) *day15Tile {
	if v, ok := d[day15Tile{X: x, Y: y}.key()]; ok {
		return v
	}
	return nil
}

func (d day15Grid) getTileType(x, y int64) day15TileType {
	if v := d.getTile(x, y); v != nil {
		return v.Type
	}
	return day15TileTypeUnknown
}

func (d day15Grid) print() {
	minX, minY, maxX, maxY := d.bounds()
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if x == 0 && y == 0 {
				fmt.Printf("@")
				continue
			}

			t := d.getTileType(x, y)
			switch t {
			case day15TileTypeFloor:
				fmt.Printf(".")
			case day15TileTypeWall:
				fmt.Printf("\u2588")
			case day15TileTypeOxygen:
				fmt.Printf("X")
			case day15TileTypeUnknown:
				fmt.Printf("\u2593")
			}
		}
		fmt.Println()
	}
}

func day15ScanGrid(code []int64) (day15Grid, error) {
	var (
		grid       = make(day15Grid)
		posX, posY int64
		rotation   int64 = 1 // start facing north
	)

	recordPosition := func(success bool, tile day15TileType) bool {
		var nPosX, nPosY = posX, posY

		if success {
			defer func() { posX, posY = nPosX, nPosY }()
		}

		switch rotation {
		case 1:
			nPosX -= 1
		case 2:
			nPosX += 1
		case 3:
			nPosY -= 1
		case 4:
			nPosY += 1
		}

		if t := grid.getTileType(nPosX, nPosY); t == tile {
			return true
		}

		nT := day15Tile{X: nPosX, Y: nPosY, Type: tile, distFromStart: math.MaxInt64}
		grid[nT.key()] = &nT

		return false
	}

	rotate := func(forward bool) {
		var (
			fr = []int64{1, 4, 2, 3}
			br = []int64{1, 3, 2, 4}
			r  []int64
		)

		if forward {
			r = fr
		} else {
			r = br
		}

		nP := int64IndexOf(r, rotation) + 1
		if nP == len(r) {
			nP = 0
		}

		rotation = r[nP]
	}

	var (
		in  = make(chan int64)
		out = make(chan int64)
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go executeIntcodeWithParams(intcodeParams{
		Code:    code,
		Context: ctx,
		In:      in,
		Out:     out,
	})

	// Start by moving
	in <- rotation

	for res := range out {
		var alreadyKnown bool

		switch day15TileType(res) {
		case day15TileTypeWall:
			// Ran into wall, not a successful move
			alreadyKnown = recordPosition(false, day15TileType(res))
			// Rotate once forward
			rotate(false)

		case day15TileTypeFloor, day15TileTypeOxygen:
			// Moved to new tile, successful move
			alreadyKnown = recordPosition(true, day15TileType(res))
			// Rotate once backward
			rotate(true)

		default:
			// Thefuck?
			return grid, errors.Errorf("Invalid tile type detected: %d", res)
		}

		_ = alreadyKnown
		if posX == 0 && posY == 0 {
			// We've reached a position twice, let's quit the program but not
			// yet the function as the input then will hang
			cancel()
		}

		in <- rotation
	}

	return grid, nil
}

func int64IndexOf(s []int64, e int64) int {
	for i, se := range s {
		if se == e {
			return i
		}
	}
	return -1
}

func solveDay15Part1(inFile string) (int64, error) {
	raw, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read input")
	}

	code, err := parseIntcode(strings.TrimSpace(string(raw)))
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse Intcode")
	}

	grid, err := day15ScanGrid(code)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to scan grid")
	}

	grid.annotateDistance(0, 0, 0)

	var oxygen *day15Tile
	for _, t := range grid {
		if t.Type == day15TileTypeOxygen {
			oxygen = t
			break
		}
	}

	return oxygen.distFromStart, nil
}

func solveDay15Part2(inFile string) (int64, error) {
	raw, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read input")
	}

	code, err := parseIntcode(strings.TrimSpace(string(raw)))
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse Intcode")
	}

	grid, err := day15ScanGrid(code)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to scan grid")
	}

	var oxygenSystem *day15Tile
	for _, t := range grid {
		if t.Type == day15TileTypeOxygen {
			oxygenSystem = t
			break
		}
	}

	grid.annotateDistance(oxygenSystem.X, oxygenSystem.Y, 0)

	var farthestTile = oxygenSystem
	for _, t := range grid {
		if t.distFromStart > farthestTile.distFromStart && t.Type == day15TileTypeFloor {
			farthestTile = t
		}
	}

	return farthestTile.distFromStart, nil
}
