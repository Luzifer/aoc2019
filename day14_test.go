package aoc2019

import (
	"strings"
	"testing"
)

func TestDay14OreForFuel(t *testing.T) {
	in := strings.NewReader("10 ORE => 10 A\n1 ORE => 1 B\n7 A, 1 B => 1 C\n7 A, 1 C => 1 D\n7 A, 1 D => 1 E\n7 A, 1 E => 1 FUEL")
	factory, err := day14ParseReactionChain(in)
	if err != nil {
		t.Fatalf("Parsing reaction chain caused an error: %s", err)
	}

	if ore := factory.calculateOreForFuel(1); ore != 31 {
		t.Errorf("Unexpected amount of ore: exp=31 got=%d", ore)
	}

	// More complex example
	in = strings.NewReader(`171 ORE => 8 CNZTR
7 ZLQW, 3 BMBT, 9 XCVML, 26 XMNCP, 1 WPTQ, 2 MZWV, 1 RJRHP => 4 PLWSL
114 ORE => 4 BHXH
14 VRPVC => 6 BMBT
6 BHXH, 18 KTJDG, 12 WPTQ, 7 PLWSL, 31 FHTLT, 37 ZDVW => 1 FUEL
6 WPTQ, 2 BMBT, 8 ZLQW, 18 KTJDG, 1 XMNCP, 6 MZWV, 1 RJRHP => 6 FHTLT
15 XDBXC, 2 LTCX, 1 VRPVC => 6 ZLQW
13 WPTQ, 10 LTCX, 3 RJRHP, 14 XMNCP, 2 MZWV, 1 ZLQW => 1 ZDVW
5 BMBT => 4 WPTQ
189 ORE => 9 KTJDG
1 MZWV, 17 XDBXC, 3 XCVML => 2 XMNCP
12 VRPVC, 27 CNZTR => 2 XDBXC
15 KTJDG, 12 BHXH => 5 XCVML
3 BHXH, 2 VRPVC => 7 MZWV
121 ORE => 7 VRPVC
7 XCVML => 6 RJRHP
5 BHXH, 4 VRPVC => 5 LTCX`)

	factory, err = day14ParseReactionChain(in)
	if err != nil {
		t.Fatalf("Parsing complex reaction chain caused an error: %s", err)
	}

	if ore := factory.calculateOreForFuel(1); ore != 2210736 {
		t.Errorf("Unexpected amount of ore: exp=2210736 got=%d", ore)
	}
}

func TestCalculateDay14_Part1(t *testing.T) {
	count, err := solveDay14Part1("day14_input.txt")
	if err != nil {
		t.Fatalf("Day 14 solver failed: %s", err)
	}

	t.Logf("Solution Day 14 Part 1: %d", count)
}

func TestCalculateDay14_Part2(t *testing.T) {
	res, err := solveDay14Part2("day14_input.txt")
	if err != nil {
		t.Fatalf("Day 14 solver failed: %s", err)
	}

	t.Logf("Solution Day 14 Part 2: %d", res)
}
