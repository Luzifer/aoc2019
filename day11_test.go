package aoc2019

import "testing"

func TestCalculateDay11_Part1(t *testing.T) {
	count, err := solveDay11Part1("day11_input.txt")
	if err != nil {
		t.Fatalf("Day 11 solver failed: %s", err)
	}

	t.Logf("Solution Day 11 Part 1: %d", count)
}

func TestCalculateDay11_Part2(t *testing.T) {
	err := solveDay11Part2("day11_input.txt")
	if err != nil {
		t.Fatalf("Day 11 solver failed: %s", err)
	}

	t.Log("Solution Day 11 Part 2: See day11_image.png")
}
