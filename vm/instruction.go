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
)

func init() {
	InstructionSet[0x00] = func(_ *Thread, _ *Frame) error { return nil } // nop
	InstructionSet[0x01] = instrAConstNull

	InstructionSet[0x02] = instrConst[int32](-1)
	InstructionSet[0x03] = instrConst[int32](0)
	InstructionSet[0x04] = instrConst[int32](1)
	InstructionSet[0x05] = instrConst[int32](2)
	InstructionSet[0x06] = instrConst[int32](3)
	InstructionSet[0x07] = instrConst[int32](4)
	InstructionSet[0x08] = instrConst[int32](5)

	InstructionSet[0x09] = instrConst[int64](0)
	InstructionSet[0x0A] = instrConst[int64](1)

	InstructionSet[0x0B] = instrConst[float32](0.0)
	InstructionSet[0x0C] = instrConst[float32](1.0)
	InstructionSet[0x0D] = instrConst[float32](2.0)

	InstructionSet[0x0E] = instrConst[float64](0.0)
	InstructionSet[0x0F] = instrConst[float64](1.0)

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

	InstructionSet[0x1E] = instrLoadN(0)
	InstructionSet[0x1F] = instrLoadN(1)
	InstructionSet[0x20] = instrLoadN(2)
	InstructionSet[0x21] = instrLoadN(3)

	InstructionSet[0x22] = instrLoadN(0)
	InstructionSet[0x23] = instrLoadN(1)
	InstructionSet[0x24] = instrLoadN(2)
	InstructionSet[0x25] = instrLoadN(3)

	InstructionSet[0x26] = instrLoadN(0)
	InstructionSet[0x27] = instrLoadN(1)
	InstructionSet[0x28] = instrLoadN(2)
	InstructionSet[0x29] = instrLoadN(3)

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

	InstructionSet[0x3F] = instrStoreN(0)
	InstructionSet[0x40] = instrStoreN(1)
	InstructionSet[0x41] = instrStoreN(2)
	InstructionSet[0x42] = instrStoreN(3)

	InstructionSet[0x47] = instrStoreN(0)
	InstructionSet[0x48] = instrStoreN(1)
	InstructionSet[0x49] = instrStoreN(2)
	InstructionSet[0x4A] = instrStoreN(3)

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

	InstructionSet[0x59] = instrDup
	InstructionSet[0x5A] = instrDupX1
	InstructionSet[0x5B] = instrDupX2
	InstructionSet[0x5C] = instrDup2
	InstructionSet[0x5D] = instrDup2X1
	InstructionSet[0x5E] = instrDup2X2

	InstructionSet[0x60] = instrAdd[int32]()
	InstructionSet[0x61] = instrAdd[int64]()

	InstructionSet[0x63] = instrBiOp[float64]("dadd", func(v1 float64, v2 float64) float64 { return v1 + v2 })
	InstructionSet[0x64] = instrBiOp[int32]("isub", func(v1 int32, v2 int32) int32 { return v1 - v2 })
	InstructionSet[0x65] = instrBiOp[int64]("lsub", func(v1 int64, v2 int64) int64 { return v1 - v2 })
	InstructionSet[0x67] = instrBiOp[float64]("dsub", func(v1 float64, v2 float64) float64 { return v1 - v2 })
	InstructionSet[0x68] = instrBiOp[int32]("imul", func(v1 int32, v2 int32) int32 { return v1 * v2 })
	InstructionSet[0x69] = instrBiOp[int64]("lmul", func(v1 int64, v2 int64) int64 { return v1 * v2 })
	InstructionSet[0x6A] = instrBiOp[float32]("fmul", func(v1 float32, v2 float32) float32 { return v1 * v2 })
	InstructionSet[0x6B] = instrBiOp[float64]("dmul", func(v1 float64, v2 float64) float64 { return v1 * v2 })
	InstructionSet[0x6C] = instrBiOp[int32]("idiv", func(v1 int32, v2 int32) int32 { return v1 / v2 })
	InstructionSet[0x6D] = instrBiOp[int64]("ldiv", func(v1 int64, v2 int64) int64 { return v1 / v2 })
	InstructionSet[0x6E] = instrBiOp[float32]("fdiv", func(v1 float32, v2 float32) float32 { return v1 / v2 })
	InstructionSet[0x70] = instrBiOp[int32]("irem", func(v1 int32, v2 int32) int32 { return v1 % v2 })
	InstructionSet[0x71] = instrBiOp[int64]("lrem", func(v1 int64, v2 int64) int64 { return v1 % v2 })

	InstructionSet[0x74] = instrINeg

	InstructionSet[0x78] = instrShiftLeft[int32](0x1F)
	InstructionSet[0x79] = instrShiftLeft[int64](0x3F)
	InstructionSet[0x7A] = instrShiftRight[int32](0x1F) // TODO: Arithmetic
	InstructionSet[0x7B] = instrShiftRight[int64](0x3F) // TODO: Arithmetic
	InstructionSet[0x7C] = instrShiftRight[int32](0x1F)
	InstructionSet[0x7D] = instrShiftRight[int64](0x3F)

	InstructionSet[0x7E] = instrAnd[int32]
	InstructionSet[0x7F] = instrAnd[int64]
	InstructionSet[0x80] = instrBiOp[int32]("ior", func(v1 int32, v2 int32) int32 { return v1 | v2 })
	InstructionSet[0x81] = instrBiOp[int64]("lor", func(v1 int64, v2 int64) int64 { return v1 | v2 })
	InstructionSet[0x82] = instrBiOp[int32]("ixor", func(v1 int32, v2 int32) int32 { return v1 ^ v2 })
	InstructionSet[0x83] = instrBiOp[int64]("lxor", func(v1 int64, v2 int64) int64 { return v1 ^ v2 })

	InstructionSet[0x84] = instrIInc

	InstructionSet[0x85] = InstrI2L
	InstructionSet[0x86] = InstrI2F
	InstructionSet[0x87] = InstrI2D
	InstructionSet[0x88] = InstrL2I
	InstructionSet[0x89] = InstrL2F

	InstructionSet[0x8B] = InstrF2I
	InstructionSet[0x8D] = InstrF2D
	InstructionSet[0x8E] = InstrD2I
	InstructionSet[0x8F] = InstrD2L
	InstructionSet[0x91] = instrI2B
	InstructionSet[0x92] = instrI2C
	InstructionSet[0x93] = instrI2S

	InstructionSet[0x94] = instrLCmp
	InstructionSet[0x95] = instrFCmp(-1)
	InstructionSet[0x96] = instrFCmp(1)

	InstructionSet[0x99] = instrIf(func(i int32) bool { return i == 0 })
	InstructionSet[0x9A] = instrIf(func(i int32) bool { return i != 0 })
	InstructionSet[0x9B] = instrIf(func(i int32) bool { return i < 0 })
	InstructionSet[0x9C] = instrIf(func(i int32) bool { return i >= 0 })
	InstructionSet[0x9D] = instrIf(func(i int32) bool { return i > 0 })
	InstructionSet[0x9E] = instrIf(func(i int32) bool { return i <= 0 })

	InstructionSet[0x9F] = instrIfICmp(func(v1 int32, v2 int32) bool { return v1 == v2 })
	InstructionSet[0xA0] = instrIfICmp(func(v1 int32, v2 int32) bool { return v1 != v2 })
	InstructionSet[0xA1] = instrIfICmp(func(v1 int32, v2 int32) bool { return v1 < v2 })
	InstructionSet[0xA2] = instrIfICmp(func(v1 int32, v2 int32) bool { return v1 >= v2 })
	InstructionSet[0xA3] = instrIfICmp(func(v1 int32, v2 int32) bool { return v1 > v2 })
	InstructionSet[0xA4] = instrIfICmp(func(v1 int32, v2 int32) bool { return v1 <= v2 })

	InstructionSet[0xA5] = instrIfACmpEq
	InstructionSet[0xA6] = instrIfACmpNe

	InstructionSet[0xA7] = instrGoTo

	InstructionSet[0xAA] = instrTableSwitch
	InstructionSet[0xAB] = instrLookupSwitch

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

	InstructionSet[0xC4] = instrWide

	InstructionSet[0xC6] = instrIfNull
	InstructionSet[0xC7] = instrIfNonNull
}

