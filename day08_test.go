package aoc2019

import (
	"reflect"
	"testing"
)

func TestDay08ComposeLayers(t *testing.T) {
	var (
		exp = []int{0, 1, 1, 0}
		in  = "0222112222120000"
		w   = 2
		h   = 2
	)

	layers, err := day08ParseToLayers(in, w, h)
	if err != nil {
		t.Fatalf("Layer parser failed: %s", err)
	}

	layer := day08ComposeLayers(layers)
	if !reflect.DeepEqual(layer.data, exp) {
		t.Errorf("Did not yield expected composed layer: exp=%+v got=%+v", exp, layer.data)
	}
}

func TestDay08ParseToLayers(t *testing.T) {
	var (
		in = "123456789012"
		w  = 3
		h  = 2
	)

	layers, err := day08ParseToLayers(in, w, h)
	if err != nil {
		t.Fatalf("Layer parser failed: %s", err)
	}

	if len(layers) != 2 || !reflect.DeepEqual(layers[0].data, []int{1, 2, 3, 4, 5, 6}) || !reflect.DeepEqual(layers[1].data, []int{7, 8, 9, 0, 1, 2}) {
		t.Errorf("Did not yield expected layers (123 456) (789 012): got=%#v", layers)
	}
}

func TestCalculateDay8_Part1(t *testing.T) {
	count, err := solveDay8Part1("day08_input.txt")
	if err != nil {
		t.Fatalf("Day 8 solver failed: %s", err)
	}

	t.Logf("Solution Day 8 Part 1: %d", count)
}

func TestCalculateDay8_Part2(t *testing.T) {
	if err := solveDay8Part2("day08_input.txt"); err != nil {
		t.Fatalf("Day 8 solver failed: %s", err)
	}

	t.Log("Solution Day 8 Part 2: See day08_image.png")
}
