package aoc2019

import (
	"math"
	"testing"
)

func TestDay18GetDistances(t *testing.T) {
	grid := day18ReadGrid([]byte("########################\n#...............b.C.D.f#\n#.######################\n#.....@.a.B.c.d.A.e.F.g#\n########################"))

	var maxInt int64 = math.MaxInt64

	dists, err := grid.getPathDistances("", &maxInt)
	if err != nil {
		t.Fatalf("GetPathDistances caused an error: %s", err)
	}

	if d := dists["bacdfeg"]; d != 132 {
		t.Errorf("Invalid distance for bacdfeg: exp=132 got=%d", d)
	}

	grid = day18ReadGrid([]byte(`#################
#i.G..c...e..H.p#
########.########
#j.A..b...f..D.o#
########@########
#k.E..a...g..B.n#
########.########
#l.F..d...h..C.m#
#################`))
	maxInt = math.MaxInt64

	dists, err = grid.getPathDistances("", &maxInt)
	if err != nil {
		t.Fatalf("GetPathDistances caused an error: %s", err)
	}

	if d := dists["afbjgnhdloepcikm"]; d != 136 {
		t.Errorf("Invalid distance for bacdfeg: exp=136 got=%d", d)
	}

}

func TestCalculateDay18_Part1(t *testing.T) {
	res, err := solveDay18Part1("day18_input.txt")
	if err != nil {
		t.Fatalf("Day 18 solver failed: %s", err)
	}

	t.Logf("Solution Day 18 Part 1: %d", res)
}

func TestCalculateDay18_Part2(t *testing.T) {
	res, err := solveDay18Part2("day18_input.txt")
	if err != nil {
		t.Fatalf("Day 18 solver failed: %s", err)
	}

	t.Logf("Solution Day 18 Part 2: %d", res)
}
