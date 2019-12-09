package aoc2019

import (
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

func solveDay5FromFile(inFile string, diagProgram int64) (int64, error) {
	raw, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read input")
	}

	code, err := parseIntcode(strings.TrimSpace(string(raw)))
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse Intcode")
	}

	var (
		in      = make(chan int64, 1)
		out     = make(chan int64, len(code)) // There cannot be more outputs than number of code parts
		outputs []int64
	)

	/*
	 * The TEST diagnostic program will start by requesting from the user
	 * the ID of the system to test by running an input instruction - provide
	 * it 1, the ID for the ship's air conditioner unit.
	 */
	in <- diagProgram

	if _, err := executeIntcode(code, in, out); err != nil {
		return 0, errors.Wrap(err, "Program execution failed")
	}

	for o := range out {
		outputs = append(outputs, o)
	}

	if len(outputs) < 1 {
		return 0, errors.New("Program did not yield any output")
	}

	for i, check := range outputs[:len(outputs)-1] {
		if check != 0 {
			return 0, errors.Errorf("Non-zero check result in position %d: %d", i, check)
		}
	}

	return outputs[len(outputs)-1], nil
}

func solveDay5Part1(inFile string) (int64, error) {
	/*
	 * The TEST diagnostic program will start by requesting from the user
	 * the ID of the system to test by running an input instruction - provide
	 * it 1, the ID for the ship's air conditioner unit.
	 */
	return solveDay5FromFile(inFile, 1)
}

func solveDay5Part2(inFile string) (int64, error) {
	/*
	 * This time, when the TEST diagnostic program runs its input
	 * instruction to get the ID of the system to test, provide it 5,
	 * the ID for the ship's thermal radiator controller. This
	 * diagnostic test suite only outputs one number, the diagnostic code.
	 */
	return solveDay5FromFile(inFile, 5)
}
