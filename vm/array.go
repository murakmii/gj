package vm

type Array struct {
	elem       []interface{}
	descOfElem string
}

func NewArray(descOfElem string, length int) *Array {
	return &Array{
		elem:       make([]interface{}, length),
		descOfElem: descOfElem,
	}
}

func (array *Array) Length() int {
	return len(array.elem)
}

func (array *Array) Set(i int, value interface{}) {
	array.elem[i] = value
}

func (array *Array) Get(i int) interface{} {
	return array.elem[i]
}
