package main

import (
	"github.com/murakmii/gj/class_file"
	"os"
)

func main() {
	_, err := class_file.ReadClassFile(os.Args[1])
	if err != nil {
		panic(err)
	}
}
