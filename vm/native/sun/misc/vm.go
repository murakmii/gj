package misc

import "github.com/murakmii/gojiai/vm"

func init() {
	class := "sun/misc/VM"

	vm.NativeMethods.Register(class, "initialize", "()V", vm.NopNativeMethod)
}
