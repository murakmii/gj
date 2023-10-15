package lang

import (
	"fmt"
	"github.com/murakmii/gj/vm"
)

func ObjectRegisterNatives(thread *vm.Thread, args []interface{}) error {
	fmt.Println("execute java/lang/Object.registerNatives")
	return nil
}
