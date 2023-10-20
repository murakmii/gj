package native

import (
	"github.com/murakmii/gj/vm"
	"github.com/murakmii/gj/vm/native/java/lang"
	"github.com/murakmii/gj/vm/native/java/security"
	"github.com/murakmii/gj/vm/native/java/util/concurrent/atomic"
	"github.com/murakmii/gj/vm/native/sun/misc"
	"github.com/murakmii/gj/vm/native/sun/reflect"
)

var nop = func(_ *vm.Thread, _ []interface{}) error {
	return nil
}

func init() {
	vm.RegisterNativeMethod("java/security/AccessController/getStackAccessControlContext()Ljava/security/AccessControlContext;", security.AccessControllerGetStackAccessControlContext)
	vm.RegisterNativeMethod("java/security/AccessController/doPrivileged(Ljava/security/PrivilegedAction;)Ljava/lang/Object;", security.AccessControllerDoPrivilegedL)
	vm.RegisterNativeMethod("java/security/AccessController/doPrivileged(Ljava/security/PrivilegedExceptionAction;)Ljava/lang/Object;", security.AccessControllerDoPrivilegedL)

	vm.RegisterNativeMethod("java/io/FileDescriptor/initIDs()V", nop)

	vm.RegisterNativeMethod("java/io/FileInputStream/initIDs()V", nop)

	vm.RegisterNativeMethod("java/io/FileOutputStream/initIDs()V", nop)

	vm.RegisterNativeMethod("java/io/UnixFileSystem/initIDs()V", nop)

	vm.RegisterNativeMethod("java/lang/Class/desiredAssertionStatus0(Ljava/lang/Class;)Z", lang.ClassDesiredAssertionStatus0)
	vm.RegisterNativeMethod("java/lang/Class/getDeclaredConstructors0(Z)[Ljava/lang/reflect/Constructor;", lang.ClassGetDeclaredConstructors)
	vm.RegisterNativeMethod("java/lang/Class/getDeclaredFields0(Z)[Ljava/lang/reflect/Field;", lang.ClassGetDeclaredFields0)
	vm.RegisterNativeMethod("java/lang/Class/getModifiers()I", lang.ClassGetModifiers)
	vm.RegisterNativeMethod("java/lang/Class/getName0()Ljava/lang/String;", lang.ClassGetName0)
	vm.RegisterNativeMethod("java/lang/Class/getSuperclass()Ljava/lang/Class;", lang.ClassGetSuperClass)
	vm.RegisterNativeMethod("java/lang/Class/isAssignableFrom(Ljava/lang/Class;)Z", lang.ClassIsAssignableFrom)
	vm.RegisterNativeMethod("java/lang/Class/isInterface()Z", lang.ClassIsInterface)
	vm.RegisterNativeMethod("java/lang/Class/isPrimitive()Z", lang.ClassIsPrimitive)
	vm.RegisterNativeMethod("java/lang/Class/forName0(Ljava/lang/String;ZLjava/lang/ClassLoader;Ljava/lang/Class;)Ljava/lang/Class;", lang.ClassForName0)
	vm.RegisterNativeMethod("java/lang/Class/registerNatives()V", nop)
	vm.RegisterNativeMethod("java/lang/Class/getPrimitiveClass(Ljava/lang/String;)Ljava/lang/Class;", lang.ClassGetPrimitiveClass)

	vm.RegisterNativeMethod("java/lang/ClassLoader/registerNatives()V", nop)

	vm.RegisterNativeMethod("java/lang/Double/doubleToRawLongBits(D)J", lang.DoubleDoubleToRawLongBits)
	vm.RegisterNativeMethod("java/lang/Double/longBitsToDouble(J)D", lang.DoubleLongBitsToDouble)

	vm.RegisterNativeMethod("java/lang/Float/floatToRawIntBits(F)I", lang.FloatFloatToRawIntBits)

	vm.RegisterNativeMethod("java/lang/Object/getClass()Ljava/lang/Class;", lang.ObjectGetClass)
	vm.RegisterNativeMethod("java/lang/Object/registerNatives()V", nop)
	vm.RegisterNativeMethod("java/lang/Object/hashCode()I", lang.ObjectHashCode)

	vm.RegisterNativeMethod("java/lang/System/arraycopy(Ljava/lang/Object;ILjava/lang/Object;II)V", lang.SystemArrayCopy)
	vm.RegisterNativeMethod("java/lang/System/initProperties(Ljava/util/Properties;)Ljava/util/Properties;", lang.SystemInitProperties)
	vm.RegisterNativeMethod("java/lang/System/registerNatives()V", nop)
	vm.RegisterNativeMethod("java/lang/System/setIn0(Ljava/io/InputStream;)V", lang.SystemSetArg0ToField("in", "Ljava/io/InputStream;"))
	vm.RegisterNativeMethod("java/lang/System/setOut0(Ljava/io/PrintStream;)V", lang.SystemSetArg0ToField("out", "Ljava/io/PrintStream;"))
	vm.RegisterNativeMethod("java/lang/System/setErr0(Ljava/io/PrintStream;)V", lang.SystemSetArg0ToField("err", "Ljava/io/PrintStream;"))

	vm.RegisterNativeMethod("java/lang/String/intern()Ljava/lang/String;", lang.StringIntern)

	vm.RegisterNativeMethod("java/lang/Thread/currentThread()Ljava/lang/Thread;", lang.ThreadCurrentThread)
	vm.RegisterNativeMethod("java/lang/Thread/isAlive()Z", lang.ThreadIsAlive)
	vm.RegisterNativeMethod("java/lang/Thread/start0()V", lang.ThreadStart0)
	vm.RegisterNativeMethod("java/lang/Thread/registerNatives()V", nop)
	vm.RegisterNativeMethod("java/lang/Thread/setPriority0(I)V", nop)

	vm.RegisterNativeMethod("java/lang/Throwable/fillInStackTrace(I)Ljava/lang/Throwable;", lang.ThrowableFillInStackTrace)

	vm.RegisterNativeMethod("java/util/concurrent/atomic/AtomicLong/VMSupportsCS8()Z", atomic.AtomicLongVMSupportsCS8)

	vm.RegisterNativeMethod("sun/misc/Unsafe/allocateMemory(J)J", misc.UnsafeAllocateMemory)
	vm.RegisterNativeMethod("sun/misc/Unsafe/freeMemory(J)V", misc.UnsafeFreeMemory)
	vm.RegisterNativeMethod("sun/misc/Unsafe/getByte(J)B", misc.UnsafeGetByteNativeMem)
	vm.RegisterNativeMethod("sun/misc/Unsafe/putLong(JJ)V", misc.UnsafePutLongNativeMem)
	vm.RegisterNativeMethod("sun/misc/Unsafe/addressSize()I", misc.UnsafeAddressSize)
	vm.RegisterNativeMethod("sun/misc/Unsafe/arrayBaseOffset(Ljava/lang/Class;)I", misc.UnsafeArrayBaseOffset)
	vm.RegisterNativeMethod("sun/misc/Unsafe/arrayIndexScale(Ljava/lang/Class;)I", misc.UnsafeArrayIndexScale)
	vm.RegisterNativeMethod("sun/misc/Unsafe/compareAndSwapInt(Ljava/lang/Object;JII)Z", misc.UnsafeCompareAndSwapInt)
	vm.RegisterNativeMethod("sun/misc/Unsafe/compareAndSwapObject(Ljava/lang/Object;JLjava/lang/Object;Ljava/lang/Object;)Z", misc.UnsafeCompareAndSwapObject)
	vm.RegisterNativeMethod("sun/misc/Unsafe/getIntVolatile(Ljava/lang/Object;J)I", misc.UnsafeGetIntVolatile)
	vm.RegisterNativeMethod("sun/misc/Unsafe/objectFieldOffset(Ljava/lang/reflect/Field;)J", misc.UnsafeObjectFieldOffset)
	vm.RegisterNativeMethod("sun/misc/Unsafe/registerNatives()V", nop)

	vm.RegisterNativeMethod("sun/misc/Signal/findSignal(Ljava/lang/String;)I", misc.SignalFindSignal)
	vm.RegisterNativeMethod("sun/misc/Signal/handle0(IJ)J", misc.SignalHandle)

	vm.RegisterNativeMethod("sun/reflect/NativeConstructorAccessorImpl/newInstance0(Ljava/lang/reflect/Constructor;[Ljava/lang/Object;)Ljava/lang/Object;", reflect.NativeConstructorAccessorImplNewInstance0)

	vm.RegisterNativeMethod("sun/reflect/Reflection/getCallerClass()Ljava/lang/Class;", reflect.ReflectionGetCallerClassV)
	vm.RegisterNativeMethod("sun/reflect/Reflection/getClassAccessFlags(Ljava/lang/Class;)I", reflect.ReflectionGetClassAccessFlags)

	vm.RegisterNativeMethod("sun/misc/VM/initialize()V", nop)
}
