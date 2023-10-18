package io

import (
	"fmt"
	"github.com/murakmii/gj/vm"
)

func FileDescriptorInitIDs(thread *vm.Thread, args []interface{}) error {
	fmt.Println("execute java/io/FileDescriptor.initIDs")
	return nil
}
