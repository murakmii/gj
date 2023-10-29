package vm

import (
	"bytes"
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
		syncObj   *Instance
	}

	StackTraceElement struct {
		class  string
		method string
		file   *string
		line   int32
	}
)

func NewFrame(curClass *Class, curMethod *class_file.MethodInfo) *Frame {
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

func (frame *Frame) SetLocal(index int, v interface{}) *Frame {
	frame.locals[index] = v
	return frame
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

func (frame *Frame) NextParamUint32() uint32 {
	return frame.code.ReadUint32()
}

func (frame *Frame) NextAlign(align int) {
	frame.code.SkipToAlign(align)
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
			if exTable.CatchType() == 0 {
				handler := exTable.HandlerPC()
				return &handler
			}

			catchType := frame.curClass.File().ConstantPool().ClassInfo(exTable.CatchType())
			if thrown.Class().IsSubClassOf(catchType) {
				handler := exTable.HandlerPC()
				return &handler
			}
		}
	}
	return nil
}

func (frame *Frame) Trace() *StackTraceElement {
	var file *string
	fileAttr := frame.curClass.File().SourceFile()
	if fileAttr != 0 {
		file = frame.curClass.File().ConstantPool().Utf8(uint16(fileAttr))
	}

	line := int32(-1)
	table := frame.curMethod.Code().LineNumberTable()
	if table != nil && table[frame.pc] > 0 {
		line = int32(table[frame.pc])
	}

	return NewStackTraceElement(
		frame.curClass.File().ThisClass(),
		*(frame.curMethod.Name()),
		file,
		line,
	)
}

func NewStackTraceElement(class, method string, file *string, line int32) *StackTraceElement {
	return &StackTraceElement{
		class:  class,
		method: method,
		file:   file,
		line:   line,
	}
}

func (trace *StackTraceElement) ToJava(vm *VM) *Instance {
	traceClass, _ := vm.Class("java/lang/StackTraceElement", nil)

	javaTrace := NewInstance(traceClass)

	javaTrace.PutField("declaringClass", "Ljava/lang/String;", NewString(vm, trace.class))
	javaTrace.PutField("methodName", "Ljava/lang/String;", NewString(vm, trace.method))
	javaTrace.PutField("lineNumber", "I", trace.line)

	if trace.file != nil {
		javaTrace.PutField("fileName", "Ljava/lang/String;", NewString(vm, *trace.file))
	}

	return javaTrace
}
