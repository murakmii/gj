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
	thread.CurrentFrame().PushOperand(1)
	return nil
}

func FileInputStreamClose0(_ *vm.Thread, args []interface{}) error {
	fis := args[0].(*vm.Instance)

	fdName := "fd"
	fdDesc := "Ljava/io/FileDescriptor;"
	if err := fis.GetField(&fdName, &fdDesc).(*vm.Instance).AsFile().Close(); err != nil {
		return err
	}

	return nil
}

func FileInputStreamOpen0(thread *vm.Thread, args []interface{}) error {
	name := args[1].(*vm.Instance).GetCharArrayField("value")
	file, err := os.Open(name)
	if err != nil {
		return err
	}

	fdClass, err := thread.VM().Class("java/io/FileDescriptor", thread)
	if err != nil {
		return err
	}
	_, constr := fdClass.ResolveMethod("<init>", "(I)V")

	fd := vm.NewInstance(fdClass)
	if err = thread.Execute(vm.NewFrame(fdClass, constr).SetLocals([]interface{}{fd, int(file.Fd())})); err != nil {
		return err
	}
	fd.SetVMData(file)

	fis := args[0].(*vm.Instance)

	fdName := "fd"
	fdDesc := "Ljava/io/FileDescriptor;"
	fis.PutField(&fdName, &fdDesc, fd)

	return nil
}

func FileInputStreamReadBytes(thread *vm.Thread, args []interface{}) error {
	fdName := "fd"
	fdDesc := "Ljava/io/FileDescriptor;"
	file := args[0].(*vm.Instance).GetField(&fdName, &fdDesc).(*vm.Instance).AsFile()
	dst := args[1].(*vm.Instance).AsArray()
	off := args[2].(int)
	size := args[3].(int)

	buf := make([]byte, size)
	n, err := file.Read(buf)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
		n = -1
	}

	for i := 0; i < n; i++ {
		dst[off+i] = int(buf[i])
	}

	thread.CurrentFrame().PushOperand(n)
	return nil
}
