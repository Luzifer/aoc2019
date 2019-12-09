package aoc2019

import (
	"io/ioutil"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

func solveDay9Part1(inFile string) (int64, error) {
	raw, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read input file")
	}

	code, err := parseIntcode(strings.TrimSpace(string(raw)))
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse intcode program")
	}

	var (
		inChan  = make(chan int64, 1)
		outChan = make(chan int64, 1)
		output  []int64
		wg      sync.WaitGroup
	)

	inChan <- 1
	wg.Add(1)

	go func() {
		for v := range outChan {
			output = append(output, v)
		}
		wg.Done()
	}()

	if _, err := executeIntcode(code, inChan, outChan); err != nil {
		return 0, errors.Wrap(err, "Unable to execute intcode")
	}

	wg.Wait()

	if len(output) != 1 {
		return 0, errors.Errorf("Got malfunction information: %+v", output)
	}

	return output[0], nil
}

func solveDay9Part2(inFile string) (int64, error) {
	raw, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read input file")
	}

	code, err := parseIntcode(strings.TrimSpace(string(raw)))
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse intcode program")
	}

	var (
		inChan  = make(chan int64, 1)
		outChan = make(chan int64, 1)
	)

	inChan <- 2

	if _, err := executeIntcode(code, inChan, outChan); err != nil {
		return 0, errors.Wrap(err, "Unable to execute intcode")
	}

	return <-outChan, nil
}