func ExecInstr(thread *Thread, frame *Frame, op byte) error {
	if InstructionSet[op] == nil {
		return fmt.Errorf("op(code = %#x) has been NOT implemented", op)
	}
	return InstructionSet[op](thread, frame)
}

func instrBiOp[T int32 | int64 | float32 | float64](name string, op func(T, T) T) Instruction {
	return func(thread *Thread, frame *Frame) error {
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

func instrConst[T int32 | int64 | float32 | float64](n T) Instruction {
	return func(_ *Thread, frame *Frame) error {
		frame.PushOperand(n)
		return nil
	}
}

func InstrBiPush(_ *Thread, frame *Frame) error {
	frame.PushOperand(int32(int8(frame.NextParamByte())))
	return nil
}

func InstrSiPush(_ *Thread, frame *Frame) error {
	frame.PushOperand(int32(int16(frame.NextParamUint16())))
	return nil
}

func instrLdc(idxLoader func(*Frame) uint16) Instruction {
	return func(thread *Thread, frame *Frame) error {
		switch entry := frame.CurrentClass().File().ConstantPool().Entry(idxLoader(frame)).(type) {
		case int32, float32, int64, float64:
			frame.PushOperand(entry)

		case class_file.StringCpInfo:
			js, err := thread.VM().JavaString(thread, frame.CurrentClass().File().ConstantPool().Utf8(uint16(entry)))
			if err != nil {
				return err
			}
			frame.PushOperand(js)

		case class_file.ClassCpInfo:
			name := frame.CurrentClass().File().ConstantPool().Utf8(uint16(entry))
			class, err := thread.VM().Class(*name, thread)
			if err != nil {
				return err
			}
			frame.PushOperand(class.Java())

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
	index := frame.PopOperand().(int32)
	frame.PushOperand(frame.PopOperand().(*Instance).AsArray()[index])
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
	index := frame.PopOperand().(int32)
	frame.PopOperand().(*Instance).AsArray()[index] = value

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

func instrDup(_ *Thread, frame *Frame) error {
	value := frame.PopOperand()
	frame.PushOperand(value)
	frame.PushOperand(value)

	return nil
}

func instrDupX1(_ *Thread, frame *Frame) error {
	v1 := frame.PopOperand()
	v2 := frame.PopOperand()

	frame.PushOperand(v1)
	frame.PushOperand(v2)
	frame.PushOperand(v1)

	return nil
}

func instrDupX2(_ *Thread, frame *Frame) error {
	v1 := frame.PopOperand()
	v2 := frame.PopOperand()

	switch v2.(type) {
	case int64, float64:
		frame.PushOperand(v1)
	default:
		v3 := frame.PopOperand()
		frame.PushOperand(v1)
		frame.PushOperand(v3)
	}

	frame.PushOperand(v2)
	frame.PushOperand(v1)

	return nil
}

func instrDup2(_ *Thread, frame *Frame) error {
	v1 := frame.PopOperand()

	switch v1.(type) {
	case int64, float64:
		frame.PushOperand(v1)
	default:
		v2 := frame.PopOperand()
		frame.PushOperand(v2)
		frame.PushOperand(v1)
		frame.PushOperand(v2)
	}

	frame.PushOperand(v1)
	return nil
}

func instrDup2X1(_ *Thread, frame *Frame) error {
	v1 := frame.PopOperand()
	v2 := frame.PopOperand()

	switch v1.(type) {
	case int64, float64:
		frame.PushOperand(v1)
	default:
		v3 := frame.PopOperand()
		frame.PushOperand(v2)
		frame.PushOperand(v1)
		frame.PushOperand(v3)
	}

	frame.PushOperand(v2)
	frame.PushOperand(v1)
	return nil
}

func instrDup2X2(_ *Thread, frame *Frame) error {
	v1 := frame.PopOperand()
	v2 := frame.PopOperand()

	v1cat2 := false
	switch v1.(type) {
	case int64, float64:
		v1cat2 = true
	}

	v2cat2 := false
	switch v2.(type) {
	case int64, float64:
		v2cat2 = true
	}

	if v1cat2 && v2cat2 {
		frame.PushOperand(v1)
		frame.PushOperand(v2)
		frame.PushOperand(v1)
		return nil
	}

	v3 := frame.PopOperand()

	if v1cat2 {
		frame.PushOperand(v1)
		frame.PushOperand(v3)
		frame.PushOperand(v2)
		frame.PushOperand(v1)
		return nil
	}

	switch v3.(type) {
	case int64, float64:
		frame.PushOperand(v2)
		frame.PushOperand(v1)
		frame.PushOperand(v3)
		frame.PushOperand(v2)
		frame.PushOperand(v1)
		return nil
	}

	v4 := frame.PopOperand()

	frame.PushOperand(v2)
	frame.PushOperand(v1)
	frame.PushOperand(v4)
	frame.PushOperand(v3)
	frame.PushOperand(v2)
	frame.PushOperand(v1)

	return nil
}

func instrAdd[T int32 | int64]() Instruction {
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

func instrINeg(_ *Thread, frame *Frame) error {
	frame.PushOperand(-1 * frame.PopOperand().(int32))
	return nil
}

func instrShiftLeft[T int32 | int64](mask int32) Instruction {
	return func(_ *Thread, frame *Frame) error {
		v2, ok := frame.PopOperand().(int32)
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

func instrShiftRight[T int32 | int64](mask int32) Instruction {
	return func(thread *Thread, frame *Frame) error {
		v2, ok := frame.PopOperand().(int32)
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

func instrAnd[T int32 | int64](thread *Thread, frame *Frame) error {
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
	count := int32(int8(frame.NextParamByte()))

	value := frame.Locals()[index].(int32)
	frame.SetLocal(int(index), value+count)

	return nil
}

func InstrI2F(_ *Thread, frame *Frame) error {
	i, ok := frame.PopOperand().(int32)
	if !ok {
		return fmt.Errorf("popped value for i2f is NOT int")
	}

	frame.PushOperand(float32(i))
	return nil
}

func InstrI2D(_ *Thread, frame *Frame) error {
	i, ok := frame.PopOperand().(int32)
	if !ok {
		return fmt.Errorf("popped value for i2d is NOT int")
	}

	frame.PushOperand(float64(i))
	return nil
}

func InstrL2I(_ *Thread, frame *Frame) error {
	i, ok := frame.PopOperand().(int64)
	if !ok {
		return fmt.Errorf("popped value for l2i is NOT int64")
	}

	frame.PushOperand(int32(i))
	return nil
}

func InstrL2F(_ *Thread, frame *Frame) error {
	i, ok := frame.PopOperand().(int64)
	if !ok {
		return fmt.Errorf("popped value for l2f is NOT int64")
	}

	frame.PushOperand(float32(i))
	return nil
}

func InstrF2I(_ *Thread, frame *Frame) error {
	f, ok := frame.PopOperand().(float32)
	if !ok {
		return fmt.Errorf("popped value for f2i is NOT float32")
	}

	frame.PushOperand(int32(f))
	return nil
}

func InstrF2D(_ *Thread, frame *Frame) error {
	f, ok := frame.PopOperand().(float32)
	if !ok {
		return fmt.Errorf("popped value for f2d is NOT float32")
	}

	frame.PushOperand(float64(f))
	return nil
}

func InstrD2I(_ *Thread, frame *Frame) error {
	d, ok := frame.PopOperand().(float64)
	if !ok {
		return fmt.Errorf("popped value for d2i is NOT float64")
	}

	frame.PushOperand(int32(d))
	return nil
}

func InstrD2L(_ *Thread, frame *Frame) error {
	f, ok := frame.PopOperand().(float64)
	if !ok {
		return fmt.Errorf("popped value for d2l is NOT float64")
	}

	frame.PushOperand(int64(f))
	return nil
}

func instrI2B(_ *Thread, frame *Frame) error {
	i, ok := frame.PopOperand().(int32)
	if !ok {
		return fmt.Errorf("popped value for i2b is NOT int")
	}

	frame.PushOperand(i & 0xFF)
	return nil
}

func instrI2C(_ *Thread, frame *Frame) error {
	i, ok := frame.PopOperand().(int32)
	if !ok {
		return fmt.Errorf("popped value for f2i is NOT float32")
	}

	frame.PushOperand(i & 0xFFFF)
	return nil
}

func instrI2S(_ *Thread, frame *Frame) error {
	i, ok := frame.PopOperand().(int32)
	if !ok {
		return fmt.Errorf("popped value for i2s is NOT int")
	}

	frame.PushOperand(i & 0xFFFF)
	return nil
}

func InstrI2L(_ *Thread, frame *Frame) error {
	i, ok := frame.PopOperand().(int32)
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

	var result int32
	if v1 > v2 {
		result = 1
	} else if v1 < v2 {
		result = -1
	}
	frame.PushOperand(result)
	return nil
}

func instrFCmp(nanResult int32) Instruction {
	return func(_ *Thread, frame *Frame) error {
		v2, ok := frame.PopOperand().(float32)
		if !ok {
			return fmt.Errorf("popped value2 for fcmp is NOT float32")
		}
		v1, ok := frame.PopOperand().(float32)
		if !ok {
			return fmt.Errorf("popped value1 for fcmp is NOT float32")
		}

		var result int32
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

func instrIf(matcher func(int32) bool) Instruction {
	return func(thread *Thread, frame *Frame) error {
		branch := int16(frame.NextParamUint16())
		value, ok := frame.PopOperand().(int32)
		if !ok {
			return fmt.Errorf("popped value for if<cond> is NOT int")
		}

		if matcher(value) {
			frame.JumpPC(uint16(int16(frame.PC()) + branch))
		}
		return nil
	}
}

func instrIfICmp(comparator func(int32, int32) bool) Instruction {
	return func(thread *Thread, frame *Frame) error {
		branch := int16(frame.NextParamUint16())
		v2, ok := frame.PopOperand().(int32)
		if !ok {
			return fmt.Errorf("popped value2 for if_icmp<cond> is NOT int")
		}
		v1, ok := frame.PopOperand().(int32)
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

func instrIfACmpNe(thread *Thread, frame *Frame) error {
	branch := int16(frame.NextParamUint16())
	value2 := frame.PopOperand()
	value1 := frame.PopOperand()

	/*if value1 == nil || value2 == nil {
		if !(value1 == nil && value2 == nil) {
			frame.JumpPC(uint16(int16(frame.PC()) + branch))
		}
		return nil
	}

	fmt.Printf("if_acmpeq %s:%s\n", frame.CurrentClass().File().ThisClass(), *frame.CurrentMethod().Name())

	v1, ok := value1.(*Instance)
	if !ok {
		return fmt.Errorf("popped value1 for if_acmpeq is NOT instance: %+v", value1)
	}

	v2, ok := value2.(*Instance)
	if !ok {
		return fmt.Errorf("popped value2 for if_acmpeq is NOT instance")
	}

	if v1 != v2 {
		frame.JumpPC(uint16(int16(frame.PC()) + branch))
	}*/

	if value1 != value2 {
		frame.JumpPC(uint16(int16(frame.PC()) + branch))
	}
	return nil
}

func instrGoTo(_ *Thread, frame *Frame) error {
	branch := int16(frame.NextParamUint16())
	frame.JumpPC(uint16(int16(frame.PC()) + branch))
	return nil
}

func instrTableSwitch(thread *Thread, frame *Frame) error {
	frame.NextAlign(4)

	defaultVal := int32(frame.NextParamUint32())
	low := int32(frame.NextParamUint32())
	high := int32(frame.NextParamUint32())

	index := frame.PopOperand().(int32)
	if index < low || index > high {
		frame.JumpPC(frame.PC() + uint16(defaultVal))
		return nil
	}

	offsets := make([]int, high-low+1)
	for i := range offsets {
		offsets[i] = int(frame.NextParamUint32())
	}

	frame.JumpPC(frame.PC() + uint16(offsets[index-low]))
	return nil
}

func instrLookupSwitch(_ *Thread, frame *Frame) error {
	frame.NextAlign(4)

	defaultVal := int32(frame.NextParamUint32())
	key := frame.PopOperand().(int32)

	for npairs := frame.NextParamUint32(); npairs > 0; npairs-- {
		match := int32(frame.NextParamUint32())
		offset := int32(frame.NextParamUint32())

		if match == key {
			frame.JumpPC(frame.PC() + uint16(offset))
			return nil
		}
	}

	frame.JumpPC(frame.PC() + uint16(defaultVal))
	return nil
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

func instrGetStatic(thread *Thread, frame *Frame) error {
	className, name, desc := frame.CurrentClass().File().ConstantPool().Reference(frame.NextParamUint16())
	class, err := thread.VM().Class(*className, thread)
	if err != nil {
		return err
	}

	resolvedClass, resolvedField := class.ResolveField(*name, *desc)
	frame.PushOperand(resolvedClass.GetStaticField(resolvedField))

	return nil
}

func instrPutStatic(thread *Thread, frame *Frame) error {
	className, name, desc := frame.CurrentClass().File().ConstantPool().Reference(frame.NextParamUint16())
	class, err := thread.VM().Class(*className, thread)
	if err != nil {
		return err
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

	switch i := frame.PeekFromTop(len(class_file.MethodDescriptor(*desc).Params())).(type) {
	case *Instance:
		resolvedClass, resolvedMethod := i.Class().ResolveMethod(*name, *desc)
		if resolvedClass == nil || !resolvedMethod.IsCallableForInstance() {
			return fmt.Errorf("method not found: %s.%s", *name, *desc)
		}
		return thread.ExecMethod(resolvedClass, resolvedMethod)

	default:
		fmt.Printf("operand stack: %+v\nPC: %d\n", frame.opStack, frame.PC())

		return fmt.Errorf("invokevirtual: receiver is invalid object: %+v", i)
	}
}

func instrInvokeSpecial(thread *Thread, frame *Frame) error {
	className, name, desc := frame.curClass.File().ConstantPool().Reference(frame.NextParamUint16())
	class, err := thread.VM().Class(*className, thread)
	if err != nil {
		return err
	}

	resolvedClass, resolvedMethod := class.ResolveMethod(*name, *desc)
	if resolvedClass == nil || !resolvedMethod.IsCallableForInstance() {
		return fmt.Errorf("method not found: %s.%s", *name, *desc)
	}

	return thread.ExecMethod(resolvedClass, resolvedMethod)
}

func instrInvokeStatic(thread *Thread, frame *Frame) error {
	className, name, desc := frame.curClass.File().ConstantPool().Reference(frame.NextParamUint16())
	class, err := thread.VM().Class(*className, thread)
	if err != nil {
		return err
	}

	resolvedClass, resolvedMethod := class.ResolveMethod(*name, *desc)
	if resolvedClass == nil || !resolvedMethod.IsCallableAsStatic() {
		return fmt.Errorf("method not found: %s.%s", *name, *desc)
	}

	return thread.ExecMethod(resolvedClass, resolvedMethod)
}

func instrInvokeInterface(thread *Thread, frame *Frame) error {
	_, name, desc := frame.curClass.File().ConstantPool().Reference(frame.NextParamUint16())
	instance := frame.PeekFromTop(len(class_file.MethodDescriptor(*desc).Params())).(*Instance)
	if instance == nil {
		return fmt.Errorf("receiver instance is null")
	}

	frame.NextParamUint16() // Skip 'count' and '0'

	resolvedClass, resolvedMethod := instance.Class().ResolveMethod(*name, *desc)
	if resolvedClass == nil || !resolvedMethod.IsCallableForInstance() {
		return fmt.Errorf("method not found: %s.%s", *name, *desc)
	}

	return thread.ExecMethod(resolvedClass, resolvedMethod)
}

func instrNew(thread *Thread, frame *Frame) error {
	className := frame.curClass.File().ConstantPool().ClassInfo(frame.NextParamUint16())
	class, err := thread.VM().Class(*className, thread)
	if err != nil {
		return err
	}

	frame.PushOperand(NewInstance(class))
	return nil
}

var typeCodes = []string{"Z", "C", "F", "D", "B", "S", "I", "J"}

func instrNewArray(thread *Thread, frame *Frame) error {
	arrayClass := "[" + typeCodes[frame.NextParamByte()-4]
	array, _ := NewArray(thread.VM(), arrayClass, int(frame.PopOperand().(int32)))
	frame.PushOperand(array)
	return nil
}

func instrANewArray(thread *Thread, frame *Frame) error {
	className := *(frame.curClass.File().ConstantPool().ClassInfo(frame.NextParamUint16()))
	if className[0] != '[' {
		className = "L" + className + ";"
	}

	array, _ := NewArray(thread.VM(), "["+className, int(frame.PopOperand().(int32)))
	frame.PushOperand(array)
	return nil
}

func instrArrayLength(_ *Thread, frame *Frame) error {
	frame.PushOperand(int32(len(frame.PopOperand().(*Instance).AsArray())))
	return nil
}

func instrAThrow(_ *Thread, frame *Frame) error {
	return NewJavaErr(frame.PopOperand().(*Instance))
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
			_, err := thread.VM().Class(*className, thread)
			if err != nil {
				return err
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
		frame.PushOperand(int32(0))
		return nil
	}

	var result int32
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
	frame.PopOperand().(*Instance).Monitor().Enter(thread, -1)
	return nil
}

func instrMonitorExit(thread *Thread, frame *Frame) error {
	frame.PopOperand().(*Instance).Monitor().Exit(thread)
	return nil
}

func instrWide(_ *Thread, frame *Frame) error {
	op := frame.NextParamByte()
	index := frame.NextParamUint16()

	if op == 0x84 { // iinc
		inc := int32(int16(frame.NextParamUint16()))
		frame.Locals()[index] = frame.Locals()[index].(int32) + inc
		return nil
	} else if op >= 0x15 && op <= 0x19 { // (i|f|a|l|d)load
		frame.PushOperand(frame.Locals()[index])
	} else if op >= 0x36 && op <= 0x3A { // (i|f|a|l|d)store
		frame.Locals()[index] = frame.PopOperand()
	} else {
		return fmt.Errorf("unknown wide op = %#x", op)
	}

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
