package aoc2019

import (
	"bufio"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type day14NanoFactory struct {
	reactions map[string]day14Reaction
}

func (d day14NanoFactory) calculateOreForFuel(n int64) int64 {
	var (
		availableMats = map[string]int64{}
		oreUsed       int64
	)

	var produce func(mat string, amount int64)
	produce = func(mat string, amount int64) {
		var (
			reaction = d.reactions[mat]
			reactN   = int64(math.Ceil(float64(amount) / float64(reaction.yield)))
		)

		for reqMat, reqAm := range reaction.materials {
			// Multiply for requested amount
			reqAm *= reactN

			if reqMat == "ORE" {
				oreUsed += reqAm
				continue
			}

			if availableMats[reqMat] < reqAm {
				produce(reqMat, reqAm-availableMats[reqMat])
			}

			availableMats[reqMat] -= reqAm
		}

		availableMats[mat] += reactN * reaction.yield
	}

	produce("FUEL", n)

	return oreUsed
}

type day14Reaction struct {
	materials map[string]int64
	yield     int64

	factory *day14NanoFactory
}

func day14ParseReactionChain(r io.Reader) (*day14NanoFactory, error) {
	var factory = &day14NanoFactory{reactions: make(map[string]day14Reaction)}

	parse := func(in string) (int64, string, error) {
		var parts = strings.Split(in, " ")
		count, err := strconv.ParseInt(parts[0], 10, 64)
		return count, parts[1], errors.Wrap(err, "Unable to parse count")
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var (
			sides  = strings.Split(scanner.Text(), " => ")
			inputs = strings.Split(sides[0], ", ")
		)

		yield, outMaterial, err := parse(sides[1])
		if err != nil {
			return nil, errors.Wrap(err, "Unable to parse output side")
		}

		rea := day14Reaction{factory: factory, yield: yield, materials: make(map[string]int64)}

		for _, in := range inputs {
			required, inMaterial, err := parse(in)
			if err != nil {
				return nil, errors.Wrap(err, "Unable to parse input material")
			}
			rea.materials[inMaterial] = required
		}

		factory.reactions[outMaterial] = rea
	}

	return factory, errors.Wrap(scanner.Err(), "Unable to scan input")
}

func solveDay14Part1(inFile string) (int64, error) {
	f, err := os.Open(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to open input file")
	}
	defer f.Close()

	factory, err := day14ParseReactionChain(f)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse reaction chain")
	}

	return factory.calculateOreForFuel(1), nil
}

func solveDay14Part2(inFile string) (int64, error) {
	f, err := os.Open(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to open input file")
	}
	defer f.Close()

	factory, err := day14ParseReactionChain(f)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse reaction chain")
	}

	var (
		limit      int64 = 1000000000000
		orePerFuel       = factory.calculateOreForFuel(1)
		minBound         = limit / orePerFuel
		maxBound         = minBound * 2
	)

	for {
		var (
			test    = (minBound + maxBound) / 2
			oreUsed = factory.calculateOreForFuel(test)
		)

		log.Printf("limit=%d used=%d min=%d max=%d test=%d", limit, oreUsed, minBound, maxBound, test)

		switch {
		case minBound+1 == maxBound && oreUsed < limit:
			return test, nil

		case oreUsed > limit:
			maxBound = test

		case oreUsed < limit:
			minBound = test

		case oreUsed == limit:
			// Doubt this will happen but well
			return test, nil
		}

	}
}
