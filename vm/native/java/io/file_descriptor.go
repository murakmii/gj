package io

import "github.com/murakmii/gj/vm"

func init() {
	class := "java/io/FileDescriptor"

	vm.NativeMethods.Register(class, "initIDs", "()V", vm.NopNativeMethod)
}
