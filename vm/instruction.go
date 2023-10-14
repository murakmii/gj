package vm

import (
	"fmt"
	"github.com/murakmii/gj/class_file"
)

type (
	Instruction           func(frame *Frame) (CurrentFrameOperation, error)
	CurrentFrameOperation uint8
)

const (
	NoFrameOp CurrentFrameOperation = iota
	ReturnFromFrame
	ThrowFromFrame
)

var InstructionSet [255]Instruction

func init() {
	InstructionSet[0x01] = instrAConstNull

	InstructionSet[0x02] = instrIconst(-1)
	InstructionSet[0x03] = instrIconst(0)
	InstructionSet[0x04] = instrIconst(1)
	InstructionSet[0x05] = instrIconst(2)
	InstructionSet[0x06] = instrIconst(3)
	InstructionSet[0x07] = instrIconst(4)
	InstructionSet[0x08] = instrIconst(5)

	InstructionSet[0x59] = instrDup(1, 0)
	InstructionSet[0x5A] = instrDup(1, 1)
	InstructionSet[0x5B] = instrDup(1, 2)
	InstructionSet[0x5C] = instrDup(2, 0)
	InstructionSet[0x5D] = instrDup(2, 1)
	InstructionSet[0x5E] = instrDup(2, 2)

	InstructionSet[0xAC] = instrReturn
	InstructionSet[0xAD] = instrReturn
	InstructionSet[0xAE] = instrReturn
	InstructionSet[0xAF] = instrReturn
	InstructionSet[0xB0] = instrReturn
	InstructionSet[0xB1] = instrReturnVoid

	InstructionSet[0xB3] = instrPutStatic

	InstructionSet[0xB7] = instrInvokeVirtual
	InstructionSet[0xB7] = instrInvokeSpecial
	InstructionSet[0xB8] = instrInvokeStatic

	InstructionSet[0xBB] = instrNew

	InstructionSet[0xBD] = instrANewArray

	InstructionSet[0xB8] = func(frame *Frame) (CurrentFrameOperation, error) { // invokestatic
		return NoFrameOp, nil
	}
}

func ExecInstr(frame *Frame, op byte) (CurrentFrameOperation, error) {
	instr := InstructionSet[op]
	if instr == nil {
		return NoFrameOp, fmt.Errorf("OP code(%#X) is NOT implemented", op)
	}
	return instr(frame)
}

func instrAConstNull(frame *Frame) (CurrentFrameOperation, error) {
	frame.PushOperand(nil)
	return NoFrameOp, nil
}

func instrIconst(n int) Instruction {
	return func(frame *Frame) (CurrentFrameOperation, error) {
		frame.PushOperand(n)
		return NoFrameOp, nil
	}
}

func instrDup(n, x int) Instruction {
	return func(frame *Frame) (CurrentFrameOperation, error) {
		top := make([]interface{}, n+x)
		for i := len(top) - 1; i >= 0; i-- {
			top[i] = frame.PopOperand()
		}

		for i := 0; i < n; i++ {
			frame.PushOperand(top[x+i])
		}

		for _, v := range top {
			frame.PushOperand(v)
		}

		return NoFrameOp, nil
	}
}

func instrReturn(frame *Frame) (CurrentFrameOperation, error) {
	ret := frame.PopOperand()
	frame.ClearOperand()
	frame.PushOperand(ret)
	return ReturnFromFrame, nil
}

func instrReturnVoid(frame *Frame) (CurrentFrameOperation, error) {
	frame.ClearOperand()
	return ReturnFromFrame, nil
}

func instrPutStatic(frame *Frame) (CurrentFrameOperation, error) {
	className, name, desc := frame.CurrentClass().File().ConstantPool().Reference(frame.PC().ReadUint16())
	class, state, err := frame.Thread().VM().FindInitializedClass(className, frame.Thread())
	if err != nil {
		return NoFrameOp, err
	}
	if state == FailedInitialization {
		return NoFrameOp, fmt.Errorf("failed initialization of waiting class: %s", *className)
	}

	resolvedClass, resolvedField := class.ResolveField(*name, *desc)
	resolvedClass.SetStaticField(resolvedField.Name(), frame.PopOperand())

	return NoFrameOp, nil
}

