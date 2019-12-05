package aoc2019

import "testing"

func TestParseOpCode(t *testing.T) {
	for code, expOpCode := range map[int]opCode{
		1002: {Type: opCodeTypeMultiplication, flags: []opCodeFlag{opCodeFlagPosition, opCodeFlagImmediate}},
		1101: {Type: opCodeTypeAddition, flags: []opCodeFlag{opCodeFlagImmediate, opCodeFlagImmediate}},
	} {
		if op := parseOpCode(code); !op.eq(expOpCode) {
			t.Errorf("OpCode execution of code %d yield unexpected result: exp=%+v got=%+v", code, expOpCode, op)
		}
	}
}

func TestExecuteIntcodeIO(t *testing.T) {
	code, _ := parseIntcode("3,0,4,0,99")

	var (
		exp = 25
		in  = make(chan int, 1)
		out = make(chan int, 1)
	)

	in <- exp

	if _, err := executeIntcode(code, in, out); err != nil {
		t.Fatalf("Intcode execution failed: %s", err)
	}

	if r := <-out; r != exp {
		t.Errorf("Program yield unexpected result: exp=%d got=%d", exp, r)
	}
}

func TestExecuteIntcodeImmediateFlag(t *testing.T) {
	// 102,4,7,0 = Multiply 4 by pos_7, store to pos_0
	// 4,0       = Output pos_0
	// 99        = Exit
	// 3         = pos_7
	code, _ := parseIntcode("102,4,7,0,4,0,99,3")

	var (
		exp = 12
		out = make(chan int, 1)
	)

	if _, err := executeIntcode(code, nil, out); err != nil {
		t.Fatalf("Intcode execution failed: %s", err)
	}

	if r := <-out; r != exp {
		t.Errorf("Program yield unexpected result: exp=%d got=%d", exp, r)
	}
}

func TestExecuteIntcodeEquals(t *testing.T) {
	for mode, codeStr := range map[string]string{
		"position":  "3,9,8,9,10,9,4,9,99,-1,8",
		"immediate": "3,3,1108,-1,8,3,4,3,99",
	} {

		for input, exp := range map[int]int{
			1:  0,
			8:  1,
			20: 0,
			-8: 0,
		} {
			var (
				in  = make(chan int, 1)
				out = make(chan int, 10)
			)

			code, _ := parseIntcode(codeStr)
			in <- input

			if _, err := executeIntcode(code, in, out); err != nil {
				t.Fatalf("Execute in mode %q with input %d caused an error: %s", mode, input, err)
			}

			if r := <-out; r != exp {
				t.Errorf("Execute in mode %q with input %d yield unexpected result: exp=%d got=%d", mode, input, exp, r)
			}

		}

	}
}

func TestExecuteIntcodeLessThan(t *testing.T) {
	for mode, codeStr := range map[string]string{
		"position":  "3,9,7,9,10,9,4,9,99,-1,8",
		"immediate": "3,3,1107,-1,8,3,4,3,99",
	} {

		for input, exp := range map[int]int{
			1:  1,
			8:  0,
			20: 0,
			-8: 1,
		} {
			var (
				in  = make(chan int, 1)
				out = make(chan int, 10)
			)

			code, _ := parseIntcode(codeStr)
			in <- input

			if _, err := executeIntcode(code, in, out); err != nil {
				t.Fatalf("Execute in mode %q with input %d caused an error: %s", mode, input, err)
			}

			if r := <-out; r != exp {
				t.Errorf("Execute in mode %q with input %d yield unexpected result: exp=%d got=%d", mode, input, exp, r)
			}

		}

	}
}

func TestExecuteIntcodeJump(t *testing.T) {
	for mode, codeStr := range map[string]string{
		"position":  "3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9",
		"immediate": "3,3,1105,-1,9,1101,0,0,12,4,12,99,1",
	} {

		for input, exp := range map[int]int{
			5: 1,
			0: 0,
		} {
			var (
				in  = make(chan int, 1)
				out = make(chan int, 10)
			)

			code, _ := parseIntcode(codeStr)
			in <- input

			if _, err := executeIntcode(code, in, out); err != nil {
				t.Fatalf("Execute in mode %q with input %d caused an error: %s", mode, input, err)
			}

			if r := <-out; r != exp {
				t.Errorf("Execute in mode %q with input %d yield unexpected result: exp=%d got=%d", mode, input, exp, r)
			}

		}

	}
}
