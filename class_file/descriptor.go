package class_file

import "strings"

type (
	MethodDescriptor string
	FieldType        string
)

func (m MethodDescriptor) Params() []FieldType {
	var result []FieldType
	arrayDim := 0

	for i := 1; m[i] != ')'; i++ {
		switch m[i] {
		case '[':
			arrayDim++
		case 'L':
			idx := strings.Index(string(m)[i:], ";")
			result = append(result, FieldType(strings.Repeat("[", arrayDim)+string(m)[i:i+idx+1]))
			arrayDim = 0
			i += idx
		default:
			result = append(result, FieldType(strings.Repeat("[", arrayDim)+string(m[i])))
			arrayDim = 0
		}
	}

	return result
}

func (m MethodDescriptor) String() string {
	return string(m)
}

func (f FieldType) Type() string {
	if f[0] == '[' {
		return string(f)
	}

	switch f {
	case "B":
		return "byte"
	case "C":
		return "char"
	case "D":
		return "double"
	case "F":
		return "float"
	case "I":
		return "int"
	case "J":
		return "long"
	case "S":
		return "short"
	case "Z":
		return "boolean"
	}

	l := len(f)
	return string(f)[1 : l-1]
}
