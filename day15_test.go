package aoc2019

import "testing"

func TestCalculateDay15_Part1(t *testing.T) {
	count, err := solveDay15Part1("day15_input.txt")
	if err != nil {
		t.Fatalf("Day 15 solver failed: %s", err)
	}

	t.Logf("Solution Day 15 Part 1: %d", count)
}

func TestCalculateDay15_Part2(t *testing.T) {
	res, err := solveDay15Part2("day15_input.txt")
	if err != nil {
		t.Fatalf("Day 15 solver failed: %s", err)
	}

	t.Logf("Solution Day 15 Part 2: %d", res)
}
