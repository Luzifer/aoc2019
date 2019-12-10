package aoc2019

import "testing"

func TestGreatestCommonDivisor(t *testing.T) {
	for expDivisor, vals := range map[int64][2]int64{
		1:  {2, 3},
		2:  {4, 6},
		-5: {-5, 15},
	} {
		if r := greatestCommonDivisor(vals[0], vals[1]); r != expDivisor {
			t.Errorf("Unexpected divisor for values %+v: exp=%d got=%d", vals, expDivisor, r)
		}
	}
}

func TestManhattenDistance(t *testing.T) {
	for expDist, p := range map[int][4]int{
		2: {0, 0, 0, 2},
		3: {0, 0, 3, 0},
		6: {-3, 0, 3, 0},
		7: {-3, 0, 0, 4},
	} {
		if dist := manhattenDistance(p[0], p[1], p[2], p[3]); dist != expDist {
			t.Errorf("Unexpected distance for %+v: exp=%d got=%d", p, expDist, dist)
		}
	}
}
