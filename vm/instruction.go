package vm

import (
	"fmt"
	"github.com/murakmii/gj/class_file"
)

type (
	Instruction func(*Thread, *Frame) error
)

var InstructionSet [255]Instruction

func init() {
	InstructionSet[0x00] = func(_ *Thread, _ *Frame) error { return nil } // nop
	InstructionSet[0x01] = instrAConstNull

	InstructionSet[0x02] = instrIconst(-1)
	InstructionSet[0x03] = instrIconst(0)
	InstructionSet[0x04] = instrIconst(1)
	InstructionSet[0x05] = instrIconst(2)
	InstructionSet[0x06] = instrIconst(3)
	InstructionSet[0x07] = instrIconst(4)
	InstructionSet[0x08] = instrIconst(5)

	InstructionSet[0x15] = instrLoad
	InstructionSet[0x17] = instrLoad
	InstructionSet[0x19] = instrLoad

	InstructionSet[0x1A] = instrLoadN(0)
	InstructionSet[0x1B] = instrLoadN(1)
	InstructionSet[0x1C] = instrLoadN(2)
	InstructionSet[0x1D] = instrLoadN(3)

	InstructionSet[0x22] = instrLoadN(0)
	InstructionSet[0x23] = instrLoadN(1)
	InstructionSet[0x24] = instrLoadN(2)
	InstructionSet[0x25] = instrLoadN(3)

	InstructionSet[0x2A] = instrLoadN(0)
	InstructionSet[0x2B] = instrLoadN(1)
	InstructionSet[0x2C] = instrLoadN(2)
	InstructionSet[0x2D] = instrLoadN(3)

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
}

func ExecInstr(thread *Thread, frame *Frame, op byte) error {
	fmt.Printf("exec instr %#x\n", op)
	if InstructionSet[op] == nil {
		return fmt.Errorf("op(code = %#x) has been NOT implemented", op)
	}
	return InstructionSet[op](thread, frame)
}

func instrAConstNull(thread *Thread, frame *Frame) error {
	frame.PushOperand(nil)
	return nil
}

func instrIconst(n int) Instruction {
	return func(_ *Thread, frame *Frame) error {
		frame.PushOperand(n)
		return nil
	}
}

func instrLoad(_ *Thread, frame *Frame) error {
	frame.PushOperand(frame.Locals()[frame.NextParamByte()])
	return nil
}

func instrLoadN(n int) Instruction {
	return func(_ *Thread, frame *Frame) error {
		frame.PushOperand(frame.Locals()[n])
		return nil
	}
}

func instrDup(n, x int) Instruction {
	return func(_ *Thread, frame *Frame) error {
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

		return nil
	}
}

func instrReturn(thread *Thread, frame *Frame) error {
	thread.PopFrame()
	if thread.CurrentFrame() != nil {
		thread.CurrentFrame().PushOperand(frame.PopOperand())
	}
	return nil
}

func instrReturnVoid(thread *Thread, _ *Frame) error {
	thread.PopFrame()
	return nil
}

func instrPutStatic(thread *Thread, frame *Frame) error {
	className, name, desc := frame.CurrentClass().File().ConstantPool().Reference(frame.NextParamUint16())
	class, state, err := thread.VM().FindInitializedClass(className, thread)
	if err != nil {
		return err
	}
	if state == FailedInitialization {
		return fmt.Errorf("failed initialization of waiting class: %s", *className)
	}

	resolvedClass, resolvedField := class.ResolveField(*name, *desc)
	resolvedClass.SetStaticField(resolvedField.Name(), frame.PopOperand())

	return nil
}

func instrInvokeVirtual(thread *Thread, frame *Frame) error {
	_, name, desc := frame.curClass.File().ConstantPool().Reference(frame.NextParamUint16())
	instance := frame.PeekFromTop(class_file.ParseDescriptor(desc)).(*Instance)
	if instance == nil {
		return fmt.Errorf("receiver instance is null")
	}

	resolvedClass, resolvedMethod := instance.Class().ResolveMethod(*name, *desc)
	if resolvedClass == nil || !resolvedMethod.IsCallableForInstance() {
		return fmt.Errorf("method not found: %s.%s", *name, *desc)
	}

	return thread.ExecMethod(resolvedClass, resolvedMethod)
}

func instrInvokeSpecial(thread *Thread, frame *Frame) error {
	className, name, desc := frame.curClass.File().ConstantPool().Reference(frame.NextParamUint16())
	class, state, err := thread.VM().FindInitializedClass(className, thread)
	if err != nil {
		return err
	}
	if state == FailedInitialization {
		return fmt.Errorf("failed initialization of waiting class: %s", *className)
	}

	resolvedClass, resolvedMethod := class.ResolveMethod(*name, *desc)
	if resolvedClass == nil || !resolvedMethod.IsCallableForInstance() {
		return fmt.Errorf("method not found: %s.%s", *name, *desc)
	}

	return thread.ExecMethod(resolvedClass, resolvedMethod)
}

func instrInvokeStatic(thread *Thread, frame *Frame) error {
	className, name, desc := frame.curClass.File().ConstantPool().Reference(frame.NextParamUint16())
	class, state, err := thread.VM().FindInitializedClass(className, thread)
	if err != nil {
		return err
	}
	if state == FailedInitialization {
		return fmt.Errorf("failed initialization of waiting class: %s", *className)
	}

	resolvedClass, resolvedMethod := class.ResolveMethod(*name, *desc)
	if resolvedClass == nil || !resolvedMethod.IsCallableAsStatic() {
		return fmt.Errorf("method not found: %s.%s", *name, *desc)
	}

	fmt.Printf("invoke static %s.%s:%s\n", resolvedClass.File().ThisClass(), *resolvedMethod.Name(), *resolvedMethod.Descriptor())

	return thread.ExecMethod(resolvedClass, resolvedMethod)
}

func instrNew(thread *Thread, frame *Frame) error {
	className := frame.curClass.File().ConstantPool().ClassInfo(frame.NextParamUint16())
	class, state, err := thread.VM().FindInitializedClass(className, thread)
	if err != nil {
		return err
	}
	if state == FailedInitialization {
		return fmt.Errorf("failed initialization of waiting class: %s", *className)
	}

	frame.PushOperand(NewInstance(class))
	return nil
}

func instrANewArray(_ *Thread, frame *Frame) error {
	className := frame.curClass.File().ConstantPool().ClassInfo(frame.NextParamUint16())
	frame.PushOperand(NewArray(*className, frame.PopOperand().(int)))
	return nil
}
