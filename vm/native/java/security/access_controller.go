package security

import (
	"github.com/murakmii/gj/vm"
)

func init() {
	class := "java/security/AccessController"

	vm.NativeMethods.Register(class, "initIDs", "()V", vm.NopNativeMethod)

	doPrivileged := func(thread *vm.Thread, args []interface{}) error {
		action := args[0].(*vm.Instance)
		runClass, runMethod := action.Class().ResolveMethod("run", "()Ljava/lang/Object;")
		thread.CurrentFrame().PushOperand(action)
		return thread.ExecMethod(runClass, runMethod)
	}

	for _, desc := range []string{
		"(Ljava/security/PrivilegedAction;)Ljava/lang/Object;",
		"(Ljava/security/PrivilegedAction;Ljava/security/AccessControlContext;)Ljava/lang/Object;",
		"(Ljava/security/PrivilegedExceptionAction;Ljava/security/AccessControlContext;)Ljava/lang/Object;",
		"(Ljava/security/PrivilegedExceptionAction;)Ljava/lang/Object;",
	} {
		vm.NativeMethods.Register(class, "doPrivileged", desc, doPrivileged)
	}

	vm.NativeMethods.Register(class, "getStackAccessControlContext", "()Ljava/security/AccessControlContext;", func(thread *vm.Thread, args []interface{}) error {
		thread.CurrentFrame().PushOperand(nil)
		return nil
	})
}
