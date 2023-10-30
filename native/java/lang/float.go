package lang

import (
	"fmt"
	"github.com/murakmii/gojiai/vm"
	"math"
)

func init() {
	class := "java/lang/Float"

	vm.NativeMethods.Register(class, "floatToRawIntBits", "(F)I", func(thread *vm.Thread, args []interface{}) error {
		float, ok := args[0].(float32)
		if !ok {
			return fmt.Errorf("Float.floatToRawIntBits received not float(%+v)", args[0])
		}

		thread.CurrentFrame().PushOperand(int32(math.Float32bits(float)))
		return nil
	})
}
