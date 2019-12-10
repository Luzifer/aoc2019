package aoc2019

import (
	"strings"
	"testing"
)

func TestDay10ReadAsteroidMap(t *testing.T) {
	grid, err := day10ReadAsteroidMap(strings.NewReader(".#..#\n.....\n#####\n....#\n...##"))
	if err != nil {
		t.Fatalf("Asteroid map parser failed: %s", err)
	}

	if grid.width != 5 {
		t.Errorf("Wrong width detected: exp=5 got=%d", grid.width)
	}

	if l := len(grid.asteroidMap); l != 25 {
		t.Errorf("Wrong length of asteroid map detected: exp=25 got=%d", l)
	}

	if c := grid.asteroidCount(); c != 10 {
		t.Errorf("Wrong number of asteroids detected: exp=10 got=%d", c)
	}
}

func TestDay10GetCleanedGrid(t *testing.T) {
	grid, err := day10ReadAsteroidMap(strings.NewReader(".#..#\n.....\n#####\n....#\n...##"))
	if err != nil {
		t.Fatalf("Asteroid map parser failed: %s", err)
	}

	for expCount, coords := range map[int][2]int{
		5: {4, 2},
		6: {0, 2},
		7: {1, 0},
		8: {3, 4},
	} {
		rGrid := grid.getCleanedGrid(coords[0], coords[1])

		if c := rGrid.asteroidCount(); c != expCount {
			t.Errorf("Wrong number of asteroids detected: exp=%d got=%d (grid=%s)", expCount, c, rGrid.asteroidMap)
		}
	}
}

func TestDay10StepToDeg(t *testing.T) {
	// Is a function on a day10MonitorGrid, grid itself is not used
	grid, err := day10ReadAsteroidMap(strings.NewReader(".#..#\n.....\n#####\n....#\n...##"))
	if err != nil {
		t.Fatalf("Asteroid map parser failed: %s", err)
	}

	for expDeg, vals := range map[float64][2]int{
		0:   {0, -5},
		45:  {5, -5},
		90:  {5, 0},
		135: {5, 5},
		180: {0, 5},
		225: {-5, 5},
		270: {-5, 0},
		315: {-5, -5},
	} {
		if d := grid.step2deg(vals[0], vals[1]); d != expDeg {
			t.Errorf("Step to Degree yield unexpected result for %+v: exp=%.2f got=%.2f", vals, expDeg, d)
		}
	}
}

func TestCalculateDay10_Part1(t *testing.T) {
	count, err := solveDay10Part1("day10_input.txt")
	if err != nil {
		t.Fatalf("Day 10 solver failed: %s", err)
	}

	t.Logf("Solution Day 10 Part 1: %d", count)
}

func TestCalculateDay10_Part2(t *testing.T) {
	res, err := solveDay10Part2("day10_input.txt")
	if err != nil {
		t.Fatalf("Day 10 solver failed: %s", err)
	}

	t.Logf("Solution Day 10 Part 2: %d", res)
}
