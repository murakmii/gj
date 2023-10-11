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
		pc := frame.pc.Pos()
		frameOp, err := ExecInstr(frame, frame.pc.ReadByte())
		if err != nil {
			return NoFrameOp, err
		}

		switch frameOp {
		case ThrowFromFrame:
			thrown := frame.PopOperand().(*Instance)
			frame.ClearOperand()
			frame.PushOperand(thrown)

			handlerPC := frame.findExceptionHandler(uint16(pc), thrown)
			if handlerPC == -1 {
				return ThrowFromFrame, nil
			}
			frame.pc.Seek(handlerPC)

		case ReturnFromFrame:
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

func (frame *Frame) findExceptionHandler(curPC uint16, thrown *Instance) int {
	for _, exceptionTable := range frame.curMethod.Code().ExceptionTable() {
		if exceptionTable.HandlerStart() <= curPC && curPC < exceptionTable.HandlerEnd() {
			catchType := frame.curClass.File().ConstantPool().ClassInfo(exceptionTable.CatchType())
			if thrown.IsSubClassOf(catchType) {
				return int(exceptionTable.HandlerPC())
			}
		}
	}
	return -1
}
