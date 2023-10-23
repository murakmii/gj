package vm

import (
	"fmt"
	"testing"
)

func TestX(t *testing.T) {
	c1 := make(chan struct{})
	c2 := make(chan struct{})
	c3 := c1

	fmt.Printf("c1 == c2 = %t", c1 == c2)
	fmt.Printf("c2 == c3 = %t", c3 == c2)
	fmt.Printf("c1 == c3 = %t", c3 == c1)
}
