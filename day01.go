package aoc2019

import (
	"bufio"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

func calculateDay1FuelForMass(mass int64) int64 {
	return mass/3 - 2
}

func calculateDay1FuelForMassRecurse(mass int64) int64 {
	sumFuel := mass/3 - 2

	var addFuel = calculateDay1FuelForMass(sumFuel)
	for addFuel > 0 {
		sumFuel += addFuel
		addFuel = calculateDay1FuelForMass(addFuel)
	}

	return sumFuel
}

func solveDay1FromInput(inFile string, sumFn func(int64) int64) (int64, error) {
	f, err := os.Open(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to open input file")
	}
	defer f.Close()

	var sumFuel int64

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return 0, errors.Wrap(err, "Unable to scan file")
		}

		mass, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			return 0, errors.Wrapf(err, "Unable to parse integer input %q", scanner.Text())
		}

		sumFuel += sumFn(mass)
	}

	return sumFuel, nil
}

func solveDay1Part1(inFile string) (int64, error) {
	return solveDay1FromInput(inFile, calculateDay1FuelForMass)
}

func solveDay1Part2(inFile string) (int64, error) {
	return solveDay1FromInput(inFile, calculateDay1FuelForMassRecurse)
}
