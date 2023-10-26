package io

import (
	"github.com/murakmii/gj/vm"
)

func FileOutputStreamWriteBytes(thread *vm.Thread, args []interface{}) error {
	fos := args[0].(*vm.Instance)
	bytes := vm.JavaByteArrayToGo(args[1].(*vm.Instance), args[2].(int), args[3].(int))

	fdFieldName := "fd"

	descriptorFieldDesc := "Ljava/io/FileDescriptor;"
	descriptor := fos.GetField(&fdFieldName, &descriptorFieldDesc).(*vm.Instance)

	_, err := descriptor.AsFile().Write(bytes)
	return err
}
