package native

import (
	"github.com/murakmii/gj/vm"
	"github.com/murakmii/gj/vm/native/java/lang"
)

func init() {
	vm.RegisterNativeMethod("java/lang/Class/desiredAssertionStatus0", lang.ClassDesiredAssertionStatus0)
	vm.RegisterNativeMethod("java/lang/Class/registerNatives", lang.ClassRegisterNatives)
	vm.RegisterNativeMethod("java/lang/Class/getPrimitiveClass", lang.ClassGetPrimitiveClass)

	vm.RegisterNativeMethod("java/lang/Double/doubleToRawLongBits", lang.DoubleDoubleToRawLongBits)
	vm.RegisterNativeMethod("java/lang/Double/longBitsToDouble", lang.DoubleLongBitsToDouble)

	vm.RegisterNativeMethod("java/lang/Float/floatToRawIntBits", lang.FloatFloatToRawIntBits)

	vm.RegisterNativeMethod("java/lang/Object/registerNatives", lang.ObjectRegisterNatives)

	vm.RegisterNativeMethod("java/lang/System/arraycopy", lang.SystemArrayCopy)
	vm.RegisterNativeMethod("java/lang/System/registerNatives", lang.SystemRegisterNatives)
}
