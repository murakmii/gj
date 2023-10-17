package vm

import (
	"bytes"
	"fmt"
	"github.com/murakmii/gj/class_file"
	"github.com/murakmii/gj/util"
)

type (
	Frame struct {
		locals    []interface{}
		curClass  *Class
		curMethod *class_file.MethodInfo
		opStack   []interface{}
		code      *util.BinReader
		pc        uint16
	}

	FrameResult struct {
		OperandStack   []interface{}
		ThrowException bool
	}
)

func NewFrame(curClass *Class, curMethod *class_file.MethodInfo) *Frame {
	fmt.Printf("enter new frame: %s.%s:%s\n", curClass.File().ThisClass(), *curMethod.Name(), *curMethod.Descriptor())

	code := curMethod.Code()
	codeReader, _ := util.NewBinReader(bytes.NewReader(code.Code()))

	return &Frame{
		locals:    make([]interface{}, code.MaxLocals()),
		curClass:  curClass,
		curMethod: curMethod,
		opStack:   nil,
		code:      codeReader,
		pc:        0,
	}
}

func (frame *Frame) SetLocal(index int, v interface{}) {
	frame.locals[index] = v
}

func (frame *Frame) SetLocals(vars []interface{}) *Frame {
	i := 0
	for _, v := range vars {
		frame.locals[i] = v

		switch v.(type) {
		case int64, float64:
			i += 2
		default:
			i++
		}
	}

	return frame
}

func (frame *Frame) CurrentClass() *Class {
	return frame.curClass
}

func (frame *Frame) CurrentMethod() *class_file.MethodInfo {
	return frame.curMethod
}

func (frame *Frame) NextInstr() byte {
	frame.pc = uint16(frame.code.Pos())
	return frame.code.ReadByte()
}

func (frame *Frame) NextParamByte() byte {
	return frame.code.ReadByte()
}

func (frame *Frame) NextParamUint16() uint16 {
	return frame.code.ReadUint16()
}

func (frame *Frame) PC() uint16 {
	return frame.pc
}

func (frame *Frame) JumpPC(pc uint16) {
	frame.pc = pc
	frame.code.Seek(int(pc))
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

func (frame *Frame) PopOperands(n int) []interface{} {
	popped := make([]interface{}, n)
	for i := n - 1; i >= 0; i-- {
		popped[i] = frame.PopOperand()
	}
	return popped
}

func (frame *Frame) PeekFromTop(index int) interface{} {
	i := len(frame.opStack) - 1 - index
	if i < 0 {
		return nil
	}
	return frame.opStack[i]
}

func (frame *Frame) ClearOperand() {
	frame.opStack = nil
}

func (frame *Frame) Locals() []interface{} {
	return frame.locals
}

func (frame *Frame) FindCurrentExceptionHandler(thrown *Instance) *uint16 {
	for _, exTable := range frame.curMethod.Code().ExceptionTable() {
		if exTable.HandlerStart() <= frame.pc && frame.pc < exTable.HandlerEnd() {
			catchType := frame.curClass.File().ConstantPool().ClassInfo(exTable.CatchType())

			if thrown.Class().IsSubClassOf(catchType) {
				handler := exTable.HandlerPC()
				return &handler
			}
		}
	}
	return nil
}
