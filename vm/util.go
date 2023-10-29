package vm

func ByteSliceToJavaArray(vm *VM, bytes []byte) *Instance {
	instance, slice := NewArray(vm, "[B", len(bytes))
	for i, b := range bytes {
		slice[i] = int32(b)
	}
	return instance
}

func JavaByteArrayToGo(array *Instance, offset, size int) []byte {
	slice := array.AsArray()
	bytes := make([]byte, size)

	for i := 0; i < size; i++ {
		bytes[i] = byte(slice[offset+i].(int32))
	}
	return bytes
}
