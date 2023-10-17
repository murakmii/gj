package vm

import (
	"fmt"
	"github.com/murakmii/gj/class_file"
)

type (
	Instruction func(*Thread, *Frame) error
)

var (
	InstructionSet [255]Instruction

	typeCodes = []string{"Z", "C", "F", "D", "B", "S", "I", "J"}
)

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

	InstructionSet[0x10] = InstrBiPush
	InstructionSet[0x11] = InstrSiPush

	InstructionSet[0x12] = instrLdc(func(f *Frame) uint16 { return uint16(f.NextParamByte()) })
	InstructionSet[0x13] = instrLdc(func(f *Frame) uint16 { return f.NextParamUint16() })
	InstructionSet[0x14] = InstructionSet[0x13]

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

	InstructionSet[0x36] = instrStore
	InstructionSet[0x37] = instrStore
	InstructionSet[0x38] = instrStore
	InstructionSet[0x39] = instrStore
	InstructionSet[0x3A] = instrStore
	InstructionSet[0x3B] = instrStoreN(0)
	InstructionSet[0x3C] = instrStoreN(1)
	InstructionSet[0x3D] = instrStoreN(2)
	InstructionSet[0x3E] = instrStoreN(3)

	InstructionSet[0x4B] = instrStoreN(0)
	InstructionSet[0x4C] = instrStoreN(1)
	InstructionSet[0x4D] = instrStoreN(2)
	InstructionSet[0x4E] = instrStoreN(3)

	InstructionSet[0x59] = instrDup(1, 0)
	InstructionSet[0x5A] = instrDup(1, 1)
	InstructionSet[0x5B] = instrDup(1, 2)
	InstructionSet[0x5C] = instrDup(2, 0)
	InstructionSet[0x5D] = instrDup(2, 1)
	InstructionSet[0x5E] = instrDup(2, 2)

	InstructionSet[0x60] = instrAdd[int]()
	InstructionSet[0x61] = instrAdd[int64]()

	InstructionSet[0x78] = instrShiftLeft[int](0x1F)
	InstructionSet[0x79] = instrShiftLeft[int64](0x3F)
	InstructionSet[0x7A] = instrShiftRight[int](0x1F)
	InstructionSet[0x7B] = instrShiftRight[int64](0x3F)

	InstructionSet[0x7E] = instrAnd[int]
	InstructionSet[0x7F] = instrAnd[int64]

	InstructionSet[0x85] = InstrI2L

	InstructionSet[0x99] = instrIf(func(i int) bool { return i == 0 })
	InstructionSet[0x9A] = instrIf(func(i int) bool { return i != 0 })
	InstructionSet[0x9B] = instrIf(func(i int) bool { return i < 0 })
	InstructionSet[0x9C] = instrIf(func(i int) bool { return i <= 0 })
	InstructionSet[0x9D] = instrIf(func(i int) bool { return i > 0 })
	InstructionSet[0x9E] = instrIf(func(i int) bool { return i >= 0 })

	InstructionSet[0x9F] = instrIfICmp(func(v1 int, v2 int) bool { return v1 == v2 })
	InstructionSet[0xA0] = instrIfICmp(func(v1 int, v2 int) bool { return v1 != v2 })
	InstructionSet[0xA1] = instrIfICmp(func(v1 int, v2 int) bool { return v1 < v2 })
	InstructionSet[0xA2] = instrIfICmp(func(v1 int, v2 int) bool { return v1 <= v2 })
	InstructionSet[0xA3] = instrIfICmp(func(v1 int, v2 int) bool { return v1 > v2 })
	InstructionSet[0xA4] = instrIfICmp(func(v1 int, v2 int) bool { return v1 >= v2 })

	InstructionSet[0xA7] = instrGoTo

	InstructionSet[0xAC] = instrReturn
	InstructionSet[0xAD] = instrReturn
	InstructionSet[0xAE] = instrReturn
	InstructionSet[0xAF] = instrReturn
	InstructionSet[0xB0] = instrReturn
	InstructionSet[0xB1] = instrReturnVoid

	InstructionSet[0xB2] = instrGetStatic
	InstructionSet[0xB3] = instrPutStatic
	InstructionSet[0xB4] = instrGetField
	InstructionSet[0xB5] = instrPutField

	InstructionSet[0xB6] = instrInvokeVirtual
	InstructionSet[0xB7] = instrInvokeSpecial
	InstructionSet[0xB8] = instrInvokeStatic

	InstructionSet[0xBB] = instrNew
	InstructionSet[0xBC] = instrNewArray
	InstructionSet[0xBD] = instrANewArray
	InstructionSet[0xBE] = instrArrayLength

	InstructionSet[0xC6] = instrIfNull
	InstructionSet[0xC7] = instrIfNonNull
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

