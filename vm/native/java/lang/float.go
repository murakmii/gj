package lang

import (
	"fmt"
	"github.com/murakmii/gj/vm"
	"math"
)

func FloatFloatToRawIntBits(thread *vm.Thread, args []interface{}) error {
	float, ok := args[0].(float32)
	if !ok {
		return fmt.Errorf("Float.floatToRawIntBits received not float(%+v)", args[0])
	}

	thread.CurrentFrame().PushOperand(int32(math.Float32bits(float)))
	return nil
}
