package aoc2019

import (
	"io"
	"io/ioutil"
	"math"
	"os"
	"sort"

	"github.com/pkg/errors"
)

type day10MonitorGrid struct {
	asteroidMap []byte
	width       int
}

func (d day10MonitorGrid) asteroidCount() int {
	var count int
	for _, c := range d.asteroidMap {
		if c == '#' {
			count++
		}
	}

	return count
}

func (d day10MonitorGrid) clone() *day10MonitorGrid {
	var grid = &day10MonitorGrid{asteroidMap: make([]byte, len(d.asteroidMap)), width: d.width}
	for i, c := range d.asteroidMap {
		grid.asteroidMap[i] = c
	}
	return grid
}

func (d day10MonitorGrid) getAsteroidPositions() []int {
	var knownPositions []int
	for i, c := range d.asteroidMap {
		if c != '#' {
			// Not an asteroid, don't care
			continue
		}

		knownPositions = append(knownPositions, i)
	}
	return knownPositions
}

func (d day10MonitorGrid) getCleanedGrid(x, y int) *day10MonitorGrid {
	// Clone the map to work on
	var grid = d.clone()

	// Collect positions of all known asteroids
	var knownPositions = grid.getAsteroidPositions()

	// Mark observer (does not count into observable asteroids)
	grid.asteroidMap[grid.coordToPos(x, y)] = '@'

	// Iterate all positions and remove covered (invisible) asteroids
	for _, pos := range knownPositions {
		var aX, aY = d.posToCoord(pos)
		if grid.isObstructed(x, y, aX, aY) {
			grid.asteroidMap[pos] = '-'
		}
	}

	return grid
}

func (d *day10MonitorGrid) isObstructed(observX, observY, x, y int) bool {
	var distX, distY = x - observX, y - observY

	if distX == 0 && distY == 0 {
		// No steps, observer equals asteroid, needless calculation
		return false
	}

	var (
		div          = int(math.Abs(float64(greatestCommonDivisor(int64(distX), int64(distY)))))
		stepX, stepY = distX / div, distY / div
	)

	for i := 1; i < math.MaxInt64; i++ {
		var rPosX, rPosY = observX + stepX*i, observY + stepY*i

		if rPosX < 0 || rPosX >= d.width || rPosY < 0 || rPosY >= len(d.asteroidMap)/d.width {
			// Position outside grid, stop searching
			panic(errors.Errorf("Observed position ran out of bounds (obsX=%d obsY=%d x=%d y=%d div=%d stepX=%d stepY=%d)", observX, observY, x, y, div, stepX, stepY))
		}

		if rPosX == x && rPosY == y {
			return false
		}

		if d.asteroidMap[d.coordToPos(rPosX, rPosY)] == '#' {
			return true
		}
	}

	panic(errors.Errorf("Unreachable end was reached"))
}

func (d day10MonitorGrid) coordToPos(x, y int) int       { return y*d.width + x }
func (d day10MonitorGrid) posToCoord(pos int) (int, int) { return pos % d.width, pos / d.width }
func (d day10MonitorGrid) step2deg(x, y int) float64 {
	rad := math.Atan2(float64(x), float64(y))
	deg := rad * (180 / math.Pi)
	return 180 + -1*deg
}

func day10ReadAsteroidMap(in io.Reader) (*day10MonitorGrid, error) {
	var grid = &day10MonitorGrid{}

	raw, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read asteroid map")
	}

	for _, c := range raw {
		if c == '\n' {
			if grid.width == 0 {
				grid.width = len(grid.asteroidMap)
			}
			// Skip newlines for our representation
			continue
		}

		grid.asteroidMap = append(grid.asteroidMap, c)
	}

	return grid, nil
}

func solveDay10Part1Coordinate(inFile string) (*day10MonitorGrid, int, error) {
	f, err := os.Open(inFile)
	if err != nil {
		return nil, 0, errors.Wrap(err, "Unable to open input file")
	}
	defer f.Close()

	grid, err := day10ReadAsteroidMap(f)
	if err != nil {
		return nil, 0, errors.Wrap(err, "Unable to read asteroid map")
	}

	var (
		bestMonitorPos    int
		bestAsteroidCount int
	)
	for _, pos := range grid.getAsteroidPositions() {
		var aX, aY = grid.posToCoord(pos)
		rGrid := grid.getCleanedGrid(aX, aY)

		if c := rGrid.asteroidCount(); c > bestAsteroidCount {
			bestMonitorPos = pos
			bestAsteroidCount = c
		}
	}

	return grid, bestMonitorPos, nil
}

func solveDay10Part1(inFile string) (int, error) {
	grid, bestMonitorPos, err := solveDay10Part1Coordinate(inFile)
	if err != nil {
		return 0, err
	}

	var aX, aY = grid.posToCoord(bestMonitorPos)
	return grid.getCleanedGrid(aX, aY).asteroidCount(), nil
}

func solveDay10Part2(inFile string) (int, error) {
	grid, bestMonitorPos, err := solveDay10Part1Coordinate(inFile)
	if err != nil {
		return 0, err
	}

	var (
		mX, mY    = grid.posToCoord(bestMonitorPos)
		destroyed []int
	)

	// Mark monitor / laser -- cannot be destroyed
	grid.asteroidMap[bestMonitorPos] = 'M'

	// Gradually destroy asteroids
	for grid.asteroidCount() > 0 {
		asteroidsInSight := grid.getCleanedGrid(mX, mY).getAsteroidPositions()
		var degPos [][2]int
		for _, pos := range asteroidsInSight {
			var aX, aY = grid.posToCoord(pos)

			degPos = append(degPos, [2]int{
				pos,
				int(grid.step2deg(aX-mX, aY-mY) * 1000000), // Degree to asteroid in 6-digit precision
			})
		}

		// Sort by degree low-to-high -- represents order of destruction
		sort.Slice(degPos, func(i, j int) bool { return degPos[i][1] < degPos[j][1] })

		for _, dp := range degPos {
			grid.asteroidMap[dp[0]] = '*' // Mark asteroids destroyed
			destroyed = append(destroyed, dp[0])
		}
	}

	destr200X, destr200Y := grid.posToCoord(destroyed[199])

	return 100*destr200X + destr200Y, nil
}
