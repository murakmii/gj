package lang

import (
	"fmt"
	"github.com/murakmii/gj/vm"
)

func SystemRegisterNatives(thread *vm.Thread, args []interface{}) error {
	fmt.Println("execute java/lang/System.registerNatives")
	return nil
}

func SystemArrayCopy(thread *vm.Thread, args []interface{}) error {
	src := args[0].(*vm.Array)
	srcStart := args[1].(int)
	dst := args[2].(*vm.Array)
	dstStart := args[3].(int)
	count := args[4].(int)

	for i := 0; i < count; i++ {
		dst.Set(dstStart+i, src.Get(srcStart+i))
	}

	return nil
}
