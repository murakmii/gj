package vm

import "unicode/utf16"

type GoString string

var (
	javaLangString           = "java/lang/String"
	javaLangStringValueField = "value"
)

func (s GoString) ToJavaString(thread *Thread) *Instance {
	js := NewInstance(thread.VM().JavaLangStringClass())

	u16 := utf16.Encode([]rune(s))
	charArray := NewArray("C", len(u16))
	for i, e := range u16 {
		charArray.Set(i, int(e))
	}

	js.PutField(&javaLangString, &javaLangStringValueField, charArray)
	return js
}
