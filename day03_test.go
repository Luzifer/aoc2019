package aoc2019

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseDay3LineDefinition(t *testing.T) {
	for ld, expLine := range map[string]day3Line{
		"R5":          {{0, 0, 5, 0}},
		"R8,U5,L5,D3": {{0, 0, 8, 0}, {8, 0, 8, 5}, {8, 5, 3, 5}, {3, 5, 3, 2}},
		"U7,R6,D4,L4": {{0, 0, 0, 7}, {0, 7, 6, 7}, {6, 7, 6, 3}, {6, 3, 2, 3}},
	} {
		l, err := parseDay3LineDefinition(ld, 0, 0)
		if err != nil {
			t.Fatalf("Day 3 line parser failed: %s", err)
		}

		if !reflect.DeepEqual(expLine, l) {
			t.Errorf("Mismatch in line of definition %q: exp=%+v got=%+v", ld, expLine, l)
		}
	}
}

func TestDay3LineGetIntersections(t *testing.T) {
	for ld, expInter := range map[string][][2]int{
		"R8,U5,L5,D3|U7,R6,D4,L4": {{0, 0}, {6, 5}, {3, 3}},
		"R75,D30,R83,U83,L12,D49,R71,U7,L72|U62,R66,U55,R34,D71,R55,D58,R83": {{0, 0}, {158, -12}, {146, 46}, {155, 4}, {155, 11}},
	} {
		lds := strings.Split(ld, "|")
		l1, _ := parseDay3LineDefinition(lds[0], 0, 0)
		l2, _ := parseDay3LineDefinition(lds[1], 0, 0)

		inter := l1.GetIntersections(l2)
		if !reflect.DeepEqual(expInter, inter) {
			t.Errorf("Mismatch in line inters of def %q: exp=%+v got=%+v", ld, expInter, inter)
		}
	}
}

func TestDay3MinIntersectionDistance(t *testing.T) {
	for ld, expDist := range map[string]int{
		"R8,U5,L5,D3|U7,R6,D4,L4": 6,
		"R75,D30,R83,U83,L12,D49,R71,U7,L72|U62,R66,U55,R34,D71,R55,D58,R83":               159,
		"R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51|U98,R91,D20,R16,D67,R40,U7,R15,U6,R7": 135,
	} {
		lds := strings.Split(ld, "|")
		l1, _ := parseDay3LineDefinition(lds[0], 0, 0)
		l2, _ := parseDay3LineDefinition(lds[1], 0, 0)

		if dist := getDay3MinIntersectionDistance(l1, l2, 0, 0); dist != expDist {
			t.Errorf("Mismatch in intersection distance of def %q: exp=%d got=%d", ld, expDist, dist)
		}
	}
}

func TestDay3MinIntersectionSteps(t *testing.T) {
	for ld, expDist := range map[string]int{
		"R8,U5,L5,D3|U7,R6,D4,L4": 30,
		"R75,D30,R83,U83,L12,D49,R71,U7,L72|U62,R66,U55,R34,D71,R55,D58,R83":               610,
		"R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51|U98,R91,D20,R16,D67,R40,U7,R15,U6,R7": 410,
	} {
		lds := strings.Split(ld, "|")
		l1, _ := parseDay3LineDefinition(lds[0], 0, 0)
		l2, _ := parseDay3LineDefinition(lds[1], 0, 0)

		if dist := getDay3MinIntersectionSteps(l1, l2); dist != expDist {
			t.Errorf("Mismatch in intersection steps of def %q: exp=%d got=%d", ld, expDist, dist)
		}
	}
}

func TestDay3DebugRender(t *testing.T) {
	var ld = "R75,D30,R83,U83,L12,D49,R71,U7,L72|U62,R66,U55,R34,D71,R55,D58,R83"

	lds := strings.Split(ld, "|")
	l1, _ := parseDay3LineDefinition(lds[0], 0, 0)
	l2, _ := parseDay3LineDefinition(lds[1], 0, 0)

	if err := debugDay3DrawImage(l1, l2); err != nil {
		t.Fatalf("Day 3 debug failed: %s", err)
	}
}

func TestCalculateDay3_Part1(t *testing.T) {
	codeP0, err := solveDay3Part1("day03_input.txt")
	if err != nil {
		t.Fatalf("Day 3 solver failed: %s", err)
	}

	t.Logf("Solution Day 3 Part 1: %d", codeP0)
}

func TestCalculateDay3_Part2(t *testing.T) {
	codeP0, err := solveDay3Part2("day03_input.txt")
	if err != nil {
		t.Fatalf("Day 3 solver failed: %s", err)
	}

	t.Logf("Solution Day 3 Part 2: %d", codeP0)
}
