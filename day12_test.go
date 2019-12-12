package aoc2019

import (
	"strings"
	"testing"
)

func TestDay12MoonExample(t *testing.T) {
	moons, err := day12MoonSystemFromReader(strings.NewReader("<x=-1, y=0, z=2>\n<x=2, y=-10, z=-7>\n<x=4, y=-8, z=8>\n<x=3, y=5, z=-1>"))
	if err != nil {
		t.Fatalf("Unable to parse moon system: %s", err)
	}

	if l := len(moons); l != 4 {
		t.Fatalf("Unexpected moon count: exp=4 got=%d", l)
	}

	// Check initial state
	for i, expStr := range []string{
		"pos=<x=-1, y=0, z=2>, vel=<x=0, y=0, z=0>",
		"pos=<x=2, y=-10, z=-7>, vel=<x=0, y=0, z=0>",
		"pos=<x=4, y=-8, z=8>, vel=<x=0, y=0, z=0>",
		"pos=<x=3, y=5, z=-1>, vel=<x=0, y=0, z=0>",
	} {
		if s := moons[i].String(); s != expStr {
			t.Errorf("Unexpected moon %d: exp=%q got=%q", i, expStr, s)
		}
	}

	// Move moons by one step
	moons.move(1)

	// Check move state
	for i, expStr := range []string{
		"pos=<x=2, y=-1, z=1>, vel=<x=3, y=-1, z=-1>",
		"pos=<x=3, y=-7, z=-4>, vel=<x=1, y=3, z=3>",
		"pos=<x=1, y=-7, z=5>, vel=<x=-3, y=1, z=-3>",
		"pos=<x=2, y=2, z=0>, vel=<x=-1, y=-3, z=1>",
	} {
		if s := moons[i].String(); s != expStr {
			t.Errorf("Unexpected moon %d: exp=%q got=%q", i, expStr, s)
		}
	}

	// Move moons by nine more step
	moons.move(9)

	// Check move state
	for i, expStr := range []string{
		"pos=<x=2, y=1, z=-3>, vel=<x=-3, y=-2, z=1>",
		"pos=<x=1, y=-8, z=0>, vel=<x=-1, y=1, z=3>",
		"pos=<x=3, y=-6, z=1>, vel=<x=3, y=2, z=-3>",
		"pos=<x=2, y=0, z=4>, vel=<x=1, y=-1, z=-1>",
	} {
		if s := moons[i].String(); s != expStr {
			t.Errorf("Unexpected moon %d: exp=%q got=%q", i, expStr, s)
		}
	}

	// Check energy after step 10
	for i, expE := range []int64{
		36, 45, 80, 18,
	} {
		if e := moons[i].totalEnergy(); e != expE {
			t.Errorf("Unexpected energy on moon %d: exp=%d got=%d", i, expE, e)
		}
	}
}

func TestCalculateDay12_Part1(t *testing.T) {
	count, err := solveDay12Part1("day12_input.txt")
	if err != nil {
		t.Fatalf("Day 12 solver failed: %s", err)
	}

	t.Logf("Solution Day 12 Part 1: %d", count)
}

func TestCalculateDay12_Part2(t *testing.T) {
	res, err := solveDay12Part2("day12_input.txt")
	if err != nil {
		t.Fatalf("Day 12 solver failed: %s", err)
	}

	t.Logf("Solution Day 12 Part 2: %d", res)
}
