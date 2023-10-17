package vm

type Instance struct {
	class  *Class
	fields map[string]interface{}
}

func NewInstance(class *Class) *Instance {
	return &Instance{
		class:  class,
		fields: make(map[string]interface{}),
	}
}

func (instance *Instance) Class() *Class {
	return instance.class
}

func (instance *Instance) GetField(class, name *string) interface{} {
	return instance.fields[*class+"."+*name]
}

func (instance *Instance) PutField(class, name *string, value interface{}) {
	instance.fields[*class+"."+*name] = value
}
