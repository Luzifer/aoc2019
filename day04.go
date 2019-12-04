package aoc2019

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func day4NumberToDigitSlice(n int64) []int {
	var out []int
	sn := strconv.FormatInt(n, 10)
	for _, c := range sn {
		d, err := strconv.Atoi(string(c))
		if err != nil {
			// This should not happen as the input is a guaranteed base-10 number
			panic(err)
		}
		out = append(out, d)
	}

	return out
}

func day4IsValidPassword(n, min, max int64) bool {
	sn := day4NumberToDigitSlice(n)

	// It is a six-digit number.
	if len(sn) != 6 {
		return false
	}

	// The value is within the range given in your puzzle input.
	if n < min || n > max {
		return false
	}

	// Two adjacent digits are the same (like 22 in 122345).
	var foundAdjacentMatch bool
	for i := 1; i < len(sn); i++ {
		if sn[i] == sn[i-1] {
			foundAdjacentMatch = true
		}
	}
	if !foundAdjacentMatch {
		return false
	}

	// Going from left to right, the digits never decrease; they only ever increase or stay the same (like 111123 or 135679).
	for i := 1; i < len(sn); i++ {
		if sn[i] < sn[i-1] {
			return false
		}
	}

	return true
}

func day4IsValidPasswordPart2(in, min, max int64) bool {
	// Previous rules still apply
	if !day4IsValidPassword(in, min, max) {
		return false
	}

	sn := day4NumberToDigitSlice(in)

	var (
		last  = sn[0]
		count = 1
	)

	for i := 1; i < len(sn); i++ {
		if sn[i] != last && count == 2 {
			return true
		}

		if sn[i] == last {
			count++
			continue
		}

		last = sn[i]
		count = 1
	}

	return count == 2
}

func solveDay4WithFunction(inFile string, vf func(int64, int64, int64) bool) (int, error) {
	raw, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read input file")
	}

	pts := strings.Split(strings.TrimSpace(string(raw)), "-")
	min, err := strconv.ParseInt(pts[0], 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse minimum")
	}

	max, err := strconv.ParseInt(pts[1], 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse maximum")
	}

	var count int
	for i := min; i <= max; i++ {
		if vf(i, min, max) {
			count++
		}
	}

	return count, nil
}

func solveDay4Part1(inFile string) (int, error) {
	return solveDay4WithFunction(inFile, day4IsValidPassword)
}

func solveDay4Part2(inFile string) (int, error) {
	return solveDay4WithFunction(inFile, day4IsValidPasswordPart2)
}
