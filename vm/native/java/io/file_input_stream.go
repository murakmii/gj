package io

import (
	"fmt"
	"github.com/murakmii/gj/vm"
)

func FileInputStreamInitIDs(thread *vm.Thread, args []interface{}) error {
	fmt.Println("execute java/io/FileInputStream.initIDs")
	return nil
}
