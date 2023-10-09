package vm

type Instance struct {
	class  *Class
	fields map[string]interface{}
	super  *Instance
}
