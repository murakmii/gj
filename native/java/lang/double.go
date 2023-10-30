package lang

import (
	"fmt"
	"github.com/murakmii/gojiai/vm"
	"math"
)

func init() {
	class := "java/lang/Double"

	vm.NativeMethods.Register(class, "doubleToRawLongBits", "(D)J", func(thread *vm.Thread, args []interface{}) error {
		float, ok := args[0].(float64)
		if !ok {
			return fmt.Errorf("Double.doubleToRawLongBits received not double(%+v)", args[0])
		}

		thread.CurrentFrame().PushOperand(int64(math.Float64bits(float)))
		return nil
	})

	vm.NativeMethods.Register(class, "longBitsToDouble", "(J)D", func(thread *vm.Thread, args []interface{}) error {
		i64, ok := args[0].(int64)
		if !ok {
			return fmt.Errorf("Double.longBitsToDouble received not int64(%+v)", args[0])
		}

		thread.CurrentFrame().PushOperand(math.Float64frombits(uint64(i64)))
		return nil
	})
}
