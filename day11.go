package aoc2019

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type day11PaintDirective struct {
	X, Y  int
	Color int64
}

func (d day11PaintDirective) positionKey() string { return fmt.Sprintf("%d:%d", d.X, d.Y) }

func day11ExecutePaintRobot(inFile string, startPanelColor int64) ([]day11PaintDirective, error) {
	var (
		posX, posY int
		direction  int // Clock-hand direction
		directives = map[string]day11PaintDirective{
			"0:0": {X: 0, Y: 0, Color: startPanelColor},
		}
	)

	move := func() {
		switch direction {
		case 0:
			posY -= 1
		case 3:
			posX += 1
		case 6:
			posY += 1
		case 9:
			posX -= 1
		default:
			panic(errors.Errorf("Invalid direction: %d", direction))
		}
	}

	rotate := func(dir int64) {
		var factor int
		if dir == 1 {
			// turn right
			factor = 1
		} else {
			// turn left
			factor = -1
		}
		direction = direction + factor*3
		if direction < 0 {
			direction += 12
		}
		if direction > 11 {
			direction -= 12
		}
	}

	// Initialize code
	rawCode, err := ioutil.ReadFile(inFile)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read intcode")
	}

	code, err := parseIntcode(strings.TrimSpace(string(rawCode)))
	if err != nil {
		return nil, errors.Wrap(err, "Unable to parse intcode")
	}

	// Run program in background
	var (
		in          = make(chan int64, 1)
		out         = make(chan int64, 2)
		outputs     []int64
		programExit bool
		wg          = new(sync.WaitGroup)
	)

	go executeIntcode(code, in, out)
	go func() {
		for o := range out {
			outputs = append(outputs, o)
			wg.Done()
		}
		programExit = true
	}()

	// Start feeding scan results
	for !programExit {
		// Get previous directive for current position
		var dir = day11PaintDirective{X: posX, Y: posY}
		if d, ok := directives[dir.positionKey()]; ok {
			dir = d
		}

		// Feed current color and expect two output
		wg.Add(2)
		in <- dir.Color
		wg.Wait()

		// Set current color
		dir.Color = outputs[0]
		directives[dir.positionKey()] = dir
		// Rotate robot
		rotate(outputs[1])
		// Move
		move()

		// Reset outputs
		outputs = nil
	}

	var result []day11PaintDirective
	for _, v := range directives {
		result = append(result, v)
	}

	return result, nil
}

func solveDay11Part1(inFile string) (int, error) {
	dirs, err := day11ExecutePaintRobot(inFile, 0)
	return len(dirs), errors.Wrap(err, "Unable to execute robot")
}

func solveDay11Part2(inFile string) error {
	dirs, err := day11ExecutePaintRobot(inFile, 1)
	if err != nil {
		return errors.Wrap(err, "Unable to execute robot")
	}

	var minX, minY, maxX, maxY = math.MaxInt64, math.MaxInt64, 0, 0
	for _, dir := range dirs {
		if dir.X < minX {
			minX = dir.X
		}
		if dir.X > maxX {
			maxX = dir.X
		}
		if dir.Y < minY {
			minY = dir.Y
		}
		if dir.Y > maxY {
			maxY = dir.Y
		}
	}

	colors := map[int64]color.Color{
		0: color.RGBA{0x0, 0x0, 0x0, 0xff},
		1: color.RGBA{0xff, 0xff, 0xff, 0xff},
	}

	img := image.NewRGBA(image.Rect(minX-5, minY-5, maxX+5, maxY+5))
	for _, dir := range dirs {
		img.Set(dir.X, dir.Y, colors[dir.Color])
	}

	w, err := os.Create("day11_image.png")
	if err != nil {
		return errors.Wrap(err, "Unable to open result image file")
	}
	defer w.Close()

	return errors.Wrap(png.Encode(w, img), "Unable to store image")
}
