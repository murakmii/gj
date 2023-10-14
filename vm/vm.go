package vm

import (
	"fmt"
	"github.com/murakmii/gj"
	"sync"
)

type VM struct {
	classPaths []gj.ClassPath
	classCache map[string]*Class
	classLock  *sync.Mutex
}

func InitVM(config *gj.Config) (*VM, error) {
	var err error
	vm := &VM{
		classCache: make(map[string]*Class),
		classLock:  &sync.Mutex{},
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
		file, err := classPath.SearchClass(*name)
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
