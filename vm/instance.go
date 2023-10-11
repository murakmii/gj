package vm

type Instance struct {
	class  *Class
	fields map[string]interface{}
	super  *Instance
}

func (instance *Instance) IsSubClassOf(className *string) bool {
	return instance.class.File().ThisClass() == *className || instance.super.IsSubClassOf(className)
}
