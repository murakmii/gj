package vm

import (
	"fmt"
	"github.com/murakmii/gojiai/class_file"
	"os"
	"strings"
	"unicode/utf16"
	"unsafe"
)

type (
	Instance struct {
		class   *Class
		fields  []interface{}
		monitor *Monitor

		// Any data for VM implementation. e.g.,
		// * *vm.Class for instance of java.lang.Class
		// * *vm.Thread for instance of java.lang.Thread
		vmData interface{}
	}
)

func NewInstance(class *Class) *Instance {
	return &Instance{
		class:   class,
		fields:  make([]interface{}, class.TotalInstanceFields()),
		monitor: NewMonitor(),
	}
}

func NewArray(vm *VM, desc string, size int) (*Instance, []interface{}) {
	arrayClass, _ := vm.Class(desc, nil)
	array := make([]interface{}, size)

	compType := desc[strings.LastIndex(desc, "[")+1:]
	defaultVal := class_file.JavaTypeSignature(compType).DefaultValue()

	if defaultVal != nil {
		for i := range array {
			array[i] = defaultVal
		}
	}

	return &Instance{
		class:   arrayClass,
		fields:  array, // Array has elements in fields
		monitor: NewMonitor(),
	}, array
}

func NewString(vm *VM, str string) *Instance {
	javaStr := NewInstance(vm.SpecialClass(JavaLangStringID))

	u16 := utf16.Encode([]rune(str))
	instance, slice := NewArray(vm, "[C", len(u16))

	for i, e := range u16 {
		slice[i] = int32(e)
	}

	javaStr.PutField("value", "[C", instance)
	return javaStr
}

func (instance *Instance) Class() *Class {
	return instance.class
}

func (instance *Instance) CompareAndSwapInt(id int, expected, x int32) (bool, error) {
	// TODO: lock
	if instance.fields[id] == nil {
		instance.fields[id] = int32(0)
	}

	target, ok := instance.fields[id].(int32)
	if !ok {
		return false, fmt.Errorf("Instance.CompareAndSwapInt only suuport int value")
	}

	if target == expected {
		instance.fields[id] = x
		return true, nil
	}

	return false, nil
}

func (instance *Instance) CompareAndSwapLong(id int, expected, x int64) (bool, error) {
	// TODO: lock
	if instance.fields[id] == nil {
		instance.fields[id] = int64(0)
	}

	target, ok := instance.fields[id].(int64)
	if !ok {
		return false, fmt.Errorf("Instance.CompareAndSwapLong only suuport int64 value")
	}

	if target == expected {
		instance.fields[id] = x
		return true, nil
	}

	return false, nil
}

func (instance *Instance) CompareAndSwap(id int, expected, x *Instance) (bool, error) {
	// TODO: lock
	// TODO: default value check
	if instance.fields[id] == nil {
		if expected != nil {
			return false, nil
		}
		instance.fields[id] = x
		return true, nil
	}

	target, ok := instance.fields[id].(*Instance)
	if !ok {
		return false, fmt.Errorf("Instance.CompareAndSwap only suuport instance value")
	}

	if target == expected {
		instance.fields[id] = x
		return true, nil
	}

	return false, nil
}

func (instance *Instance) GetField(name, desc string) interface{} {
	_, field := instance.class.ResolveField(name, desc)

	value := instance.fields[field.ID()]
	if value == nil && !field.NullableDefaultValue() {
		instance.fields[field.ID()] = field.DefaultValue()
		value = instance.fields[field.ID()]
	}

	return value
}

func (instance *Instance) PutField(name, desc string, value interface{}) {
	_, field := instance.class.ResolveField(name, desc)
	instance.fields[field.ID()] = value
}

func (instance *Instance) GetFieldByID(id int) interface{} {
	return instance.fields[id]
}

func (instance *Instance) PutFieldByID(id int, value interface{}) {
	instance.fields[id] = value
}

func (instance *Instance) Monitor() *Monitor {
	return instance.monitor
}

func (instance *Instance) AsArray() []interface{} {
	return instance.fields
}

// For instance of java.lang.Class
func (instance *Instance) AsClass() *Class        { return instance.vmData.(*Class) }
func (instance *Instance) ToBeClass(class *Class) { instance.vmData = class }

// For instance of java.lang.Thread
func (instance *Instance) ToBeThread(thread *Thread) { instance.vmData = thread }
func (instance *Instance) AsThread() *Thread {
	if instance.vmData == nil { // vmData of java.lang.Thread instance is nil until calling start0
		return nil
	}
	return instance.vmData.(*Thread)
}

// For instance of java.lang.Throwable
func (instance *Instance) ToBeThrowable(traces []*StackTraceElement) { instance.vmData = traces }
func (instance *Instance) AsThrowable() []*StackTraceElement {
	// vmData of java.lang.Throwable instance is nil if fillInStackTrace is skipped
	// https://github.com/openjdk/jdk8u/blob/master/jdk/src/share/classes/java/lang/Throwable.java#L360
	if instance.vmData == nil {
		return nil
	}
	return instance.vmData.([]*StackTraceElement)
}

// For instance of java.lang.String
func (instance *Instance) AsString() string {
	// java.lang.String has value field contains string content.
	// https://github.com/openjdk/jdk8u/blob/master/jdk/src/share/classes/java/lang/String.java#L114
	slice := instance.GetField("value", "[C").(*Instance).AsArray()

	u16 := make([]uint16, len(slice))
	for i := range u16 {
		u16[i] = uint16(slice[i].(int32))
	}

	return string(utf16.Decode(u16))
}

// For instance of java.io.FileDescriptor
func (instance *Instance) AsFile() *os.File {
	if instance.vmData != nil {
		return instance.vmData.(*os.File)
	}

	var file *os.File

	// Constructor of FileDescriptor could be received fd already opened(>= 0). e.g., stdout, stdin, stderr
	// In this case, open it by os.NewFile
	// https://github.com/openjdk/jdk8u/blob/master/jdk/src/solaris/classes/java/io/FileDescriptor.java#L62
	fd := instance.GetField("fd", "I").(int32)
	if fd >= 0 {
		file = os.NewFile(uintptr(fd), "")
	}

	instance.vmData = file
	return file
}

func (instance *Instance) ToBeFile(file *os.File) {
	instance.PutField("fd", "I", int32(file.Fd()))
	instance.vmData = file
}

func (instance *Instance) HashCode() int32 {
	return int32(uintptr(unsafe.Pointer(instance)))
}

func (instance *Instance) Clone() *Instance {
	fields := make([]interface{}, len(instance.fields))
	copy(fields, instance.fields)

	return &Instance{
		class:   instance.class,
		fields:  fields,
		monitor: NewMonitor(),
		vmData:  instance.vmData,
	}
}
