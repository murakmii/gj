package vm

import (
	"unicode/utf16"
)

type GoString string

var (
	javaLangStringValueField = "value"
	javaLangStringValueDesc  = "[C"
)

func (s GoString) ToJavaString(thread *Thread) *Instance {
	js := NewInstance(thread.VM().StdClass(JavaLangString))

	u16 := utf16.Encode([]rune(s))
	instance, slice := NewArray(thread.VM(), "[C", len(u16))

	for i, e := range u16 {
		slice[i] = int(e)
	}

	js.PutField(&javaLangStringValueField, &javaLangStringValueDesc, instance)
	return js
}

func ByteSliceToJavaArray(vm *VM, bytes []byte) *Instance {
	instance, slice := NewArray(vm, "[B", len(bytes))
	for i, b := range bytes {
		slice[i] = int(b)
	}
	return instance
}

func JavaByteArrayToGo(array *Instance, offset, size int) []byte {
	slice := array.AsArray()
	bytes := make([]byte, size)

	for i := 0; i < size; i++ {
		bytes[i] = byte(slice[offset+i].(int))
	}
	return bytes
}
