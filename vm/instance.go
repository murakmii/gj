package vm

import "fmt"

type Instance struct {
	class   *Class
	fields  map[string]interface{}
	monitor *Monitor
}

func NewInstance(class *Class) *Instance {
	return &Instance{
		class:   class,
		fields:  make(map[string]interface{}),
		monitor: NewMonitor(),
	}
}

func (instance *Instance) Class() *Class {
	return instance.class
}

func (instance *Instance) Cast(class *Class) {
	if class.File().ThisClass() == "java/util/Collection" {
		fmt.Printf("--------------------- interface cast!!!!! %s to %s---------------\n", instance.Class().File().ThisClass(), class.File().ThisClass())
	}
	instance.class = class
}

func (instance *Instance) GetField(name, desc *string) interface{} {
	class, field := instance.class.ResolveField(*name, *desc)
	fName := class.File().ThisClass() + "." + *field.Name()

	value, exist := instance.fields[fName]
	if !exist {
		instance.fields[fName] = field.DefaultValue()
		value = instance.fields[fName]
	}

	return value
}

func (instance *Instance) PutField(class, name *string, value interface{}) {
	instance.fields[*class+"."+*name] = value
}

func (instance *Instance) Monitor() *Monitor {
	return instance.monitor
}
