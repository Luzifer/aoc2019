package aoc2019

import (
	"bytes"
	"io/ioutil"
	"math"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var day16BasePattern = []int64{0, 1, 0, -1}

type day16Pattern struct {
	pattern []int64
	ptr     int
}

func (d *day16Pattern) get() int64 {
	d.ptr += 1
	if d.ptr == len(d.pattern) {
		d.ptr = 0
	}

	return d.pattern[d.ptr]
}

func day16PatternForElement(idx int) *day16Pattern {
	var pattern []int64

	for _, p := range day16BasePattern {
		for i := 0; i <= idx; i++ {
			pattern = append(pattern, p)
		}
	}

	return &day16Pattern{pattern: pattern, ptr: 0}
}

func day16ReadInputSignal(in string) ([]int64, error) {
	var out []int64

	for _, c := range in {
		v, err := strconv.ParseInt(string(c), 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to parse number")
		}
		out = append(out, v)
	}

	return out, nil
}

func day16ProcessSignal(s []int64, phases int) []int64 {
	var processedSignal = make([]int64, len(s))
	copy(processedSignal, s)

	for i := 0; i < phases; i++ {
		var tmpSignal = make([]int64, len(processedSignal))
		for i := range tmpSignal {
			p := day16PatternForElement(i)

			for _, n := range processedSignal {
				tmpSignal[i] += n * p.get()
			}

			tmpSignal[i] = int64(math.Abs(float64(tmpSignal[i]))) % 10
		}

		copy(processedSignal, tmpSignal)
	}

	return processedSignal
}

func solveDay16Part1(inFile string) (string, error) {
	raw, err := ioutil.ReadFile(inFile)
	if err != nil {
		return "", errors.Wrap(err, "Unable to read input file")
	}

	s, err := day16ReadInputSignal(strings.TrimSpace(string(raw)))
	if err != nil {
		return "", errors.Wrap(err, "Unable to parse input signal")
	}

	ps := day16ProcessSignal(s, 100)[:8]

	var res string
	for _, n := range ps {
		res += strconv.FormatInt(n, 10)
	}

	return res, nil
}

func solveDay16Part2(inFile string) (string, error) {
	raw, err := ioutil.ReadFile(inFile)
	if err != nil {
		return "", errors.Wrap(err, "Unable to read input file")
	}

	raw = bytes.TrimSpace(raw)

	// Assemble "real signal"
	rs := make([]byte, 10000*len(raw))
	for i := 0; i < 10000; i++ {
		copy(rs[i*len(raw):(i+1)*len(raw)], raw)
	}

	offset, err := strconv.ParseInt(string(raw[0:7]), 10, 64)
	if err != nil {
		return "", errors.Wrap(err, "Unable to parse offset")
	}
	rs = rs[offset:]

	s, err := day16ReadInputSignal(string(rs))
	if err != nil {
		return "", errors.Wrap(err, "Unable to parse input signal")
	}

	// Calculate special solution for offset-trimmed signal
	for phase := 0; phase < 100; phase++ {
		for i := len(s) - 1; i >= 0; i-- {
			var ni int64
			if i+1 < len(s) {
				ni = s[i+1]
			}
			s[i] = int64(math.Abs(float64(ni+s[i]))) % 10
		}
	}

	// Marshal output
	var res string
	for _, n := range s[:8] {
		res += strconv.FormatInt(n, 10)
	}

	return res, nil
}
