package aoc2019

import "testing"

func TestCalculateDay5_Part1(t *testing.T) {
	result, err := solveDay5Part1("day05_input.txt")
	if err != nil {
		t.Fatalf("Day 5 solver failed: %s", err)
	}

	t.Logf("Solution Day 5 Part 1: %d", result)
}

func TestCalculateDay5_Part2(t *testing.T) {
	count, err := solveDay5Part2("day05_input.txt")
	if err != nil {
		t.Fatalf("Day 5 solver failed: %s", err)
	}

	t.Logf("Solution Day 5 Part 2: %d", count)
}
