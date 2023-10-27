package io

import (
	"github.com/murakmii/gj/vm"
	"os"
)

const (
	// See: https://github.com/openjdk/jdk/blob/jdk8-b120/jdk/src/share/classes/java/io/FileSystem.java#L102
	ufsBAExists    = int32(0x01)
	ufsBARegular   = int32(0x02)
	ufsBADirectory = int32(0x04)
)

func UnixFileSystemCanonicalize0(thread *vm.Thread, args []interface{}) error {
	thread.CurrentFrame().PushOperand(args[1]) // nop
	return nil
}

func UnixFileSystemGetBooleanAttributes0(thread *vm.Thread, args []interface{}) error {
	file := args[1].(*vm.Instance)

	pathName := "path"
	pathDesc := "Ljava/lang/String;"
	path := file.GetField(&pathName, &pathDesc).(*vm.Instance).GetCharArrayField("value")

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
}
