package vm

type Thread struct {
	vm   *VM
	java *Instance
}

func NewThread(vm *VM, java *Instance) *Thread {
	return &Thread{vm: vm, java: java}
}

func (thread *Thread) SetJavaThread(java *Instance) {
	thread.java = java
}

func (thread *Thread) VM() *VM {
	return thread.vm
}
