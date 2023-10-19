package lang

import (
	"fmt"
	"github.com/murakmii/gj/vm"
	"strings"
	"unicode/utf16"
)

func ClassDesiredAssertionStatus0(thread *vm.Thread, args []interface{}) error {
	// return false
	thread.CurrentFrame().PushOperand(0)
	return nil
}

func ClassGetPrimitiveClass(thread *vm.Thread, _ []interface{}) error {
	// TODO: generate class instance
	// return null
	thread.CurrentFrame().PushOperand(nil)
	return nil
}

func ClassGetName0(thread *vm.Thread, args []interface{}) error {
	switch javaClass := args[0].(type) {
	case *vm.Instance:
		name := strings.ReplaceAll(*javaClass.VMData().(*string), "/", ".")
		nameJS, err := thread.VM().JavaString(thread, &name)
		if err != nil {
			return err
		}
		thread.CurrentFrame().PushOperand(nameJS)

	default:
		return fmt.Errorf("Clas.getName0 does NOT support object is NOT instance")
	}

	return nil
}

func ClassForName0(thread *vm.Thread, args []interface{}) error {
	nameJS := args[0].(*vm.Instance)

	valueName := "value"
	valueDesc := "[C"
	valueArray := nameJS.GetField(&valueName, &valueDesc).(*vm.Array)

	u16 := make([]uint16, valueArray.Length())
	for i := 0; i < valueArray.Length(); i++ {
		u16[i] = uint16((valueArray.Get(i)).(int))
	}
	nameGS := strings.ReplaceAll(string(utf16.Decode(u16)), ".", "/")

	className := "java/lang/Class"
	class, state, err := thread.VM().FindInitializedClass(&className, thread)
	if err != nil {
		return err
	}
	if state == vm.FailedInitialization {
		return fmt.Errorf("failed initialization of class class in LDC")
	}

	thread.CurrentFrame().PushOperand(vm.NewInstance(class).SetVMData(&nameGS))
	return nil
}

func ClassGetDeclaredFields0(thread *vm.Thread, args []interface{}) error {
	thread.DumpFrameStack(true)
	return fmt.Errorf("Class.getDeclaredFields0")
}
