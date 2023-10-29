package misc

import "github.com/murakmii/gj/vm"

func init() {
	class := "sun/misc/URLClassPath"

	vm.NativeMethods.Register(class, "getLookupCacheURLs", "(Ljava/lang/ClassLoader;)[Ljava/net/URL;", func(thread *vm.Thread, args []interface{}) error {
		thread.CurrentFrame().PushOperand(nil)
		return nil
	})
}
