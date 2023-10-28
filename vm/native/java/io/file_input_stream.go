package io

import (
	"errors"
	"github.com/murakmii/gj/vm"
	"io"
	"os"
)

// Return value of FileInputStream.available0 is approximate.
// So, This native implementation always returns 1.
func FileInputStreamAvailable0(thread *vm.Thread, args []interface{}) error {
	thread.CurrentFrame().PushOperand(int32(1))
	return nil
}

func FileInputStreamClose0(_ *vm.Thread, args []interface{}) error {
	fis := args[0].(*vm.Instance)
	if err := fis.GetField("fd", "Ljava/io/FileDescriptor;").(*vm.Instance).AsFile().Close(); err != nil {
		return err
	}

	return nil
}

func FileInputStreamOpen0(thread *vm.Thread, args []interface{}) error {
	file, err := os.Open(args[1].(*vm.Instance).AsString())
	if err != nil {
		return err
	}

	fdClass, err := thread.VM().Class("java/io/FileDescriptor", thread)
	if err != nil {
		return err
	}
	_, constr := fdClass.ResolveMethod("<init>", "(I)V")

	fd := vm.NewInstance(fdClass)
	if err = thread.Execute(vm.NewFrame(fdClass, constr).SetLocals([]interface{}{fd, int32(file.Fd())})); err != nil {
		return err
	}
	fd.SetVMData(file)

	fis := args[0].(*vm.Instance)
	fis.PutField("fd", "Ljava/io/FileDescriptor;", fd)

	return nil
}

func FileInputStreamReadBytes(thread *vm.Thread, args []interface{}) error {
	file := args[0].(*vm.Instance).GetField("fd", "Ljava/io/FileDescriptor;").(*vm.Instance).AsFile()
	dst := args[1].(*vm.Instance).AsArray()
	off := int(args[2].(int32))
	size := args[3].(int32)

	buf := make([]byte, size)
	n, err := file.Read(buf)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
		n = -1
	}

	for i := 0; i < n; i++ {
		dst[off+i] = int32(buf[i])
	}

	thread.CurrentFrame().PushOperand(int32(n))
	return nil
}
