package vm

type Thread struct {
	vm   *VM
	java *Instance
}

func (thread *Thread) SetJavaThread(java *Instance) {
	thread.java = java
}

func (thread *Thread) VM() *VM {
	return thread.vm
}
