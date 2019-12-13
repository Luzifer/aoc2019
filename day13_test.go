package aoc2019

import "testing"

func TestCalculateDay13_Part1(t *testing.T) {
	count, err := solveDay13Part1("day13_input.txt")
	if err != nil {
		t.Fatalf("Day 13 solver failed: %s", err)
	}

	t.Logf("Solution Day 13 Part 1: %d", count)
}

func TestCalculateDay13_Part2(t *testing.T) {
	res, err := solveDay13Part2("day13_input.txt")
	if err != nil {
		t.Fatalf("Day 13 solver failed: %s", err)
	}

	t.Logf("Solution Day 13 Part 2: %d", res)
}
