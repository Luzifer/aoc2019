package aoc2019

import (
	"context"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var intcodeDebugging = false

type opCodeFlag int64

const (
	opCodeFlagPosition opCodeFlag = iota
	opCodeFlagImmediate
	opCodeFlagRelative
)

type opCodeType int64

const (
	opCodeTypeAddition       opCodeType = 1  // Day 02
	opCodeTypeMultiplication opCodeType = 2  // Day 02
	opCodeTypeInput          opCodeType = 3  // Day 05 P1
	opCodeTypeOutput         opCodeType = 4  // Day 05 P1
	opCodeTypeJumpIfTrue     opCodeType = 5  // Day 05 P2
	opCodeTypeJumpIfFalse    opCodeType = 6  // Day 05 P2
	opCodeTypeLessThan       opCodeType = 7  // Day 05 P2
	opCodeTypeEquals         opCodeType = 8  // Day 05 P2
	opCodeTypeAdjRelBase     opCodeType = 9  // Day 09
	opCodeTypeExit           opCodeType = 99 // Day 02
)

type opCode struct {
	Type  opCodeType
	flags []opCodeFlag
}

func (o opCode) GetFlag(param int64) opCodeFlag {
	if param-1 >= int64(len(o.flags)) {
		return opCodeFlagPosition
	}
	return o.flags[param-1]
}

func (o opCode) eq(in opCode) bool {
	return o.Type == in.Type && reflect.DeepEqual(o.flags, in.flags)
}

func parseOpCode(in int64) opCode {
	out := opCode{}

	out.Type = opCodeType(in % 100)

	var paramFactor int64 = 100
	for {
		if in < paramFactor {
			break
		}

		out.flags = append(out.flags, opCodeFlag((in % (paramFactor * 10) / paramFactor)))
		paramFactor *= 10
	}

	return out
}

func cloneIntcode(in []int64) []int64 {
	out := make([]int64, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}

func parseIntcode(code string) ([]int64, error) {
	parts := strings.Split(code, ",")

	var out []int64
	for _, n := range parts {
		v, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}

	return out, nil
}

type intcodeParams struct {
	// Intcode program to execute
	Code []int64
	// Context to execute the program in (program might hang on input if context is closed during input directive)
	Context context.Context
	// Channel / Callback to query on input directive
	In interface{}
	// Channel to use for output directive
	Out chan int64
}

func executeIntcode(code []int64, in interface{}, out chan int64) ([]int64, error) {
	return executeIntcodeWithParams(intcodeParams{
		Code:    code,
		Context: context.Background(),
		In:      in,
		Out:     out,
	})
}

func executeIntcodeWithParams(params intcodeParams) ([]int64, error) {
	var (
		code         = params.Code
		inCB         func() (int64, error)
		pos          int64
		relativeBase int64
	)

	if params.Out != nil {
		defer close(params.Out)
	}

	switch params.In.(type) {
	case nil:
		inCB = func() (int64, error) { return 0, errors.New("No input available") }
	case chan int64:
		inCB = func() (int64, error) { return <-(params.In.(chan int64)), nil }
	case func() (int64, error):
		inCB = params.In.(func() (int64, error))
	default:
		return nil, errors.New("Unsupported input type")
	}

	transformPos := func(param int64, op opCode, write bool) int64 {
		var addr int64

		switch op.GetFlag(param) {

		case opCodeFlagImmediate:
			if write {
				addr = code[pos+param]
			} else {
				addr = pos + param
			}

		case opCodeFlagPosition:
			addr = code[pos+param]

		case opCodeFlagRelative:
			addr = code[pos+param] + int64(relativeBase)

		default:
			panic(errors.Errorf("Unexpected opCodeFlag %d", op.GetFlag(param)))

		}

		return addr
	}

	getParamValue := func(param int64, op opCode) int64 {
		var addr = transformPos(param, op, false)

		if addr >= int64(len(code)) {
			return 0
		}

		return code[addr]
	}

	setParamValue := func(param, value int64, op opCode) {
		var addr = transformPos(param, op, false)

		if addr >= int64(len(code)) {
			// Write outside memory, increase memory
			var tmp = make([]int64, addr+1)
			for i, v := range code {
				tmp[i] = v
			}
			code = tmp
		}

		code[addr] = value
	}

	for {
		if pos >= int64(len(code)) {
			return nil, errors.Errorf("Code position out of bounds: %d (len=%d)", pos, len(code))
		}

		if err := params.Context.Err(); err != nil {
			return nil, errors.Wrap(err, "Context closed")
		}

		// Position is expected to be an OpCode
		op := parseOpCode(code[pos])

		if intcodeDebugging {
			log.Printf("OpCode execution: %#v", op)
		}

		switch op.Type {

		case opCodeTypeAddition: // p1 + p2 => p3
			setParamValue(3, getParamValue(1, op)+getParamValue(2, op), op)
			pos += 4

		case opCodeTypeMultiplication: // p1 * p2 => p3
			setParamValue(3, getParamValue(1, op)*getParamValue(2, op), op)
			pos += 4

		case opCodeTypeInput: // in => p1
			v, err := inCB()
			if err != nil {
				return nil, errors.Wrap(err, "Unable to read input")
			}
			setParamValue(1, v, op)
			pos += 2

		case opCodeTypeOutput: // p1 => out
			params.Out <- getParamValue(1, op)
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
			var res int64
			if getParamValue(1, op) < getParamValue(2, op) {
				res = 1
			}
			setParamValue(3, res, op)
			pos += 4

		case opCodeTypeEquals: // p1 == p2 => p3
			var res int64
			if getParamValue(1, op) == getParamValue(2, op) {
				res = 1
			}
			setParamValue(3, res, op)
			pos += 4

		case opCodeTypeAdjRelBase:
			relativeBase += getParamValue(1, op)
			pos += 2

		case opCodeTypeExit: // exit
			return code, nil

		default:
			return nil, errors.Errorf("Encountered invalid operation %d (parsed %#v)", code[pos], op)

		}
	}
}
