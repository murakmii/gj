package vm

import (
	"fmt"
	"unicode/utf16"
)

type Instance struct {
	class   *Class
	fields  []interface{}
	monitor *Monitor

	// Any data for VM implementation. e.g.,
	// * Class name for instance of java.lang.Class
	vmData interface{}
}

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

	return NewInstance(arrayClass).SetVMData(array), array
}

func (instance *Instance) Class() *Class {
	return instance.class
}

func (instance *Instance) CompareAndSwapInt(id int, expected, x int) (bool, error) {
	// TODO: lock
	if instance.fields[id] == nil {
		instance.fields[id] = 0
	}

	target, ok := instance.fields[id].(int)
	if !ok {
		return false, fmt.Errorf("Instance.CompareAndSwapInt only suuport int value")
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

func (instance *Instance) GetField(name, desc *string) interface{} {
	_, field := instance.class.ResolveField(*name, *desc)

	value := instance.fields[field.ID()]
	if value == nil && !field.NullableDefaultValue() {
		instance.fields[field.ID()] = field.DefaultValue()
		value = instance.fields[field.ID()]
	}

	return value
}

func (instance *Instance) GetFieldByID(id int) interface{} {
	return instance.fields[id]
}

func (instance *Instance) PutField(name, desc *string, value interface{}) {
	_, field := instance.class.ResolveField(*name, *desc)
	instance.fields[field.ID()] = value
}

func (instance *Instance) Monitor() *Monitor {
	return instance.monitor
}

func (instance *Instance) VMData() interface{} {
	return instance.vmData
}

func (instance *Instance) SetVMData(data interface{}) *Instance {
	instance.vmData = data
	return instance
}

func (instance *Instance) AsArray() []interface{} {
	return instance.vmData.([]interface{})
}

// Utility method to get value of char array field as string.
// e.g., 'value' field of java.lang.String, 'name' field of java.lang.Thread.
func (instance *Instance) GetCharArrayField(name string) string {
	desc := "[C"
	slice := instance.GetField(&name, &desc).(*Instance).AsArray()

	u16 := make([]uint16, len(slice))
	for i := 0; i < len(slice); i++ {
		u16[i] = uint16(slice[i].(int))
	}

	return string(utf16.Decode(u16))
}

func (instance *Instance) Clone() *Instance {
	fields := make([]interface{}, len(instance.fields))
	copy(fields, instance.fields)

	clone := &Instance{
		class:   instance.class,
		fields:  fields,
		monitor: NewMonitor(),
		vmData:  instance.vmData,
	}

	if instance.class.IsArray() {
		srcArray := instance.AsArray()
		dstArray := make([]interface{}, len(srcArray))
		copy(dstArray, srcArray)

		clone.vmData = dstArray
	}

	return clone
}
