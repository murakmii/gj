package vm

import (
	"fmt"
	"github.com/murakmii/gj/class_file"
	"math"
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

	InstructionSet[0x09] = instrLConst(0)
	InstructionSet[0x0A] = instrLConst(1)

	InstructionSet[0x0B] = instrFConst(0.0)
	InstructionSet[0x0C] = instrFConst(1.0)
	InstructionSet[0x0D] = instrFConst(2.0)

	InstructionSet[0x10] = InstrBiPush
	InstructionSet[0x11] = InstrSiPush

	InstructionSet[0x12] = instrLdc(func(f *Frame) uint16 { return uint16(f.NextParamByte()) })
	InstructionSet[0x13] = instrLdc(func(f *Frame) uint16 { return f.NextParamUint16() })
	InstructionSet[0x14] = InstructionSet[0x13]

	InstructionSet[0x15] = instrLoad
	InstructionSet[0x16] = instrLoad
	InstructionSet[0x17] = instrLoad
	InstructionSet[0x18] = instrLoad
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

	InstructionSet[0x2E] = instrALoad
	InstructionSet[0x2F] = instrALoad
	InstructionSet[0x30] = instrALoad
	InstructionSet[0x31] = instrALoad
	InstructionSet[0x32] = instrALoad
	InstructionSet[0x33] = instrALoad
	InstructionSet[0x34] = instrALoad
	InstructionSet[0x35] = instrALoad

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

	InstructionSet[0x4F] = InstrAStore
	InstructionSet[0x50] = InstrAStore
	InstructionSet[0x51] = InstrAStore
	InstructionSet[0x52] = InstrAStore
	InstructionSet[0x53] = InstrAStore
	InstructionSet[0x54] = InstrAStore
	InstructionSet[0x55] = InstrAStore
	InstructionSet[0x56] = InstrAStore

	InstructionSet[0x57] = InstrPop(1)
	InstructionSet[0x58] = InstrPop(2)

	InstructionSet[0x59] = instrDup(1, 0)
	InstructionSet[0x5A] = instrDup(1, 1)
	InstructionSet[0x5B] = instrDup(1, 2)
	InstructionSet[0x5C] = instrDup(2, 0)
	InstructionSet[0x5D] = instrDup(2, 1)
	InstructionSet[0x5E] = instrDup(2, 2)

	InstructionSet[0x60] = instrAdd[int]()
	InstructionSet[0x61] = instrAdd[int64]()

	InstructionSet[0x64] = instrBiOp[int]("isub", func(v1 int, v2 int) int { return v1 - v2 })
	InstructionSet[0x65] = instrBiOp[int64]("lsub", func(v1 int64, v2 int64) int64 { return v1 - v2 })
	InstructionSet[0x68] = instrBiOp[int]("imul", func(v1 int, v2 int) int { return v1 * v2 })
	InstructionSet[0x70] = instrBiOp[int]("irem", func(v1 int, v2 int) int { return v1 % v2 })

	InstructionSet[0x6A] = instrBiOp[float32]("fmul", func(v1 float32, v2 float32) float32 { return v1 * v2 })

	InstructionSet[0x78] = instrShiftLeft[int](0x1F)
	InstructionSet[0x79] = instrShiftLeft[int64](0x3F)
	InstructionSet[0x7A] = instrShiftRight[int](0x1F) // TODO: Arithmetic
	InstructionSet[0x7B] = instrShiftRight[int64](0x3F)
	InstructionSet[0x7C] = instrShiftRight[int](0x1F)

	InstructionSet[0x7E] = instrAnd[int]
	InstructionSet[0x7F] = instrAnd[int64]
	InstructionSet[0x80] = instrBiOp[int]("ior", func(v1 int, v2 int) int { return v1 | v2 })
	InstructionSet[0x82] = instrBiOp[int]("ixor", func(v1 int, v2 int) int { return v1 ^ v2 })

	InstructionSet[0x84] = instrIInc

	InstructionSet[0x85] = InstrI2L
	InstructionSet[0x86] = InstrI2F

	InstructionSet[0x8B] = InstrF2I

	InstructionSet[0x94] = instrLCmp
	InstructionSet[0x95] = instrFCmp(-1)
	InstructionSet[0x96] = instrFCmp(1)

	InstructionSet[0x99] = instrIf(func(i int) bool { return i == 0 })
	InstructionSet[0x9A] = instrIf(func(i int) bool { return i != 0 })
	InstructionSet[0x9B] = instrIf(func(i int) bool { return i < 0 })
	InstructionSet[0x9C] = instrIf(func(i int) bool { return i >= 0 })
	InstructionSet[0x9D] = instrIf(func(i int) bool { return i > 0 })
	InstructionSet[0x9E] = instrIf(func(i int) bool { return i <= 0 })

	InstructionSet[0x9F] = instrIfICmp(func(v1 int, v2 int) bool { return v1 == v2 })
	InstructionSet[0xA0] = instrIfICmp(func(v1 int, v2 int) bool { return v1 != v2 })
	InstructionSet[0xA1] = instrIfICmp(func(v1 int, v2 int) bool { return v1 < v2 })
	InstructionSet[0xA2] = instrIfICmp(func(v1 int, v2 int) bool { return v1 >= v2 })
	InstructionSet[0xA3] = instrIfICmp(func(v1 int, v2 int) bool { return v1 > v2 })
	InstructionSet[0xA4] = instrIfICmp(func(v1 int, v2 int) bool { return v1 <= v2 })

	InstructionSet[0xA5] = instrIfACmpEq
	InstructionSet[0xA6] = instrIfACmpNe

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
	InstructionSet[0xB9] = instrInvokeInterface

	InstructionSet[0xBB] = instrNew
	InstructionSet[0xBC] = instrNewArray
	InstructionSet[0xBD] = instrANewArray
	InstructionSet[0xBE] = instrArrayLength

	InstructionSet[0xBF] = instrAThrow

	InstructionSet[0xC0] = instrCheckCast
	InstructionSet[0xC1] = instrInstanceOf

	InstructionSet[0xC2] = instrMonitorEnter
	InstructionSet[0xC3] = instrMonitorExit

	InstructionSet[0xC6] = instrIfNull
	InstructionSet[0xC7] = instrIfNonNull
}