func instrInvokeVirtual(frame *Frame) (CurrentFrameOperation, error) {
	_, name, desc := frame.curClass.File().ConstantPool().Reference(frame.PC().ReadUint16())
	instance := frame.PeekFromTop(class_file.ParseDescriptor(desc)).(*Instance)
	if instance == nil {
		return NoFrameOp, fmt.Errorf("receiver instance is null")
	}

	resolvedClass, resolvedMethod := instance.Class().ResolveMethod(*name, *desc)
	if resolvedClass == nil || !resolvedMethod.IsCallableForInstance() {
		return NoFrameOp, fmt.Errorf("method not found: %s.%s", *name, *desc)
	}

	ret, frameOp, err := ExecuteFrame(frame, resolvedClass, resolvedMethod, true)
	if err != nil {
		return NoFrameOp, err
	}

	frame.PushOperand(ret)
	return frameOp, nil
}

func instrInvokeSpecial(frame *Frame) (CurrentFrameOperation, error) {
	className, name, desc := frame.curClass.File().ConstantPool().Reference(frame.PC().ReadUint16())
	class, state, err := frame.Thread().VM().FindInitializedClass(className, frame.Thread())
	if err != nil {
		return NoFrameOp, err
	}
	if state == FailedInitialization {
		return NoFrameOp, fmt.Errorf("failed initialization of waiting class: %s", *className)
	}

	resolvedClass, resolvedMethod := class.ResolveMethod(*name, *desc)
	if resolvedClass == nil || !resolvedMethod.IsCallableForInstance() {
		return NoFrameOp, fmt.Errorf("method not found: %s.%s", *name, *desc)
	}

	ret, frameOp, err := ExecuteFrame(frame, resolvedClass, resolvedMethod, true)
	if err != nil {
		return NoFrameOp, err
	}

	frame.PushOperand(ret)
	return frameOp, nil
}

func instrInvokeStatic(frame *Frame) (CurrentFrameOperation, error) {
	className, name, desc := frame.curClass.File().ConstantPool().Reference(frame.PC().ReadUint16())
	class, state, err := frame.Thread().VM().FindInitializedClass(className, frame.Thread())
	if err != nil {
		return NoFrameOp, err
	}
	if state == FailedInitialization {
		return NoFrameOp, fmt.Errorf("failed initialization of waiting class: %s", *className)
	}

	resolvedClass, resolvedMethod := class.ResolveMethod(*name, *desc)
	if resolvedClass == nil || !resolvedMethod.IsCallableAsStatic() {
		return NoFrameOp, fmt.Errorf("method not found: %s.%s", *name, *desc)
	}

	ret, frameOp, err := ExecuteFrame(frame, resolvedClass, resolvedMethod, false)
	if err != nil {
		return NoFrameOp, err
	}

	frame.PushOperand(ret)
	return frameOp, nil
}

func instrNew(frame *Frame) (CurrentFrameOperation, error) {
	className := frame.curClass.File().ConstantPool().ClassInfo(frame.PC().ReadUint16())
	class, state, err := frame.Thread().VM().FindInitializedClass(className, frame.Thread())
	if err != nil {
		return NoFrameOp, err
	}
	if state == FailedInitialization {
		return NoFrameOp, fmt.Errorf("failed initialization of waiting class: %s", *className)
	}

	frame.PushOperand(NewInstance(class))
	return NoFrameOp, nil
}

func instrANewArray(frame *Frame) (CurrentFrameOperation, error) {
	className := frame.curClass.File().ConstantPool().ClassInfo(frame.PC().ReadUint16())
	frame.PushOperand(NewArray(*className, frame.PopOperand().(int)))

	return NoFrameOp, nil
}
