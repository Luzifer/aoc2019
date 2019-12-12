package aoc2019

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

type day12MoonSystem []*day12Moon

func day12MoonSystemFromFile(inFile string) (day12MoonSystem, error) {
	f, err := os.Open(inFile)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to open inFile")
	}
	defer f.Close()

	return day12MoonSystemFromReader(f)
}

func day12MoonSystemFromReader(r io.Reader) (day12MoonSystem, error) {
	var moons day12MoonSystem

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		m, err := day12MoonFromScan(scanner.Text())
		if err != nil {
			return nil, errors.Wrap(err, "Unable to read input file")
		}
		moons = append(moons, m)
	}

	return moons, errors.Wrap(scanner.Err(), "Unable to scan file")
}

func (d day12MoonSystem) move(steps int) {
	for i := 0; i < steps; i++ {
		d.updateVelocity()
		for _, m := range d {
			m.move()
		}
	}
}

func (d day12MoonSystem) totalEnergy() int64 {
	var e int64
	for _, m := range d {
		e += m.totalEnergy()
	}
	return e
}

func (d day12MoonSystem) updateVelocity() {
	for _, affected := range d {
		for _, affector := range d {
			if affector == affected {
				// Moons does not affect themselves
				continue
			}

			if affector.PX > affected.PX {
				affected.VX += 1
			} else if affector.PX < affected.PX {
				affected.VX -= 1
			}

			if affector.PY > affected.PY {
				affected.VY += 1
			} else if affector.PY < affected.PY {
				affected.VY -= 1
			}

			if affector.PZ > affected.PZ {
				affected.VZ += 1
			} else if affector.PZ < affected.PZ {
				affected.VZ -= 1
			}
		}
	}
}

type day12Moon struct {
	PX, PY, PZ int64
	VX, VY, VZ int64
}

func day12MoonFromScan(scan string) (*day12Moon, error) {
	grps := regexp.MustCompile(`^<x=([0-9-]+), y=([0-9-]+), z=([0-9-]+)>$`).FindStringSubmatch(scan)
	if len(grps) != 4 {
		return nil, errors.Errorf("Invalid input %q: %#v", scan, grps)
	}

	var (
		err error
		m   = &day12Moon{}
	)

	if m.PX, err = strconv.ParseInt(grps[1], 10, 64); err != nil {
		return nil, errors.Wrap(err, "Invalid X")
	}
	if m.PY, err = strconv.ParseInt(grps[2], 10, 64); err != nil {
		return nil, errors.Wrap(err, "Invalid Y")
	}
	if m.PZ, err = strconv.ParseInt(grps[3], 10, 64); err != nil {
		return nil, errors.Wrap(err, "Invalid Z")
	}

	return m, nil
}

func (d day12Moon) String() string {
	return fmt.Sprintf(
		"pos=<x=%d, y=%d, z=%d>, vel=<x=%d, y=%d, z=%d>",
		d.PX, d.PY, d.PZ,
		d.VX, d.VY, d.VZ,
	)
}

func (d *day12Moon) clearKinetic() { d.VX, d.VY, d.VZ = 0, 0, 0 }

func (d day12Moon) eq(t *day12Moon, component int) bool {
	switch component {
	case 0:
		return d.PX == t.PX && d.VX == t.VX
	case 1:
		return d.PY == t.PY && d.VY == t.VY
	case 2:
		return d.PZ == t.PZ && d.VZ == t.VZ
	default:
		panic("Invalid component compared")
	}
}

func (d *day12Moon) move() {
	d.PX += d.VX
	d.PY += d.VY
	d.PZ += d.VZ
}

func (d day12Moon) totalEnergy() int64 {
	return int64((math.Abs(float64(d.PX)) + math.Abs(float64(d.PY)) + math.Abs(float64(d.PZ))) * (math.Abs(float64(d.VX)) + math.Abs(float64(d.VY)) + math.Abs(float64(d.VZ))))
}

func solveDay12Part1(inFile string) (int64, error) {
	moons, err := day12MoonSystemFromFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read system scan")
	}

	moons.move(1000)

	return moons.totalEnergy(), nil
}

func solveDay12Part2(inFile string) (int64, error) {
	moons, err := day12MoonSystemFromFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read system scan")
	}

	initMoons, _ := day12MoonSystemFromFile(inFile)
	var (
		iterations [3]int64
		completed  [3]bool
	)

	for {
		moons.move(1)

		for c := 0; c < len(completed); c++ {
			if completed[c] {
				continue
			}

			var matches = true
			for i := range moons {
				if !moons[i].eq(initMoons[i], c) {
					matches = false
				}
			}

			iterations[c]++
			if matches {
				completed[c] = true
			}
		}

		var allComplete = true
		for _, c := range completed {
			if !c {
				allComplete = false
			}
		}

		if allComplete {
			break
		}
	}

	return leastCommonMultiple(leastCommonMultiple(iterations[0], iterations[1]), iterations[2]), nil
}
