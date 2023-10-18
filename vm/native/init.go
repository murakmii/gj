package native

import (
	"github.com/murakmii/gj/vm"
	"github.com/murakmii/gj/vm/native/java/io"
	"github.com/murakmii/gj/vm/native/java/lang"
	"github.com/murakmii/gj/vm/native/java/security"
	"github.com/murakmii/gj/vm/native/sun/misc"
)

func init() {
	vm.RegisterNativeMethod("java/security/AccessController/getStackAccessControlContext", security.AccessControllerGetStackAccessControlContext)

	vm.RegisterNativeMethod("java/io/FileDescriptor/initIDs", io.FileDescriptorInitIDs)

	vm.RegisterNativeMethod("java/io/FileInputStream/initIDs", io.FileInputStreamInitIDs)

	vm.RegisterNativeMethod("java/lang/Class/desiredAssertionStatus0", lang.ClassDesiredAssertionStatus0)
	vm.RegisterNativeMethod("java/lang/Class/registerNatives", lang.ClassRegisterNatives)
	vm.RegisterNativeMethod("java/lang/Class/getPrimitiveClass", lang.ClassGetPrimitiveClass)

	vm.RegisterNativeMethod("java/lang/Double/doubleToRawLongBits", lang.DoubleDoubleToRawLongBits)
	vm.RegisterNativeMethod("java/lang/Double/longBitsToDouble", lang.DoubleLongBitsToDouble)

	vm.RegisterNativeMethod("java/lang/Float/floatToRawIntBits", lang.FloatFloatToRawIntBits)

	vm.RegisterNativeMethod("java/lang/Object/registerNatives", lang.ObjectRegisterNatives)
	vm.RegisterNativeMethod("java/lang/Object/hashCode", lang.ObjectHashCode)

	vm.RegisterNativeMethod("java/lang/System/arraycopy", lang.SystemArrayCopy)
	vm.RegisterNativeMethod("java/lang/System/initProperties", lang.SystemInitProperties)
	vm.RegisterNativeMethod("java/lang/System/registerNatives", lang.SystemRegisterNatives)

	vm.RegisterNativeMethod("java/lang/Thread/currentThread", lang.ThreadCurrentThread)
	vm.RegisterNativeMethod("java/lang/Thread/registerNatives", lang.ThreadRegisterNatives)
	vm.RegisterNativeMethod("java/lang/Thread/setPriority0", lang.ThreadSetPriority0)

	vm.RegisterNativeMethod("sun/misc/Unsafe/addressSize", misc.UnsafeAddressSize)
	vm.RegisterNativeMethod("sun/misc/Unsafe/arrayBaseOffset", misc.UnsafeArrayBaseOffset)
	vm.RegisterNativeMethod("sun/misc/Unsafe/arrayIndexScale", misc.UnsafeArrayIndexScale)
	vm.RegisterNativeMethod("sun/misc/Unsafe/registerNatives", misc.UnsafeRegisterNatives)

	vm.RegisterNativeMethod("sun/misc/VM/initialize", misc.VMInitialize)
}
