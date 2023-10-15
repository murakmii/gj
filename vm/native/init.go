package native

import (
	"github.com/murakmii/gj/vm"
	"github.com/murakmii/gj/vm/native/java/lang"
)

func init() {
	vm.RegisterNativeMethod("java/lang/Object/registerNatives", lang.ObjectRegisterNatives)
}
