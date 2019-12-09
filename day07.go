package aoc2019

import (
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

func day07TestMaxOutputFromChain(code []int64, chainStart, chainEnd int, looped bool) int64 {
	var (
		chainLen = chainEnd - chainStart + 1
		permuts  [][]int64
		permute  func(a []int64, k int)
		rootSeq  = make([]int64, chainLen)
	)

	permute = func(a []int64, k int) {
		if k == len(a) {
			permuts = append(permuts, append([]int64{}, a...))
		} else {
			for i := k; i < len(rootSeq); i++ {
				a[k], a[i] = a[i], a[k]
				permute(a, k+1)
				a[k], a[i] = a[i], a[k]
			}
		}
	}

	for i := range rootSeq {
		rootSeq[i] = int64(chainStart + i)
	}
	permute(rootSeq, 0)

	var maxOutput int64
	for _, seq := range permuts {
		var (
			commInChans = make([]chan int64, chainLen)
			commOut     = make(chan int64, 2)
		)

		// Create channels
		for i := range commInChans {
			commInChans[i] = make(chan int64, 2)
		}

		// Build execution chain
		for i := 0; i < chainLen-1; i++ {
			go executeIntcode(cloneIntcode(code), commInChans[i], commInChans[i+1])
		}
		go executeIntcode(cloneIntcode(code), commInChans[chainLen-1], commOut)

		// Initialize chain
		for i, v := range seq {
			commInChans[i] <- v
		}
		commInChans[0] <- 0 // Input signal

		var lastOutput int64
		for r := range commOut {
			lastOutput = r
			if looped {
				commInChans[0] <- r
			}
		}

		// Test output of last execution
		if lastOutput > maxOutput {
			maxOutput = lastOutput
		}
	}

	return maxOutput
}

func solveDay7Part1(inFile string) (int64, error) {
	raw, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read input file")
	}

	code, err := parseIntcode(strings.TrimSpace(string(raw)))
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse intcode program")
	}

	return day07TestMaxOutputFromChain(code, 0, 4, false), nil
}

func solveDay7Part2(inFile string) (int64, error) {
	raw, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read input file")
	}

	code, err := parseIntcode(strings.TrimSpace(string(raw)))
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse intcode program")
	}

	return day07TestMaxOutputFromChain(code, 5, 9, true), nil
}
