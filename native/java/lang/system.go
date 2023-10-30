package lang

import (
	"fmt"
	"github.com/murakmii/gojiai/vm"
	"time"
)

func init() {
	_class := "java/lang/System"

	vm.NativeMethods.Register(_class, "arraycopy", "(Ljava/lang/Object;ILjava/lang/Object;II)V", func(thread *vm.Thread, args []interface{}) error {
		src := args[0].(*vm.Instance).AsArray()
		srcStart := args[1].(int32)
		dst := args[2].(*vm.Instance).AsArray()
		dstStart := args[3].(int32)
		count := args[4].(int32)

		// TODO: copy
		for i := int32(0); i < count; i++ {
			dst[dstStart+i] = src[srcStart+i]
		}

		return nil
	})

	vm.NativeMethods.Register(_class, "currentTimeMillis", "()J", func(thread *vm.Thread, args []interface{}) error {
		thread.CurrentFrame().PushOperand(time.Now().UnixMilli())
		return nil
	})

	vm.NativeMethods.Register(_class, "identityHashCode", "(Ljava/lang/Object;)I", func(thread *vm.Thread, args []interface{}) error {
		thread.CurrentFrame().PushOperand(args[0].(*vm.Instance).HashCode())
		return nil
	})

	vm.NativeMethods.Register(_class, "initProperties", "(Ljava/util/Properties;)Ljava/util/Properties;", func(thread *vm.Thread, args []interface{}) error {
		props, ok := args[0].(*vm.Instance)
		if !ok {
			return fmt.Errorf("argument of System.initProperties is NOT class instance")
		}

		class, method := props.Class().ResolveMethod("setProperty", "(Ljava/lang/String;Ljava/lang/String;)Ljava/lang/Object;")

		for k, v := range thread.VM().SysProps() {
			if err := thread.Execute(vm.NewFrame(class, method).SetLocals([]interface{}{
				props,
				thread.VM().JavaString(k),
				thread.VM().JavaString(v),
			})); err != nil {
				return err
			}
		}

		thread.CurrentFrame().PushOperand(props)
		return nil
	})

	vm.NativeMethods.Register(_class, "nanoTime", "()J", func(thread *vm.Thread, args []interface{}) error {
		thread.CurrentFrame().PushOperand(time.Now().UnixNano())
		return nil
	})

	vm.NativeMethods.Register(_class, "registerNatives", "()V", vm.NopNativeMethod)

	// For setIn0, setOut0, setErr0
	streamSetter := func(name, desc string) vm.NativeMethodFunc {
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

	vm.NativeMethods.Register(_class, "setErr0", "(Ljava/io/PrintStream;)V", streamSetter("err", "Ljava/io/PrintStream;"))
	vm.NativeMethods.Register(_class, "setIn0", "(Ljava/io/InputStream;)V", streamSetter("in", "Ljava/io/InputStream;"))
	vm.NativeMethods.Register(_class, "setOut0", "(Ljava/io/PrintStream;)V", streamSetter("out", "Ljava/io/PrintStream;"))
}
