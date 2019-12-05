package aoc2019

import (
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

func parseDay02Intcode(code string) ([]int, error) { return parseIntcode(code) }

func executeDay02Intcode(code []int) ([]int, error) {
	return executeIntcode(code, nil, nil) // Day02 intcode may not contain I/O
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
