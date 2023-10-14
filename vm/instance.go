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
