package aoc2019

import "testing"

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
