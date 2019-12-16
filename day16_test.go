package aoc2019

import (
	"reflect"
	"testing"
)

func TestDay16TestPatternForElement(t *testing.T) {
	for idx, expPattern := range [][]int64{
		{0, 1, 0, -1},
		{0, 0, 1, 1, 0, 0, -1, -1},
		{0, 0, 0, 1, 1, 1, 0, 0, 0, -1, -1, -1},
	} {
		if p := day16PatternForElement(idx); !reflect.DeepEqual(expPattern, p.pattern) {
			t.Errorf("Unexpected pattern for idx=%d: exp=%+v got=%+v", idx, expPattern, p.pattern)
		}
	}
}

func TestDay16PatternLoop(t *testing.T) {
	p1 := day16PatternForElement(0)

	for i, expN := range []int64{
		// MUST skip first element at first execution
		// 0,
		1, 0, -1, 0, 1, 0, -1, 0, 1, 0, -1,
	} {
		if n := p1.get(); n != expN {
			t.Errorf("Unexpected number for pattern p1 at index %d: exp=%d got=%d", i, expN, n)
		}
	}
}

func TestDay16ParseInputSignal(t *testing.T) {
	s, err := day16ReadInputSignal("12345678")
	if err != nil {
		t.Fatalf("Unable to parse input signal: %s", err)
	}

	if !reflect.DeepEqual(s, []int64{1, 2, 3, 4, 5, 6, 7, 8}) {
		t.Errorf("Unexpected input signal: %+v", s)
	}
}

func TestDay16ProcessSignal(t *testing.T) {
	s, err := day16ReadInputSignal("12345678")
	if err != nil {
		t.Fatalf("Unable to parse input signal: %s", err)
	}

	for i, expS := range [][]int64{
		{1, 2, 3, 4, 5, 6, 7, 8},
		{4, 8, 2, 2, 6, 1, 5, 8},
		{3, 4, 0, 4, 0, 4, 3, 8},
		{0, 3, 4, 1, 5, 5, 1, 8},
		{0, 1, 0, 2, 9, 4, 9, 8},
	} {
		if rs := day16ProcessSignal(s, i); !reflect.DeepEqual(rs, expS) {
			t.Errorf("Unexpected processed signal after %d phases: exp=%+v got=%+v", i, expS, rs)
		}
	}
}

func TestDay16ProcessSignal8OfLonger(t *testing.T) {
	for signal, expPSig8 := range map[string][]int64{
		"80871224585914546619083218645595": {2, 4, 1, 7, 6, 1, 7, 6},
		"19617804207202209144916044189917": {7, 3, 7, 4, 5, 4, 1, 8},
		"69317163492948606335995924319873": {5, 2, 4, 3, 2, 1, 3, 3},
	} {
		s, err := day16ReadInputSignal(signal)
		if err != nil {
			t.Fatalf("Unable to parse input signal %q: %s", signal, err)
		}

		if rs := day16ProcessSignal(s, 100); !reflect.DeepEqual(rs[:8], expPSig8) {
			t.Errorf("Unexpected processed signal for signal %q: exp=%+v got=%+v", signal, expPSig8, rs[:8])
		}
	}
}

func TestCalculateDay16_Part1(t *testing.T) {
	res, err := solveDay16Part1("day16_input.txt")
	if err != nil {
		t.Fatalf("Day 16 solver failed: %s", err)
	}

	t.Logf("Solution Day 16 Part 1: %s", res)
}

func TestCalculateDay16_Part2(t *testing.T) {
	res, err := solveDay16Part2("day16_input.txt")
	if err != nil {
		t.Fatalf("Day 16 solver failed: %s", err)
	}

	t.Logf("Solution Day 16 Part 2: %s", res)
}
