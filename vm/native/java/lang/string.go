package lang

import (
	"github.com/murakmii/gj/vm"
	"strings"
	"unicode/utf16"
)

func StringIntern(thread *vm.Thread, args []interface{}) error {
	js := args[0].(*vm.Instance)

	valueName := "value"
	valueDesc := "[C"
	valueArray := js.GetField(&valueName, &valueDesc).(*vm.Array)

	u16 := make([]uint16, valueArray.Length())
	for i := 0; i < valueArray.Length(); i++ {
		u16[i] = uint16((valueArray.Get(i)).(int))
	}
	gs := strings.ReplaceAll(string(utf16.Decode(u16)), ".", "/")

	interned, err := thread.VM().JavaString(thread, &gs)
	if err != nil {
		return err
	}

	thread.CurrentFrame().PushOperand(interned)
	return nil
}
