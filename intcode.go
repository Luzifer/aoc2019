package aoc2019

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type opCodeFlag int

const (
	opCodeFlagPosition opCodeFlag = iota
	opCodeFlagImmediate
)

type opCodeType int

const (
	opCodeTypeAddition       opCodeType = 1  // Day 02
	opCodeTypeMultiplication opCodeType = 2  // Day 02
	opCodeTypeInput          opCodeType = 3  // Day 05 P1
	opCodeTypeOutput         opCodeType = 4  // Day 05 P1
	opCodeTypeJumpIfTrue     opCodeType = 5  // Day 05 P2
	opCodeTypeJumpIfFalse    opCodeType = 6  // Day 05 P2
	opCodeTypeLessThan       opCodeType = 7  // Day 05 P2
	opCodeTypeEquals         opCodeType = 8  // Day 05 P2
	opCodeTypeExit           opCodeType = 99 // Day 02
)

type opCode struct {
	Type  opCodeType
	flags []opCodeFlag
}

func (o opCode) GetFlag(param int) opCodeFlag {
	if param-1 >= len(o.flags) {
		return opCodeFlagPosition
	}
	return o.flags[param-1]
}

func (o opCode) eq(in opCode) bool {
	return o.Type == in.Type && reflect.DeepEqual(o.flags, in.flags)
}

func parseOpCode(in int) opCode {
	out := opCode{}

	out.Type = opCodeType(in % 100)

	var paramFactor = 100
	for {
		if in < paramFactor {
			break
		}

		out.flags = append(out.flags, opCodeFlag((in % (paramFactor * 10) / paramFactor)))
		paramFactor *= 10
	}

	return out
}

func parseIntcode(code string) ([]int, error) {
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

func executeIntcode(code []int, in, out chan int) ([]int, error) {
	var pos int

	if out != nil {
		defer close(out)
	}

	getParamValue := func(param int, op opCode) int {
		switch op.GetFlag(param) {

		case opCodeFlagImmediate:
			return code[pos+param]

		case opCodeFlagPosition:
			return code[code[pos+param]]

		default:
			panic(errors.Errorf("Unexpected opCodeFlag %d", op.GetFlag(param)))

		}
	}

	for {
		if pos >= len(code) {
			return nil, errors.Errorf("Code position out of bounds: %d (len=%d)", pos, len(code))
		}

		// Position is expected to be an OpCode
		op := parseOpCode(code[pos])
		switch op.Type {

		case opCodeTypeAddition: // p1 + p2 => p3
			code[code[pos+3]] = getParamValue(1, op) + getParamValue(2, op)
			pos += 4

		case opCodeTypeMultiplication: // p1 * p2 => p3
			code[code[pos+3]] = getParamValue(1, op) * getParamValue(2, op)
			pos += 4

		case opCodeTypeInput: // in => p1
			code[code[pos+1]] = <-in
			pos += 2

		case opCodeTypeOutput: // p1 => out
			out <- getParamValue(1, op)
			pos += 2

		case opCodeTypeJumpIfTrue: // p1 != 0 => jmp
			if getParamValue(1, op) != 0 {
				pos = getParamValue(2, op)
				continue
			}
			pos += 3

		case opCodeTypeJumpIfFalse: // p1 == 0 => jmp
			if getParamValue(1, op) == 0 {
				pos = getParamValue(2, op)
				continue
			}
			pos += 3

		case opCodeTypeLessThan: // p1 < p2 => p3
			var res int
			if getParamValue(1, op) < getParamValue(2, op) {
				res = 1
			}
			code[code[pos+3]] = res
			pos += 4

		case opCodeTypeEquals: // p1 == p2 => p3
			var res int
			if getParamValue(1, op) == getParamValue(2, op) {
				res = 1
			}
			code[code[pos+3]] = res
			pos += 4

		case opCodeTypeExit: // exit
			return code, nil

		default:
			return nil, errors.Errorf("Encountered invalid operation %d (parsed %#v)", code[pos], op)

		}
	}
}
