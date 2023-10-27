package reflect

import (
	"github.com/murakmii/gj/vm"
)

func NativeConstructorAccessorImplNewInstance0(thread *vm.Thread, args []interface{}) error {
	cstr := args[0].(*vm.Instance)

	clazzName := "clazz"
	clazzDesc := "Ljava/lang/Class;"
	class, err := thread.VM().Class(*(cstr.GetField(&clazzName, &clazzDesc).(*vm.Instance).VMData().(*string)), thread)
	if err != nil {
		return err
	}

	slotName := "slot"
	slotDesc := "I"
	method := class.File().FindMethodByID(int(cstr.GetField(&slotName, &slotDesc).(int32)))

	var cstrArgs []interface{}
	if args[1] != nil {
		cstrArgs = args[1].(*vm.Instance).AsArray()
	}

	locals := make([]interface{}, len(cstrArgs)+1)
	locals[0] = vm.NewInstance(class)
	for i, a := range cstrArgs {
		locals[i+1] = a
	}

	err = thread.Execute(vm.NewFrame(class, method).SetLocals(locals))
	if err != nil {
		return err
	}

	thread.CurrentFrame().PushOperand(locals[0])
	return nil
}
