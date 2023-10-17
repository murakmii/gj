package vm

import (
	"fmt"
	"github.com/murakmii/gj"
	"sync"
	"unicode/utf16"
)

type VM struct {
	classPaths []gj.ClassPath
	classCache map[string]*Class
	classLock  *sync.Mutex

	javaStringCache map[string]*Instance
}

func InitVM(config *gj.Config) (*VM, error) {
	var err error
	vm := &VM{
		classCache:      make(map[string]*Class),
		classLock:       &sync.Mutex{},
		javaStringCache: make(map[string]*Instance),
	}

	vm.classPaths, err = gj.InitClassPaths(config.ClassPath)
	if err != nil {
		return nil, err
	}

	return vm, nil
}

func (vm *VM) FindClass(name *string) (*Class, error) {
	vm.classLock.Lock()
	defer vm.classLock.Unlock()

	if class, ok := vm.classCache[*name]; ok {
		return class, nil
	}

	for _, classPath := range vm.classPaths {
		file, err := classPath.SearchClass(*name + ".class")
		if err != nil {
			return nil, err
		}

		if file != nil {
			vm.classCache[*name] = NewClass(file)
			return vm.classCache[*name], nil
		}
	}

	return nil, fmt.Errorf("class '%s' not found", *name)
}

func (vm *VM) FindInitializedClass(name *string, curThread *Thread) (*Class, ClassState, error) {
	class, err := vm.FindClass(name)
	if err != nil {
		return nil, NotInitialized, err
	}

	state, err := class.Initialize(curThread)
	return class, state, err
}

func (vm *VM) JavaString(thread *Thread, s *string) (*Instance, error) {
	if cache, ok := vm.javaStringCache[*s]; ok {
		return cache, nil
	}

	className := "java/lang/String"
	class, state, err := vm.FindInitializedClass(&className, thread)
	if err != nil {
		return nil, err
	}
	if state == FailedInitialization {
		return nil, fmt.Errorf("failed initialization for %s", className)
	}

	js := NewInstance(class)

	u16 := utf16.Encode([]rune(*s))
	charArray := NewArray("C", len(u16))
	for i, e := range u16 {
		charArray.Set(i, int(e))
	}

	fieldName := "value"
	js.PutField(&className, &fieldName, charArray)

	vm.javaStringCache[*s] = js
	return js, nil
}
