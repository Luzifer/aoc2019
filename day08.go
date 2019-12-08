package aoc2019

import (
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	day08ColorBlack int = iota
	day08ColorWhite
	day08ColorTransparent
)

type day08Layer struct {
	data []int
	w, h int
}

func (d day08Layer) countNumber(n int) int {
	var count int

	for _, v := range d.data {
		if v == n {
			count++
		}
	}

	return count
}

func day08ComposeLayers(layers []day08Layer) day08Layer {
	var finalLayerData = make([]int, layers[0].w*layers[0].h)

	for i := len(layers) - 1; i >= 0; i-- {
		for idx, col := range layers[i].data {
			if col == day08ColorTransparent {
				// Transparent, ignore
				continue
			}

			finalLayerData[idx] = col
		}
	}

	return day08Layer{data: finalLayerData, w: layers[0].w, h: layers[0].h}
}

func day08ParseToLayers(in string, w, h int) ([]day08Layer, error) {
	var layers []day08Layer

	if len(in)%(w*h) != 0 {
		return nil, errors.Errorf("Unexpected input length for format %dx%d: %d", w, h, len(in))
	}

	var pos int
	for pos < len(in) {
		var lData []int
		for _, c := range in[pos : pos+w*h] {
			v, err := strconv.Atoi(string(c))
			if err != nil {
				return nil, errors.Wrap(err, "Unable to parse digit")
			}
			lData = append(lData, v)
		}

		layers = append(layers, day08Layer{data: lData, w: w, h: h})
		pos += w * h
	}

	return layers, nil
}

func day08RenderLayer(layer day08Layer) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, layer.w, layer.h))

	for idx, px := range layer.data {
		var (
			x = idx % layer.w
			y = idx / layer.w
		)

		switch px {
		case day08ColorBlack:
			img.SetRGBA(x, y, color.RGBA{0x0, 0x0, 0x0, 0xff})
		case day08ColorWhite:
			img.SetRGBA(x, y, color.RGBA{0xff, 0xff, 0xff, 0xff})
		}
	}

	return img
}

func solveDay8Part1(inFile string) (int, error) {
	raw, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read input file")
	}

	layers, err := day08ParseToLayers(strings.TrimSpace(string(raw)), 25, 6)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse layers")
	}

	var (
		fewest = math.MaxInt64
		layer  day08Layer
	)
	for _, l := range layers {
		if n := l.countNumber(0); n < fewest {
			fewest = n
			layer = l
		}
	}

	return layer.countNumber(1) * layer.countNumber(2), nil
}

func solveDay8Part2(inFile string) error {
	raw, err := ioutil.ReadFile(inFile)
	if err != nil {
		return errors.Wrap(err, "Unable to read input file")
	}

	layers, err := day08ParseToLayers(strings.TrimSpace(string(raw)), 25, 6)
	if err != nil {
		return errors.Wrap(err, "Unable to parse layers")
	}

	img := day08RenderLayer(day08ComposeLayers(layers))

	f, err := os.Create("day08_image.png")
	if err != nil {
		return errors.Wrap(err, "Unable to open result file")
	}
	defer f.Close()

	return errors.Wrap(png.Encode(f, img), "Unable to render image")
}
