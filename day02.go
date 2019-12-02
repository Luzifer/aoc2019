package aoc2019

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func parseDay02Intcode(code string) ([]int, error) {
	parts := strings.Split(code, ",")

	var out []int
	for _, n := range parts {
		v, err := strconv.Atoi(n)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}

	return out, nil
}

func executeDay02Intcode(code []int) ([]int, error) {
	var (
		pos int
		run = true
	)

	for run {
		if pos >= len(code) {
			return nil, errors.Errorf("Code position out of bounds: %d (len=%d)", pos, len(code))
		}

		switch code[pos] {
		case 1:
			// addition
			p1, p2, pt := code[pos+1], code[pos+2], code[pos+3]
			code[pt] = code[p1] + code[p2]

		case 2:
			// multiplication
			p1, p2, pt := code[pos+1], code[pos+2], code[pos+3]
			code[pt] = code[p1] * code[p2]

		case 99:
			// program finished
			run = false

		default:
			return nil, errors.Errorf("Encountered invalid operation %d", code[pos])
		}

		pos += 4
	}

	return code, nil
}

func solveDay2Part1(inFile string) (int, error) {
	raw, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read input")
	}

	code, err := parseDay02Intcode(strings.TrimSpace(string(raw)))
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse Intcode")
	}

	// Modify data for "1202 program alarm"
	code[1] = 12
	code[2] = 2

	// Execute Intcode
	code, err = executeDay02Intcode(code)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to execute Intcode")
	}

	return code[0], nil
}

func solveDay2Part2(inFile string) (int, error) {
	raw, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read input")
	}

	const expectedResult int = 19690720

	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {

			// Re-initialize "memory"
			code, err := parseDay02Intcode(strings.TrimSpace(string(raw)))
			if err != nil {
				return 0, errors.Wrap(err, "Unable to parse Intcode")
			}

			// Modify code
			code[1] = noun
			code[2] = verb

			// Execute Intcode
			code, err = executeDay02Intcode(code)
			if err != nil {
				return 0, errors.Wrap(err, "Unable to execute Intcode")
			}

			if code[0] == expectedResult {
				return 100*noun + verb, nil
			}

		}
	}

	return 0, errors.New("No valid result was found")
}
