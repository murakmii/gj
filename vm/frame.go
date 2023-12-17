package vm

import (
	"bytes"
	"github.com/murakmii/gojiai/class_file"
	"github.com/murakmii/gojiai/support"
)

type (
	// フレームを表す構造体
	Frame struct {
		locals  []interface{}    // ローカル変数領域
		opStack []interface{}    // オペランドスタック
		code    *support.ByteSeq // 実行しているメソッドの命令
		pc      uint16           // プログラムカウンタ
	}
)

// 呼び出すメソッドに対応するフレームを生成
func NewFrame(method *class_file.MethodInfo) *Frame {
	// 本誌で扱う範囲では、必要な情報はCode属性のみ
	code := method.Code()
	codeReader, _ := support.NewByteSeq(bytes.NewReader(code.Code()))

	return &Frame{
		// ローカル変数領域の長さはmax_localsフィールドの値より
		locals:  make([]interface{}, code.MaxLocals()),
		opStack: nil, // オペランドスタックは空から始まる
		code:    codeReader,
		pc:      0, // プログラムカウンタは命令列の先頭を指す
	}
}

// ローカル変数領域のインデックスを指定して、値を格納する
func (frame *Frame) SetLocal(index int, v interface{}) *Frame {
	frame.locals[index] = v
	return frame
}

func (frame *Frame) SetLocals(vars []interface{}) *Frame {
	for i, v := range vars {
		frame.locals[i] = v
	}

	return frame
}

// ローカル変数領域を参照する
func (frame *Frame) Locals() []interface{} {
	return frame.locals
}

// 次の命令を返し、プログラムカウンタを更新する
func (frame *Frame) NextInstr() byte {
	frame.pc = uint16(frame.code.Pos())
	return frame.code.ReadByte()
}

// 命令中の次の1バイトを読み取る
func (frame *Frame) NextParamByte() byte {
	return frame.code.ReadByte()
}

// 命令中の次の2バイトを非負2バイト整数として読み取る
func (frame *Frame) NextParamUint16() uint16 {
	return frame.code.ReadUint16()
}

// 現在のプログラムカウンタを返す
func (frame *Frame) PC() uint16 {
	return frame.pc
}

// プログラムカウンタを移動させる
func (frame *Frame) JumpPC(pc uint16) {
	frame.pc = pc
	frame.code.Seek(int(pc))
}

// オペランドスタックに値をプッシュする
func (frame *Frame) PushOperand(value interface{}) {
	frame.opStack = append(frame.opStack, value)
}

// オペランドスタックから値をポップする
func (frame *Frame) PopOperand() interface{} {
	last := len(frame.opStack) - 1
	pop := frame.opStack[last]

	frame.opStack = frame.opStack[:last]

	return pop
}

// オペランドスタックから n 個だけポップしてスライスで返す
func (frame *Frame) PopOperands(n int) []interface{} {
	popped := make([]interface{}, n)
	for i := n - 1; i >= 0; i-- {
		popped[i] = frame.PopOperand()
	}
	return popped
}
