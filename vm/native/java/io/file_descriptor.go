package io

import "github.com/murakmii/gojiai/vm"

func init() {
	class := "java/io/FileDescriptor"

	vm.NativeMethods.Register(class, "initIDs", "()V", vm.NopNativeMethod)
}
