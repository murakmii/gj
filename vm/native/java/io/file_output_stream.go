package io

import (
	"github.com/murakmii/gj/vm"
	"os"
)

func FileOutputStreamWriteBytes(thread *vm.Thread, args []interface{}) error {
	fos := args[0].(*vm.Instance)
	bytes := vm.JavaByteArrayToGo(args[1].(*vm.Instance), args[2].(int), args[3].(int))

	fdFieldName := "fd"

	descriptorFieldDesc := "Ljava/io/FileDescriptor;"
	descriptor := fos.GetField(&fdFieldName, &descriptorFieldDesc).(*vm.Instance)

	fdFieldDesc := "I"
	fd := descriptor.GetField(&fdFieldName, &fdFieldDesc).(int)

	_, err := os.NewFile(uintptr(fd), "").Write(bytes)
	return err
}
