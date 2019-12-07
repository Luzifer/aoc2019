package aoc2019

import "testing"

func TestChainedInput(t *testing.T) {
	code, err := parseIntcode("3,15,3,16,1002,16,10,16,1,16,15,15,4,15,99,0,0")
	if err != nil {
		t.Fatalf("Intcode parser failed: %s", err)
	}

	// Create channels
	var (
		aIn  = make(chan int, 2)
		bIn  = make(chan int, 2)
		cIn  = make(chan int, 2)
		dIn  = make(chan int, 2)
		eIn  = make(chan int, 2)
		eOut = make(chan int, 2)
	)

	// Build execution chain
	go executeIntcode(cloneIntcode(code), aIn, bIn)
	go executeIntcode(cloneIntcode(code), bIn, cIn)
	go executeIntcode(cloneIntcode(code), cIn, dIn)
	go executeIntcode(cloneIntcode(code), dIn, eIn)
	go executeIntcode(cloneIntcode(code), eIn, eOut)

	// Initialize chain
	aIn <- 4 // Sequence
	aIn <- 0 // Input signal
	bIn <- 3 // Sequence
	cIn <- 2 // Sequence
	dIn <- 1 // Sequence
	eIn <- 0 // Sequence

	// Test output of last execution
	if r := <-eOut; r != 43210 {
		t.Errorf("Unexpected result from chain: exp=43210 got=%d", r)
	}
}

func TestMaxOutputFromChain(t *testing.T) {
	for codeStr, expValue := range map[string]int{
		"3,15,3,16,1002,16,10,16,1,16,15,15,4,15,99,0,0":                                                      43210,
		"3,23,3,24,1002,24,10,24,1002,23,-1,23,101,5,23,23,1,24,23,23,4,23,99,0,0":                            54321,
		"3,31,3,32,1002,32,10,32,1001,31,-2,31,1007,31,0,33,1002,33,7,33,1,33,31,31,1,32,31,31,4,31,99,0,0,0": 65210,
	} {
		code, err := parseIntcode(codeStr)
		if err != nil {
			t.Fatalf("Parsing Intcode failed: %s", err)
		}

		if r := day07TestMaxOutputFromChain(code, 0, 4, false); r != expValue {
			t.Errorf("Max output yield unexpected result: exp=%d got=%d", expValue, r)
		}
	}
}

func TestMaxOutputFromLoopedChain(t *testing.T) {
	for codeStr, expValue := range map[string]int{
		"3,26,1001,26,-4,26,3,27,1002,27,2,27,1,27,26,27,4,27,1001,28,-1,28,1005,28,6,99,0,0,5":                                                                                         139629729,
		"3,52,1001,52,-5,52,3,53,1,52,56,54,1007,54,5,55,1005,55,26,1001,54,-5,54,1105,1,12,1,53,54,53,1008,54,0,55,1001,55,1,55,2,53,55,53,4,53,1001,56,-1,56,1005,56,6,99,0,0,0,0,10": 18216,
	} {
		code, err := parseIntcode(codeStr)
		if err != nil {
			t.Fatalf("Parsing Intcode failed: %s", err)
		}

		if r := day07TestMaxOutputFromChain(code, 5, 9, true); r != expValue {
			t.Errorf("Max output yield unexpected result: exp=%d got=%d", expValue, r)
		}
	}
}

func TestCalculateDay7_Part1(t *testing.T) {
	codeP0, err := solveDay7Part1("day07_input.txt")
	if err != nil {
		t.Fatalf("Day 7 solver failed: %s", err)
	}

	t.Logf("Solution Day 7 Part 1: %d", codeP0)
}

func TestCalculateDay7_Part2(t *testing.T) {
	result, err := solveDay7Part2("day07_input.txt")
	if err != nil {
		t.Fatalf("Day 7 solver failed: %s", err)
	}

	t.Logf("Solution Day 7 Part 2: %d", result)
}
