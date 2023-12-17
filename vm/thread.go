package vm

import (
	"github.com/murakmii/gojiai/class_file"
)

type (
	// スレッドを表す構造体
	Thread struct {
		frames []*Frame
	}
)

func NewThread() *Thread {
	return &Thread{}
}

// フレームの生成(新たなカレントフレームを設定)
func (thread *Thread) PushFrame(frame *Frame) {
	thread.frames = append(thread.frames, frame)
}

// カレントフレームの破棄
func (thread *Thread) PopFrame() {
	thread.frames = thread.frames[:len(thread.frames)-1]
}

// カレントフレームを参照
func (thread *Thread) CurrentFrame() *Frame {
	return thread.frames[len(thread.frames)-1]
}

// 指定のメソッドを実行するためのフレーム操作を行う
func (thread *Thread) ExecMethod(method *class_file.MethodInfo) error {
	invokerFrame := thread.CurrentFrame() // メソッド呼び出し前のカレントフレーム

	// メソッド呼び出し前のカレントフレームのオペランドスタックから引数分だけ値をポップし
	// ローカル変数領域へ。
	// その後カレントフレームを切り替え。
	thread.PushFrame(
		NewFrame(method).SetLocals(
			invokerFrame.PopOperands(method.NumArgs()),
		),
	)

	for {
		curFrame := thread.CurrentFrame()
		if curFrame == invokerFrame {
			break
		}

		if err := ExecInstr(thread); err != nil {
			return err
		}
	}

	return nil
}
