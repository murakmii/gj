package misc

import (
	"fmt"
	"github.com/murakmii/gj/vm"
)

func VMInitialize(_ *vm.Thread, args []interface{}) error {
	fmt.Println("execute sum/misc/VM.initialize")
	return nil
}