func ExecInstr(thread *Thread, frame *Frame, op byte) error {
	if InstructionSet[op] == nil {
		return fmt.Errorf("op(code = %#x) has been NOT implemented", op)
	}
	return InstructionSet[op](thread, frame)
}

func instrBiOp[T int | int64 | float32 | float64](name string, op func(T, T) T) Instruction {
	return func(_ *Thread, frame *Frame) error {
		v2, ok := frame.PopOperand().(T)
		if !ok {
			return fmt.Errorf("popped value2 for %s is invalid type", name)
		}
		v1, ok := frame.PopOperand().(T)
		if !ok {
			return fmt.Errorf("popped value1 for %s is invalid type", name)
		}

		frame.PushOperand(op(v1, v2))
		return nil
	}
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

func instrLConst(n int64) Instruction {
	return func(_ *Thread, frame *Frame) error {
		frame.PushOperand(n)
		return nil
	}
}

func instrFConst(n float32) Instruction {
	return func(_ *Thread, frame *Frame) error {
		frame.PushOperand(n)
		return nil
	}
}

func InstrBiPush(_ *Thread, frame *Frame) error {
	frame.PushOperand(int(int8(frame.NextParamByte())))
	return nil
}

func InstrSiPush(_ *Thread, frame *Frame) error {
	frame.PushOperand(int(int16(frame.NextParamUint16())))
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
			frame.PushOperand(NewInstance(class).SetVMData(
				frame.CurrentClass().File().ConstantPool().Utf8(uint16(entry)),
			))

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

func instrALoad(_ *Thread, frame *Frame) error {
	index := frame.PopOperand().(int)
	arrayref := frame.PopOperand().(*Array)

	frame.PushOperand(arrayref.Get(index))
	return nil
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

func InstrAStore(_ *Thread, frame *Frame) error {
	value := frame.PopOperand()
	index := frame.PopOperand().(int)
	arrayref := frame.PopOperand().(*Array)

	arrayref.Set(index, value)
	return nil
}

func InstrPop(n int) Instruction {
	return func(_ *Thread, frame *Frame) error {
		for i := 0; i < n; i++ {
			frame.PopOperand()
		}
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

func instrIInc(_ *Thread, frame *Frame) error {
	index := frame.NextParamByte()
	count := int(int8(frame.NextParamByte()))

	value := frame.Locals()[index].(int)
	frame.SetLocal(int(index), value+count)

	return nil
}

func InstrI2F(_ *Thread, frame *Frame) error {
	i, ok := frame.PopOperand().(int)
	if !ok {
		return fmt.Errorf("popped value for i2f is NOT int")
	}

	frame.PushOperand(float32(i))
	return nil
}

func InstrF2I(_ *Thread, frame *Frame) error {
	f, ok := frame.PopOperand().(float32)
	if !ok {
		return fmt.Errorf("popped value for f2i is NOT float32")
	}

	frame.PushOperand(int(f))
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

func instrLCmp(_ *Thread, frame *Frame) error {
	v2, ok := frame.PopOperand().(int64)
	if !ok {
		return fmt.Errorf("popped value2 for (i|l)cmp is invalid")
	}
	v1, ok := frame.PopOperand().(int64)
	if !ok {
		return fmt.Errorf("popped value1 for (i|l)cmp is invalid")
	}

	result := 0
	if v1 > v2 {
		result = 1
	} else if v1 < v2 {
		result = -1
	}
	frame.PushOperand(result)
	return nil
}

func instrFCmp(nanResult int) Instruction {
	return func(_ *Thread, frame *Frame) error {
		v2, ok := frame.PopOperand().(float32)
		if !ok {
			return fmt.Errorf("popped value2 for fcmp is NOT float32")
		}
		v1, ok := frame.PopOperand().(float32)
		if !ok {
			return fmt.Errorf("popped value1 for fcmp is NOT float32")
		}

		var result int
		if math.IsNaN(float64(v1)) || math.IsNaN(float64(v2)) {
			result = nanResult
		} else if v1 > v2 {
			result = 1
		} else if v1 == v2 {
			result = 0
		} else {
			result = -1
		}

		frame.PushOperand(result)
		return nil
	}
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

func instrIfACmpEq(_ *Thread, frame *Frame) error {
	branch := int16(frame.NextParamUint16())
	value2 := frame.PopOperand()
	value1 := frame.PopOperand()

	if value1 == nil || value2 == nil {
		if value1 == nil && value2 == nil {
			frame.JumpPC(uint16(int16(frame.PC()) + branch))
		}
		return nil
	}

	v1, ok := value1.(*Instance)
	if !ok {
		return fmt.Errorf("popped value2 for if_acmpeq is NOT instance")
	}

	v2, ok := value2.(*Instance)
	if !ok {
		return fmt.Errorf("popped value1 for if_acmpeq is NOT instance")
	}

	if v1 == v2 {
		frame.JumpPC(uint16(int16(frame.PC()) + branch))
	}
	return nil
}

func instrIfACmpNe(_ *Thread, frame *Frame) error {
	branch := int16(frame.NextParamUint16())
	value2 := frame.PopOperand()
	value1 := frame.PopOperand()

	if value1 == nil || value2 == nil {
		if !(value1 == nil && value2 == nil) {
			frame.JumpPC(uint16(int16(frame.PC()) + branch))
		}
		return nil
	}

	v1, ok := value1.(*Instance)
	if !ok {
		return fmt.Errorf("popped value2 for if_acmpeq is NOT instance")
	}

	v2, ok := value2.(*Instance)
	if !ok {
		return fmt.Errorf("popped value1 for if_acmpeq is NOT instance")
	}

	if v1 != v2 {
		frame.JumpPC(uint16(int16(frame.PC()) + branch))
	}
	return nil
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
	frame.PushOperand(resolvedClass.GetStaticField(resolvedField))

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
	resolvedClass.SetStaticField(resolvedField, frame.PopOperand())

	return nil
}

func instrGetField(_ *Thread, frame *Frame) error {
	_, name, desc := frame.curClass.File().ConstantPool().Reference(frame.NextParamUint16())
	instance := frame.PopOperand().(*Instance)
	if instance == nil {
		return fmt.Errorf("objectref for getfield is null")
	}

	frame.PushOperand(instance.GetField(name, desc))
	return nil
}

func instrPutField(_ *Thread, frame *Frame) error {
	_, name, desc := frame.curClass.File().ConstantPool().Reference(frame.NextParamUint16())
	value := frame.PopOperand()
	instance := frame.PopOperand().(*Instance)
	if instance == nil {
		return fmt.Errorf("objectref for getfield is null")
	}

	instance.PutField(name, desc, value)
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

func instrInvokeInterface(thread *Thread, frame *Frame) error {
	_, name, desc := frame.curClass.File().ConstantPool().Reference(frame.NextParamUint16())
	instance := frame.PeekFromTop(class_file.ParseDescriptor(desc)).(*Instance)
	if instance == nil {
		return fmt.Errorf("receiver instance is null")
	}

	frame.NextParamUint16() // Skip 'count' and '0'

	thread.DumpFrameStack(true)

	resolvedClass, resolvedMethod := instance.Class().ResolveMethod(*name, *desc)
	if resolvedClass == nil || !resolvedMethod.IsCallableForInstance() {
		return fmt.Errorf("method not found: %s.%s", *name, *desc)
	}

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

	fmt.Printf("****** new instance: %s\n", class.File().ThisClass())

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

func instrAThrow(thread *Thread, frame *Frame) error {
	thread.HandleException(frame.PopOperand().(*Instance))
	return nil
}

func instrCheckCast(thread *Thread, frame *Frame) error {
	className := frame.curClass.File().ConstantPool().ClassInfo(frame.NextParamUint16())

	objRef := frame.PopOperand()
	if objRef == nil {
		frame.PushOperand(nil)
		return nil
	}

	result := objRef
	switch o := objRef.(type) {
	case *Instance:
		if o.Class().IsInstanceOf(className) {
			_, state, err := thread.VM().FindInitializedClass(className, thread)
			if err != nil {
				return err
			}
			if state == FailedInitialization {
				return fmt.Errorf("failed initialization of waiting class: %s", *className)
			}
			//o.Cast(class)

			result = o
		}
	default:
		//return fmt.Errorf("instanceof does NOT support object is NOT class instance: %+v", o)
	}

	frame.PushOperand(result)
	return nil
}

func instrInstanceOf(_ *Thread, frame *Frame) error {
	className := frame.curClass.File().ConstantPool().ClassInfo(frame.NextParamUint16())

	objRef := frame.PopOperand()
	if objRef == nil {
		frame.PushOperand(0)
		return nil
	}

	result := 0
	switch o := objRef.(type) {
	case *Instance:
		if o.Class().IsInstanceOf(className) {
			result = 1
		}
	default:
		return fmt.Errorf("instanceof does NOT support object is NOT class instance")
	}

	frame.PushOperand(result)
	return nil
}

func instrMonitorEnter(thread *Thread, frame *Frame) error {
	frame.PopOperand().(*Instance).Monitor().Enter(thread)
	return nil
}

func instrMonitorExit(thread *Thread, frame *Frame) error {
	frame.PopOperand().(*Instance).Monitor().Exit()
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
