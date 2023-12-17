package main

import (
	"fmt"
	"github.com/murakmii/gojiai/vm"
)

func main() {
	v := vm.NewVM()
	ret, err := v.Execute(
		"MakeJVM.class",
		"compute",
		"(I)I",
		[]interface{}{int32(10)},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("return: %+v\n", ret)
}
