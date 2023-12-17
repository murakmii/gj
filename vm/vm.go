package vm

import (
	"github.com/murakmii/gojiai/class_file"
	"os"
)

type (
	VM struct {
		thread *Thread
	}
)

func NewVM() *VM {
	return &VM{thread: NewThread()}
}

// classファイルのパス、実行するメソッドの名前とシグネチャ、メソッドの引数を受け取り、
// そのメソッドを実行した後戻り値を返す。
func (v *VM) Execute(
	classFilePath, methodName, methodDesc string,
	args []interface{},
) (interface{}, error) {

	f, err := os.Open(classFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// classファイルのデコード
	classFile, err := class_file.ReadClassFile(f)
	if err != nil {
		return nil, err
	}

	// 呼び出し元となるフレームをカレントフレームとしておく
	// このフレームのオペランドスタックに呼び出すメソッドの引数をプッシュする
	invokerFrame := &Frame{}
	for _, arg := range args {
		invokerFrame.PushOperand(arg)
	}
	v.thread.PushFrame(invokerFrame)

	// 呼び出すメソッドをclassファイルから探し、実行
	method := classFile.FindMethod(methodName, methodDesc)
	if err = v.thread.ExecMethod(method); err != nil {
		return nil, err
	}

	// 戻り値は呼び出し元フレームのオペランドスタックにプッシュされているはずなので、それを返す
	return invokerFrame.PopOperand(), nil
}
