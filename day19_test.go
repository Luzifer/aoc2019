package aoc2019

import "testing"

func TestCalculateDay19_Part1(t *testing.T) {
	res, err := solveDay19Part1("day19_input.txt")
	if err != nil {
		t.Fatalf("Day 19 solver failed: %s", err)
	}

	t.Logf("Solution Day 19 Part 1: %d", res)
}

func TestCalculateDay19_Part2(t *testing.T) {
	res, err := solveDay19Part2("day19_input.txt")
	if err != nil {
		t.Fatalf("Day 19 solver failed: %s", err)
	}

	t.Logf("Solution Day 19 Part 2: %d", res)
}
