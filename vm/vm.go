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

func (vm *VM) FindClass(name string) (*Class, error) {
	vm.classLock.Lock()
	defer vm.classLock.Unlock()

	if class, ok := vm.classCache[name]; ok {
		return class, nil
	}

	for _, classPath := range vm.classPaths {
		file, err := classPath.SearchClass(name)
		if err != nil {
			return nil, err
		}

		if file != nil {
			vm.classCache[name] = NewClass(file)
			return vm.classCache[name], nil
		}
	}

	return nil, fmt.Errorf("class not found")
}

func (vm *VM) initializeClasses(mainThread *Thread, classNames []string) error {
	for _, className := range classNames {
		class, err := vm.FindClass(className)
		if err != nil {
			return err
		}

		state, err := class.Initialize(mainThread)
		if err != nil {
			return err
		}

		if state == FailedInitialization {
			return fmt.Errorf("failed to initialize classes in VM initialization")
		}
	}

	return nil
}
