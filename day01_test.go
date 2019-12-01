package aoc2019

import "testing"

func TestCalculateDay1_Examples(t *testing.T) {
	for mass, expFuel := range map[int64]int64{
		12:     2,
		14:     2,
		1969:   654,
		100756: 33583,
	} {
		if f := calculateDay1FuelForMass(mass); f != expFuel {
			t.Errorf("Mismatch for mass of %d, expected %d, got %d", mass, expFuel, f)
		}
	}

	for mass, expFuel := range map[int64]int64{
		14:     2,
		1969:   966,
		100756: 50346,
	} {
		if f := calculateDay1FuelForMassRecurse(mass); f != expFuel {
			t.Errorf("Mismatch in recurse for mass of %d, expected %d, got %d", mass, expFuel, f)
		}
	}
}

func TestCalculateDay1_Part1(t *testing.T) {
	fuel, err := solveDay1Part1("day01_input.txt")
	if err != nil {
		t.Fatalf("Day 1 solver failed: %s", err)
	}

	t.Logf("Solution Day 1 Part 1: %d", fuel)
}

func TestCalculateDay1_Part2(t *testing.T) {
	fuel, err := solveDay1Part2("day01_input.txt")
	if err != nil {
		t.Fatalf("Day 1 solver failed: %s", err)
	}

	t.Logf("Solution Day 1 Part 2: %d", fuel)
}
