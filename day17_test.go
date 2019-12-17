package aoc2019

import "testing"

func TestCalculateDay17_Part1(t *testing.T) {
	res, err := solveDay17Part1("day17_input.txt")
	if err != nil {
		t.Fatalf("Day 17 solver failed: %s", err)
	}

	t.Logf("Solution Day 17 Part 1: %d", res)
}

func TestCalculateDay17_Part2(t *testing.T) {
	res, err := solveDay17Part2("day17_input.txt")
	if err != nil {
		t.Fatalf("Day 17 solver failed: %s", err)
	}

	t.Logf("Solution Day 17 Part 2: %d", res)
}
