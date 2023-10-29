package io

import (
	"github.com/murakmii/gojiai/vm"
	"os"
	"syscall"
)

const (
	// See: https://github.com/openjdk/jdk/blob/jdk8-b120/jdk/src/share/classes/java/io/FileSystem.java#L102
	ufsBAExists    = int32(0x01)
	ufsBARegular   = int32(0x02)
	ufsBADirectory = int32(0x04)
)

func init() {
	class := "java/io/UnixFileSystem"

	vm.NativeMethods.Register(class, "initIDs", "()V", vm.NopNativeMethod)

	vm.NativeMethods.Register(class, "canonicalize0", "(Ljava/lang/String;)Ljava/lang/String;", func(thread *vm.Thread, args []interface{}) error {
		thread.CurrentFrame().PushOperand(args[1]) // nop
		return nil
	})

	vm.NativeMethods.Register(class, "checkAccess", "(Ljava/io/File;I)Z", func(thread *vm.Thread, args []interface{}) error {
		file := args[1].(*vm.Instance)
		path := file.GetField("path", "Ljava/lang/String;").(*vm.Instance).AsString()

		ret := int32(1)
		if err := syscall.Access(path, uint32(args[2].(int32))); err != nil {
			ret = 0
		}

		thread.CurrentFrame().PushOperand(ret)
		return nil
	})

	vm.NativeMethods.Register(class, "getBooleanAttributes0", "(Ljava/io/File;)I", func(thread *vm.Thread, args []interface{}) error {
		file := args[1].(*vm.Instance)
		path := file.GetField("path", "Ljava/lang/String;").(*vm.Instance).AsString()

		stat, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				thread.CurrentFrame().PushOperand(int32(0))
				return nil
			}
			return err
		}

		ba := ufsBAExists
		if stat.IsDir() {
			ba |= ufsBADirectory
		} else {
			ba |= ufsBARegular
		}

		thread.CurrentFrame().PushOperand(ba)
		return nil
	})
}