func InstrBiPush(_ *Thread, frame *Frame) error {
	frame.PushOperand(int(frame.NextParamByte()))
	return nil
}

func InstrSiPush(_ *Thread, frame *Frame) error {
	frame.PushOperand(int(frame.NextParamUint16()))
	return nil
}

func instrLdc(idxLoader func(*Frame) uint16) Instruction {
	return func(thread *Thread, frame *Frame) error {
		switch entry := frame.CurrentClass().File().ConstantPool().Entry(idxLoader(frame)).(type) {
		case int, float32, int64, float64:
			frame.PushOperand(entry)

		case class_file.StringCpInfo:
			js, err := thread.VM().JavaString(thread, frame.CurrentClass().File().ConstantPool().Utf8(uint16(entry)))
			if err != nil {
				return err
			}
			frame.PushOperand(js)

		case class_file.ClassCpInfo:
			className := "java/lang/Class"
			class, state, err := thread.VM().FindInitializedClass(&className, thread)
			if err != nil {
				return err
			}
			if state == FailedInitialization {
				return fmt.Errorf("failed initialization of class class in LDC")
			}
			frame.PushOperand(NewInstance(class))

		default:
			return fmt.Errorf("LDC unsupport %T:%+v", entry, entry)
		}

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

func instrStore(_ *Thread, frame *Frame) error {
	frame.SetLocal(int(frame.NextParamByte()), frame.PopOperand())
	return nil
}

func instrStoreN(n int) Instruction {
	return func(thread *Thread, frame *Frame) error {
		frame.SetLocal(n, frame.PopOperand())
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

func instrAdd[T int | int64]() Instruction {
	return func(thread *Thread, frame *Frame) error {
		v2, ok := frame.PopOperand().(T)
		if !ok {
			return fmt.Errorf("popped value2 for (i|l)add is invalid type")
		}
		v1, ok := frame.PopOperand().(T)
		if !ok {
			return fmt.Errorf("popped value1 for (i|l)add is invalid type")
		}

		frame.PushOperand(v1 + v2)
		return nil
	}
}

func instrShiftLeft[T int | int64](mask int) Instruction {
	return func(_ *Thread, frame *Frame) error {
		v2, ok := frame.PopOperand().(int)
		if !ok {
			return fmt.Errorf("popped value2 for (i|l)shl is invalid type")
		}
		v1, ok := frame.PopOperand().(T)
		if !ok {
			return fmt.Errorf("popped value1 for (i|l)shl is invalid type")
		}

		frame.PushOperand(v1 << (v2 & mask))
		return nil
	}
}

func instrShiftRight[T int | int64](mask int) Instruction {
	return func(_ *Thread, frame *Frame) error {
		v2, ok := frame.PopOperand().(int)
		if !ok {
			return fmt.Errorf("popped value2 for (i|l)shr is invalid type")
		}
		v1, ok := frame.PopOperand().(T)
		if !ok {
			return fmt.Errorf("popped value1 for (i|l)shr is invalid type")
		}

		frame.PushOperand(v1 >> (v2 & mask))
		return nil
	}
}

func instrAnd[T int | int64](_ *Thread, frame *Frame) error {
	v2, ok := frame.PopOperand().(T)
	if !ok {
		return fmt.Errorf("popped value2 for (i|l)and is invalid type")
	}
	v1, ok := frame.PopOperand().(T)
	if !ok {
		return fmt.Errorf("popped value1 for (i|l)and is invalid type")
	}

	frame.PushOperand(v1 & v2)
	return nil
}

func InstrI2L(_ *Thread, frame *Frame) error {
	i, ok := frame.PopOperand().(int)
	if !ok {
		return fmt.Errorf("popped value for i2l is NOT int")
	}

	frame.PushOperand(int64(i))
	return nil
}

func instrIf(matcher func(int) bool) Instruction {
	return func(_ *Thread, frame *Frame) error {
		branch := int16(frame.NextParamUint16())
		value, ok := frame.PopOperand().(int)
		if !ok {
			return fmt.Errorf("popped value for if<cond> is NOT int")
		}

		if matcher(value) {
			frame.JumpPC(uint16(int16(frame.PC()) + branch))
		}
		return nil
	}
}

func instrIfICmp(comparator func(int, int) bool) Instruction {
	return func(_ *Thread, frame *Frame) error {
		branch := int16(frame.NextParamUint16())
		v2, ok := frame.PopOperand().(int)
		if !ok {
			return fmt.Errorf("popped value2 for if_icmp<cond> is NOT int")
		}
		v1, ok := frame.PopOperand().(int)
		if !ok {
			return fmt.Errorf("popped value1 for if_icmp<cond> is NOT int")
		}

		if comparator(v1, v2) {
			frame.JumpPC(uint16(int16(frame.PC()) + branch))
		}
		return nil
	}
}

func instrGoTo(_ *Thread, frame *Frame) error {
	branch := int16(frame.NextParamUint16())
	frame.JumpPC(uint16(int16(frame.PC()) + branch))
	return nil
}

func instrReturn(thread *Thread, frame *Frame) error {
	thread.PopFrame()
	if thread.CurrentFrame() != nil {
		thread.CurrentFrame().PushOperand(frame.PopOperand())
	} else {
		thread.SetResult(frame.PopOperand())
	}
	return nil
}

func instrReturnVoid(thread *Thread, _ *Frame) error {
	thread.PopFrame()
	return nil
}

func instrGetStatic(thread *Thread, frame *Frame) error {
	className, name, desc := frame.CurrentClass().File().ConstantPool().Reference(frame.NextParamUint16())
	class, state, err := thread.VM().FindInitializedClass(className, thread)
	if err != nil {
		return err
	}
	if state == FailedInitialization {
		return fmt.Errorf("failed initialization of waiting class: %s", *className)
	}

	resolvedClass, resolvedField := class.ResolveField(*name, *desc)
	frame.PushOperand(resolvedClass.GetStaticField(resolvedField.Name()))

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

func instrGetField(_ *Thread, frame *Frame) error {
	class, name, _ := frame.curClass.File().ConstantPool().Reference(frame.NextParamUint16())
	instance := frame.PopOperand().(*Instance)
	if instance == nil {
		return fmt.Errorf("objectref for getfield is null")
	}

	frame.PushOperand(instance.GetField(class, name))
	return nil
}

func instrPutField(_ *Thread, frame *Frame) error {
	class, name, _ := frame.curClass.File().ConstantPool().Reference(frame.NextParamUint16())
	value := frame.PopOperand()
	instance := frame.PopOperand().(*Instance)
	if instance == nil {
		return fmt.Errorf("objectref for getfield is null")
	}

	instance.PutField(class, name, value)
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

func instrNewArray(_ *Thread, frame *Frame) error {
	frame.PushOperand(NewArray(typeCodes[frame.NextParamByte()-4], frame.PopOperand().(int)))
	return nil
}

func instrANewArray(_ *Thread, frame *Frame) error {
	className := frame.curClass.File().ConstantPool().ClassInfo(frame.NextParamUint16())
	frame.PushOperand(NewArray(*className, frame.PopOperand().(int)))
	return nil
}

func instrArrayLength(_ *Thread, frame *Frame) error {
	array, ok := frame.PopOperand().(*Array)
	if !ok {
		return fmt.Errorf("called arraylength for instance is NOT array")
	}
	frame.PushOperand(array.Length())
	return nil
}

func instrIfNonNull(_ *Thread, frame *Frame) error {
	branch := int16(frame.NextParamUint16())
	if frame.PopOperand() == nil {
		return nil
	}

	frame.JumpPC(uint16(int16(frame.PC()) + branch))
	return nil
}

func instrIfNull(_ *Thread, frame *Frame) error {
	branch := int16(frame.NextParamUint16())
	if frame.PopOperand() != nil {
		return nil
	}

	frame.JumpPC(uint16(int16(frame.PC()) + branch))
	return nil
}
