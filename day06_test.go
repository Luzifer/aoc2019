package aoc2019

import (
	"strings"
	"testing"
)

var day06OrbitMapExample = strings.TrimSpace(`
COM)B
B)C
C)D
D)E
E)F
B)G
G)H
D)I
E)J
J)K
K)L
K)YOU
I)SAN
`)

func TestDay6OrbitPathLengthToCOM(t *testing.T) {
	oMap, err := day06ParseOrbitMap(strings.NewReader(day06OrbitMapExample))
	if err != nil {
		t.Fatalf("Orbit map parser failed: %s", err)
	}

	for start, expCount := range map[string]int{
		"D":   3,
		"L":   7,
		"COM": 0,
	} {
		if c := oMap.getPathLengthToCOMByName(start); c != expCount {
			t.Errorf("Number of recursive orbits to %q yield unexpected result: exp=%d got=%d", start, expCount, c)
		}
	}
}

func TestDay6OrbitMapCommonPlanet(t *testing.T) {
	oMap, err := day06ParseOrbitMap(strings.NewReader(day06OrbitMapExample))
	if err != nil {
		t.Fatalf("Orbit map parser failed: %s", err)
	}

	p := oMap.getCommonPlanet(oMap["YOU"], oMap["SAN"])
	if p == nil {
		t.Fatalf("Found no common planet")
	}

	if p != oMap["D"] {
		t.Errorf("Found wrong common planet: exp=D, got=%s", p.Name)
	}
}

func TestCalculateDay6_Part1(t *testing.T) {
	count, err := solveDay6Part1("day06_input.txt")
	if err != nil {
		t.Fatalf("Day 6 solver failed: %s", err)
	}

	t.Logf("Solution Day 6 Part 1: %d", count)
}

func TestCalculateDay6_Part2(t *testing.T) {
	count, err := solveDay6Part2("day06_input.txt")
	if err != nil {
		t.Fatalf("Day 6 solver failed: %s", err)
	}

	t.Logf("Solution Day 6 Part 2: %d", count)
}
