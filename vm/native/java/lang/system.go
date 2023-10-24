package lang

import (
	"fmt"
	"github.com/murakmii/gj/vm"
	"time"
)

func SystemArrayCopy(thread *vm.Thread, args []interface{}) error {
	src := args[0].(*vm.Instance).AsArray()
	srcStart := args[1].(int)
	dst := args[2].(*vm.Instance).AsArray()
	dstStart := args[3].(int)
	count := args[4].(int)

	// TODO: copy
	for i := 0; i < count; i++ {
		dst[dstStart+i] = src[srcStart+i]
	}

	return nil
}

func SystemCurrentTimeMillis(thread *vm.Thread, _ []interface{}) error {
	thread.CurrentFrame().PushOperand(time.Now().UnixMilli())
	return nil
}

func SystemInitProperties(thread *vm.Thread, args []interface{}) error {
	props, ok := args[0].(*vm.Instance)
	if !ok {
		return fmt.Errorf("argument of System.initProperties is NOT class instance")
	}

	class, method := props.Class().ResolveMethod("setProperty", "(Ljava/lang/String;Ljava/lang/String;)Ljava/lang/Object;")

	for k, v := range thread.VM().SysProps() {
		kJS, err := thread.VM().JavaString(thread, &k)
		if err != nil {
			return fmt.Errorf("failed to instantiate string for system prorperty key")
		}

		vJS, err := thread.VM().JavaString(thread, &v)
		if err != nil {
			return fmt.Errorf("failed to instantiate string for system prorperty value")
		}

		thrown, err := thread.Derive().Execute(vm.NewFrame(class, method).SetLocals([]interface{}{props, kJS, vJS}))
		if err != nil {
			return err
		}
		if thrown != nil {
			return fmt.Errorf("failed to set system property %s:%s = %+v", k, v, thrown)
		}
	}

	thread.CurrentFrame().PushOperand(props)
	return nil
}

func SystemSetArg0ToField(name, desc string) vm.NativeMethodFunc {
	return func(thread *vm.Thread, args []interface{}) error {
		sys, err := thread.VM().Class("java/lang/System", thread)
		if err != nil {
			return err
		}

		_, field := sys.ResolveField(name, desc)
		sys.SetStaticField(field, args[0])

		return nil
	}
}
