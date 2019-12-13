package aoc2019

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type day13TileType int64

const (
	day13TileTypeEmpty   day13TileType = iota // No game object appears in this tile.
	day13TileTypeWall                         // Walls are indestructible barriers.
	day13TileTypeBlock                        // Blocks can be broken by the ball.
	day13TileTypeHPaddle                      // The paddle is indestructible.
	day13TileTypeBall                         // The ball moves diagonally and bounces off objects.
)

type day13Tile struct {
	X, Y int64
	Type day13TileType
}

func (d day13Tile) key() string { return fmt.Sprintf("%d:%d", d.X, d.Y) }

type day13Field struct {
	tiles      map[string]day13Tile
	score      int64
	maxX, maxY int64
	ball       day13Tile

	tileLock sync.RWMutex
}

func (d *day13Field) getPosition(tileType day13TileType) (int64, int64) {
	d.tileLock.RLock()
	defer d.tileLock.RUnlock()

	if tileType == day13TileTypeBall {
		return d.ball.X, d.ball.Y
	}

	for _, t := range d.tiles {
		if t.Type == tileType {
			return t.X, t.Y
		}
	}

	return -1, -1
}

func (d *day13Field) remainingTiles(tileType day13TileType) int {
	d.tileLock.RLock()
	defer d.tileLock.RUnlock()

	var count int
	for _, t := range d.tiles {
		if t.Type == tileType {
			count++
		}
	}
	return count
}

func day13PlayGame(code []int64, field *day13Field, in func() (int64, error)) error {
	var (
		err    error
		out    = make(chan int64)
		output []int64
		wg     = new(sync.WaitGroup)
	)

	in = func(f func() (int64, error)) func() (int64, error) {
		return func() (int64, error) {
			// Ensure every output is already processed before deciding on the input.
			// Without this there is a race which randomly breaks the program in
			// different ways m(
			for len(output) > 0 {
				time.Sleep(time.Nanosecond)
			}

			return f()
		}
	}(in)

	go func() {
		for o := range out {
			output = append(output, o)

			if len(output) >= 3 {
				field.tileLock.Lock()

				if output[0] == -1 && output[1] == 0 {
					field.score = output[2]
					output = output[3:]

					field.tileLock.Unlock()
					continue
				}

				tile := day13Tile{X: output[0], Y: output[1], Type: day13TileType(output[2])}
				field.tiles[tile.key()] = tile

				if tile.Type == day13TileTypeBall {
					// Ball sometimes disappear, store it extra
					field.ball = tile
				}

				if tile.X > field.maxX {
					field.maxX = tile.X
				}
				if tile.Y > field.maxY {
					field.maxY = tile.Y
				}
				output = output[3:]

				field.tileLock.Unlock()
			}
		}
		wg.Done()
	}()

	wg.Add(1)
	_, err = executeIntcode(code, in, out)
	wg.Wait()

	return errors.Wrap(err, "Unable to execute intcode")
}

func solveDay13Part1(inFile string) (int, error) {
	raw, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read code")
	}

	code, err := parseIntcode(strings.TrimSpace(string(raw)))
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse code")
	}

	var field = &day13Field{tiles: make(map[string]day13Tile)}
	if err := day13PlayGame(code, field, nil); err != nil {
		return 0, errors.Wrap(err, "Unable to draw field")
	}

	return field.remainingTiles(day13TileTypeBlock), nil
}

func solveDay13Part2(inFile string) (int64, error) {
	raw, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read code")
	}

	code, err := parseIntcode(strings.TrimSpace(string(raw)))
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse code")
	}

	// Let the game initialize once to get tick count for initialization
	var field = &day13Field{tiles: make(map[string]day13Tile)}

	// Start the real game
	code[0] = 2 // Insert two quarters

	// Input callback
	var in = func() (int64, error) {
		field.tileLock.RLock()
		defer field.tileLock.RUnlock()

		var (
			ballX, _   = field.getPosition(day13TileTypeBall)
			dir        int64
			paddleX, _ = field.getPosition(day13TileTypeHPaddle)
		)

		// Move paddle in ball direction
		switch {
		case ballX == paddleX || ballX == -1:
			dir = 0
		case ballX < paddleX:
			dir = -1
		case ballX > paddleX:
			dir = 1
		}

		return dir, nil
	}

	//intcodeDebugging = true
	if err := day13PlayGame(code, field, in); err != nil {
		log.Printf("err=%s", err)
	}

	if rb := field.remainingTiles(day13TileTypeBlock); rb > 0 {
		return 0, errors.Errorf("Game logic quit with blocks remaining: %d", rb)
	}

	return field.score, nil
}
