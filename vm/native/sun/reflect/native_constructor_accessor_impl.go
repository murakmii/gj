package reflect

import (
	"fmt"
	"github.com/murakmii/gj/vm"
)

func NativeConstructorAccessorImplNewInstance0(thread *vm.Thread, args []interface{}) error {
	cstr := args[0].(*vm.Instance)

	clazzName := "clazz"
	clazzDesc := "Ljava/lang/Class;"
	class, err := thread.VM().FindClass(cstr.GetField(&clazzName, &clazzDesc).(*vm.Instance).VMData().(*string))
	if err != nil {
		return err
	}

	slotName := "slot"
	slotDesc := "I"
	method := class.File().FindMethodByID(cstr.GetField(&slotName, &slotDesc).(int))

	var cstrArgs *vm.Array
	if args[1] == nil {
		cstrArgs = vm.NewArray("Ljava/lang/Object;", 0)
	}

	locals := make([]interface{}, cstrArgs.Length()+1)
	locals[0] = vm.NewInstance(class)
	for i := 0; i < cstrArgs.Length(); i++ {
		locals[i+1] = cstrArgs.Get(i)
	}

	thrown, err := thread.Derive().Execute(vm.NewFrame(class, method).SetLocals(locals))
	if err != nil {
		return err
	}
	if thrown != nil {
		return fmt.Errorf("constructor thrown exception %+v in NativeConstructorAccessorImpl.newInstance0", thrown)
	}

	thread.CurrentFrame().PushOperand(locals[0])
	return nil
}
