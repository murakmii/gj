package vm

type Array struct {
	elements []interface{}
}

func NewArray(lenEachDim []int) *Array {
	array := &Array{elements: make([]interface{}, lenEachDim[0])}
	if len(lenEachDim) > 1 {
		for i, _ := range array.elements {
			array.elements[i] = NewArray(lenEachDim[1:])
		}
	}

	return array
}

func (array *Array) Length() int {
	return len(array.elements)
}

func (array *Array) Set(i int, value interface{}) {
	array.elements[i] = value
}

func (array *Array) Get(i int) interface{} {
	return array.elements[i]
}
