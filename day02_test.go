package aoc2019

import (
	"reflect"
	"testing"
)

func TestExecuteDay02Intcode(t *testing.T) {
	for codeStr, expResult := range map[string][]int{
		"1,0,0,0,99":          {2, 0, 0, 0, 99},
		"2,3,0,3,99":          {2, 3, 0, 6, 99},
		"2,4,4,5,99,0":        {2, 4, 4, 5, 99, 9801},
		"1,1,1,4,99,5,6,0,99": {30, 1, 1, 4, 2, 5, 6, 0, 99},
	} {
		code, err := parseDay02Intcode(codeStr)
		if err != nil {
			t.Fatalf("Parsing Intcode failed: %s", err)
		}

		res, err := executeDay02Intcode(code)
		if err != nil {
			t.Fatalf("Intcode execution failed: %s", err)
		}

		if !reflect.DeepEqual(res, expResult) {
			t.Errorf("Intcode execution yield unexpected result: %+v != %+v", res, expResult)
		}
	}
}

func TestCalculateDay2_Part1(t *testing.T) {
	codeP0, err := solveDay2Part1("day02_input.txt")
	if err != nil {
		t.Fatalf("Day 2 solver failed: %s", err)
	}

	t.Logf("Solution Day 2 Part 1: %d", codeP0)
}

func TestCalculateDay2_Part2(t *testing.T) {
	result, err := solveDay2Part2("day02_input.txt")
	if err != nil {
		t.Fatalf("Day 2 solver failed: %s", err)
	}

	t.Logf("Solution Day 2 Part 2: %d", result)
}
