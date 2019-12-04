package aoc2019

import (
	"math"
	"reflect"
	"testing"
)

func TestDay4NumberToDigitSlice(t *testing.T) {
	for n, expSlice := range map[int64][]int{
		1234567:  {1, 2, 3, 4, 5, 6, 7},
		43626145: {4, 3, 6, 2, 6, 1, 4, 5},
	} {
		if s := day4NumberToDigitSlice(n); !reflect.DeepEqual(expSlice, s) {
			t.Errorf("Number to slice for number %d yield unexpected result: exp=%+v got=%+v", n, expSlice, s)
		}
	}
}

func TestDay4ValidPassword(t *testing.T) {
	for n, expValid := range map[int64]bool{
		111111: true,
		223450: false,
		123789: false,
	} {
		if v := day4IsValidPassword(n, 0, math.MaxInt64); v != expValid {
			t.Errorf("Number %d did not have expected validity: exp=%v got=%v", n, expValid, v)
		}
	}

	for n, expValid := range map[int64]bool{
		112233: true,
		123444: false,
		111122: true,
	} {
		if v := day4IsValidPasswordPart2(n, 0, math.MaxInt64); v != expValid {
			t.Errorf("Number %d did not have expected validity for part 2: exp=%v got=%v", n, expValid, v)
		}
	}
}

func TestCalculateDay4_Part1(t *testing.T) {
	count, err := solveDay4Part1("day04_input.txt")
	if err != nil {
		t.Fatalf("Day 4 solver failed: %s", err)
	}

	t.Logf("Solution Day 4 Part 1: %d", count)
}

func TestCalculateDay4_Part2(t *testing.T) {
	count, err := solveDay4Part2("day04_input.txt")
	if err != nil {
		t.Fatalf("Day 4 solver failed: %s", err)
	}

	t.Logf("Solution Day 4 Part 2: %d", count)
}
