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
	js := NewInstance(thread.VM().JavaLangStringClass())

	u16 := utf16.Encode([]rune(s))
	charArray := NewArray("C", len(u16))
	for i, e := range u16 {
		charArray.Set(i, int(e))
	}

	js.PutField(&javaLangStringValueField, &javaLangStringValueDesc, charArray)
	return js
}

func ByteSliceToJavaArray(bytes []byte) *Array {
	array := NewArray("B", len(bytes))
	for i, b := range bytes {
		array.Set(i, b)
	}
	return array
}

func JavaByteArrayToGo(array *Array, offset, size int) []byte {
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		bytes[i] = byte(array.Get(i + offset).(int))
	}
	return bytes
}
