package native

import (
	"github.com/murakmii/gj/vm"
	"github.com/murakmii/gj/vm/native/java/lang"
	"github.com/murakmii/gj/vm/native/java/security"
	"github.com/murakmii/gj/vm/native/sun/misc"
)

var nop = func(_ *vm.Thread, _ []interface{}) error {
	return nil
}

func init() {
	vm.RegisterNativeMethod("java/security/AccessController/getStackAccessControlContext()Ljava/security/AccessControlContext;", security.AccessControllerGetStackAccessControlContext)

	vm.RegisterNativeMethod("java/io/FileDescriptor/initIDs()V", nop)

	vm.RegisterNativeMethod("java/io/FileInputStream/initIDs()V", nop)

	vm.RegisterNativeMethod("java/lang/Class/desiredAssertionStatus0(Ljava/lang/Class;)Z", lang.ClassDesiredAssertionStatus0)
	vm.RegisterNativeMethod("java/lang/Class/registerNatives()V", nop)
	vm.RegisterNativeMethod("java/lang/Class/getPrimitiveClass(Ljava/lang/String;)Ljava/lang/Class;", lang.ClassGetPrimitiveClass)

	vm.RegisterNativeMethod("java/lang/Double/doubleToRawLongBits(D)J", lang.DoubleDoubleToRawLongBits)
	vm.RegisterNativeMethod("java/lang/Double/longBitsToDouble(J)D", lang.DoubleLongBitsToDouble)

	vm.RegisterNativeMethod("java/lang/Float/floatToRawIntBits(F)I", lang.FloatFloatToRawIntBits)

	vm.RegisterNativeMethod("java/lang/Object/registerNatives()V", nop)
	vm.RegisterNativeMethod("java/lang/Object/hashCode()I", lang.ObjectHashCode)

	vm.RegisterNativeMethod("java/lang/System/arraycopy", lang.SystemArrayCopy)
	vm.RegisterNativeMethod("java/lang/System/initProperties(Ljava/util/Properties;)Ljava/util/Properties;", lang.SystemInitProperties)
	vm.RegisterNativeMethod("java/lang/System/registerNatives()V", nop)

	vm.RegisterNativeMethod("java/lang/Thread/currentThread()Ljava/lang/Thread;", lang.ThreadCurrentThread)
	vm.RegisterNativeMethod("java/lang/Thread/registerNatives()V", nop)
	vm.RegisterNativeMethod("java/lang/Thread/setPriority0(I)V", nop)

	vm.RegisterNativeMethod("sun/misc/Unsafe/addressSize()I", misc.UnsafeAddressSize)
	vm.RegisterNativeMethod("sun/misc/Unsafe/arrayBaseOffset(Ljava/lang/Class;)I", misc.UnsafeArrayBaseOffset)
	vm.RegisterNativeMethod("sun/misc/Unsafe/arrayIndexScale(Ljava/lang/Class;)I", misc.UnsafeArrayIndexScale)
	vm.RegisterNativeMethod("sun/misc/Unsafe/registerNatives()V", nop)

	vm.RegisterNativeMethod("sun/misc/VM/initialize()V", nop)
}
