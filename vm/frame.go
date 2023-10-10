package vm

import (
	"bytes"
	"fmt"
	"github.com/murakmii/gj/class_file"
	"github.com/murakmii/gj/util"
)

type (
	Frame struct {
		thread    *Thread
		locals    []interface{}
		curClass  *Class
		curMethod *class_file.MethodInfo
		opStack   []interface{}
		pc        *util.BinReader
	}
)

func NewFrame(thread *Thread, curClass *Class, curMethod *class_file.MethodInfo) *Frame {
	code := curMethod.Code()
	pc, _ := util.NewBinReader(bytes.NewReader(code.Code()))

	return &Frame{
		thread:    thread,
		locals:    make([]interface{}, code.MaxLocals()),
		curClass:  curClass,
		curMethod: curMethod,
		opStack:   nil,
		pc:        pc,
	}
}

func (frame *Frame) Execute() (CurrentFrameOperation, error) {
	for frame.pc.Remain() > 0 {
		op := frame.pc.ReadByte()
		instr := InstructionSet[op]
		if instr == nil {
			return NoFrameOp, fmt.Errorf("OP code(%#X) is NOT implemented", op)
		}

		frameOp, err := instr(frame)
		if err != nil {
			return NoFrameOp, err
		}

		if frameOp != NoFrameOp {
			return frameOp, nil
		}
	}

	return NoFrameOp, fmt.Errorf("end of code")
}

func (frame *Frame) PushOperand(value interface{}) {
	frame.opStack = append(frame.opStack, value)
}

func (frame *Frame) PopOperand() interface{} {
	last := len(frame.opStack) - 1
	if last == -1 {
		return nil
	}

	pop := frame.opStack[last]
	frame.opStack = frame.opStack[:last]
	return pop
}

func (frame *Frame) ClearOperand() {
	frame.opStack = nil
}
