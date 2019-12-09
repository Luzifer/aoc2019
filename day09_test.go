package aoc2019

import "testing"

func TestCalculateDay9_Part1(t *testing.T) {
	codeP0, err := solveDay9Part1("day09_input.txt")
	if err != nil {
		t.Fatalf("Day 9 solver failed: %s", err)
	}

	t.Logf("Solution Day 9 Part 1: %d", codeP0)
}

func TestCalculateDay9_Part2(t *testing.T) {
	result, err := solveDay9Part2("day09_input.txt")
	if err != nil {
		t.Fatalf("Day 9 solver failed: %s", err)
	}

	t.Logf("Solution Day 9 Part 2: %d", result)
}
