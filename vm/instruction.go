package vm

import (
	"fmt"
)

type (
	// 1命令に対応する関数を表す型
	Instruction func(*Thread)
)

var (
	// オペコードに対応するインデックスに、対応する処理を実装した関数を登録する
	InstructionSet [256]Instruction
)

func init() {
	InstructionSet[0x00] = instrNop

	InstructionSet[0x03] = instrConst(0)
	InstructionSet[0x04] = instrConst(1)

	InstructionSet[0x1A] = instrILoad(0)
	InstructionSet[0x1B] = instrILoad(1)
	InstructionSet[0x1C] = instrILoad(2)

	InstructionSet[0x3C] = instrIStore(1)
	InstructionSet[0x3D] = instrIStore(2)

	InstructionSet[0x60] = instrIAdd

	InstructionSet[0x84] = instrIInc

	InstructionSet[0xA3] = instrIfICmpGt

	InstructionSet[0xA7] = instrGoTo

	InstructionSet[0xAC] = instrIReturn
}

// スレッド thread のカレントフレーム上で1命令実行する
func ExecInstr(thread *Thread) error {
	// オペコード読み取り
	opCode := thread.CurrentFrame().NextInstr()

	if InstructionSet[opCode] == nil {
		return fmt.Errorf("op(code = %#x) has been NOT implemented", opCode)
	}

	// オペコードから命令に対応する関数をひいて実行
	InstructionSet[opCode](thread)
	return nil
}

func instrNop(_ *Thread) {}

// iconst_0, iconst1
func instrConst(n int32) Instruction {
	// オペランドスタックに n をプッシュする Instruction を返す
	return func(thread *Thread) {
		thread.CurrentFrame().PushOperand(n)
	}
}

// iload_0, iload_1, iload_2
func instrILoad(n int) Instruction {
	// オペランドスタックに、
	// ローカル変数領域の n 番目の値をプッシュする Instruction を返す
	return func(thread *Thread) {
		frame := thread.CurrentFrame()
		frame.PushOperand(frame.Locals()[n])
	}
}

// istore_1, istore_2
func instrIStore(n int) Instruction {
	// オペランドスタックから値をポップし、
	// その値をローカル変数領域の n 番目に格納する Instruction を返す
	return func(thread *Thread) {
		frame := thread.CurrentFrame()
		frame.SetLocal(n, frame.PopOperand())
	}
}

// iadd
func instrIAdd(thread *Thread) {
	frame := thread.CurrentFrame()
	v2 := frame.PopOperand().(int32)
	v1 := frame.PopOperand().(int32)

	frame.PushOperand(v1 + v2)
}

// iinc
func instrIInc(thread *Thread) {
	frame := thread.CurrentFrame()

	// 対象のローカル変数領域のインデックス番号をオペランド空読み取り
	index := frame.NextParamByte()

	// 増分も同様にオペランドから取得
	count := int32(int8(frame.NextParamByte()))

	value := frame.Locals()[index].(int32)
	frame.SetLocal(int(index), value+count) // 増分を加算
}

// if_icmpgt
func instrIfICmpGt(thread *Thread) {
	frame := thread.CurrentFrame()

	// プログラムカウンタの移動量はオペランドより
	branch := int16(frame.NextParamUint16())

	// 不等式の右辺と左辺をオペランドスタックより
	v2 := frame.PopOperand().(int32)
	v1 := frame.PopOperand().(int32)

	// 左辺の値が右辺の値より大きい場合に、現在のプログラムカウンタを
	// オペランドの値分だけ移動させる
	if v1 > v2 {
		frame.JumpPC(uint16(int16(frame.PC()) + branch))
	}
}

// goto
func instrGoTo(thread *Thread) {
	frame := thread.CurrentFrame()
	branch := int16(frame.NextParamUint16())
	frame.JumpPC(uint16(int16(frame.PC()) + branch))
}

// ireturn
func instrIReturn(thread *Thread) {
	// オペランドスタックからポップしてから破棄
	retVal := thread.CurrentFrame().PopOperand()
	thread.PopFrame()

	// カレントフレームは呼び出し元メソッドのものになっているため、
	// ポップしておいた値をプッシュする
	thread.CurrentFrame().PushOperand(retVal)
}
