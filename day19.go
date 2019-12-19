package aoc2019

import (
	"io/ioutil"
	"log"
	"math"
	"strings"

	"github.com/pkg/errors"
)

type day19Fact uint8

const (
	day19FactNoMoveDownPossible day19Fact = 1 << iota
	day19FactNoMoveLeftPossible
	day19FactNoMoveRightPossible
	day19FactNoMoveTopPossible
	day19FactOriginNotInBeam
	day19FactX100NotInBeam
	day19FactY100NotInBeam
)

func (f day19Fact) has(t day19Fact) bool { return f&t != 0 }
func (f day19Fact) String() string {
	return map[day19Fact]string{
		day19FactNoMoveDownPossible:  "NoMoveDownPossible",
		day19FactNoMoveLeftPossible:  "NoMoveLeftPossible",
		day19FactNoMoveRightPossible: "NoMoveRightPossible",
		day19FactNoMoveTopPossible:   "NoMoveTopPossible",
		day19FactOriginNotInBeam:     "OriginNotInBeam",
		day19FactX100NotInBeam:       "X100NotInBeam",
		day19FactY100NotInBeam:       "Y100NotInBeam",
	}[f]
}

func day19CountFieldsInTractorBeam(code []int64, maxX, maxY int64) int64 {
	var count int64

	for y := int64(0); y <= maxY; y++ {
		for x := int64(0); x <= maxX; x++ {

			var (
				in  = make(chan int64)
				out = make(chan int64)
			)
			go executeIntcode(code, in, out)

			// Submit coordinates
			in <- x
			in <- y

			// Count fields with tractor beam
			o := <-out
			count += o
		}
	}

	return count
}

func day19Find100x100ShipPlace(code []int64) int64 {
	var (
		x, y     int64
		bvL, bvR int64
	)

	getBeamVectors := func() {
		for x := int64(0); x <= math.MaxInt64; x++ {

			var (
				in  = make(chan int64)
				out = make(chan int64)
			)
			go executeIntcode(code, in, out)

			// Submit coordinates
			in <- x
			in <- 100

			// Count fields with tractor beam
			o := <-out
			switch {
			case o == 0 && bvR == 0:
				// Did not yet find the beam
			case o == 1 && bvL == 0:
				// Left beam end has not been set
				bvL, bvR = x, x
			case o == 1 && x > bvR:
				// New right "edge" found
				bvR = x
			case 0 == 0 && bvR > 0:
				// End of right edge found, end
				return
			}

		}
	}

	checkCoordinate := func(x, y int64) day19Fact {
		var facts day19Fact

		for flag, c := range map[day19Fact][2]int64{
			day19FactNoMoveDownPossible:  {x, y + 100},
			day19FactNoMoveLeftPossible:  {x - 1, y + 99},
			day19FactNoMoveRightPossible: {x + 100, y},
			day19FactNoMoveTopPossible:   {x + 99, y - 1},
			day19FactOriginNotInBeam:     {x, y},
			day19FactX100NotInBeam:       {x + 99, y},
			day19FactY100NotInBeam:       {x, y + 99},
		} {
			if c[0] < 0 || c[1] < 0 {
				// Never check negative coordinates
				facts |= flag
				continue
			}

			var (
				in  = make(chan int64)
				out = make(chan int64)
			)
			defer close(in)

			go executeIntcode(code, in, out)

			// Submit coordinates
			in <- c[0]
			in <- c[1]

			o := <-out
			if o == 0 {
				facts |= flag
			}

		}

		return facts
	}

	// Initially get vectors for beam edges
	getBeamVectors()

	// Try to get to the best position through movement
	var success bool
	for !success {
		var f = checkCoordinate(x, y)

		switch {
		case f.has(day19FactOriginNotInBeam):
			panic("Origin placed outside beam")

		case f.has(day19FactX100NotInBeam) || f.has(day19FactY100NotInBeam):
			// Ship does not fit, move further away
			y += 100
			x = ((bvL * (y / 100)) + (bvR * (y / 100))) / 2

		case !f.has(day19FactNoMoveTopPossible):
			// Ship can be moved up in beam
			y -= 1

		case !f.has(day19FactNoMoveLeftPossible):
			// Ship can be moved left in beam
			x -= 1

		case f.has(day19FactNoMoveTopPossible) && f.has(day19FactNoMoveLeftPossible):
			// No movement towards ship is possible and ship fits in: Perfect
			success = true

		}
	}

	log.Printf("Result before force-move: x=%d y=%d res=%d", x, y, x*10000+y)

	// This MIGHT not be the perfect position, force further movement
	// by ignoring some factors set above in order to compensate inaccurate
	// vectors due to working with integers only
	var fX, fY = x, y
	for success {
		var f = checkCoordinate(fX, fY)

		if !f.has(day19FactX100NotInBeam) && !f.has(day19FactY100NotInBeam) {
			// Found a better position through forcing
			x, y = fX, fY
		}

		switch {
		case !f.has(day19FactNoMoveTopPossible):
			// Ship can be moved up in beam
			fY -= 1

		case f.has(day19FactOriginNotInBeam):
			// No further tests, there will be no better position
			success = false

		default:
			fX -= 1
		}
	}

	log.Printf("Result after force-move:  x=%d y=%d res=%d", x, y, x*10000+y)

	return x*10000 + y
}

func solveDay19Part1(inFile string) (int64, error) {
	rawCode, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read intcode")
	}

	code, err := parseIntcode(strings.TrimSpace(string(rawCode)))
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse intcode")
	}

	return day19CountFieldsInTractorBeam(code, 49, 49), nil
}

func solveDay19Part2(inFile string) (int64, error) {
	rawCode, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read intcode")
	}

	code, err := parseIntcode(strings.TrimSpace(string(rawCode)))
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse intcode")
	}

	return day19Find100x100ShipPlace(code), nil
}
