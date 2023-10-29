package misc

import "github.com/murakmii/gj/vm"

func init() {
	class := "sun/misc/VM"

	vm.NativeMethods.Register(class, "initialize", "()V", vm.NopNativeMethod)
}
