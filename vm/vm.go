package vm

import (
	"fmt"
	"github.com/murakmii/gj"
	"sync"
	"unicode/utf16"
)

type VM struct {
	sysProps map[string]string

	classPaths []gj.ClassPath
	classCache map[string]*Class
	classLock  *sync.Mutex

	mainThread *Thread

	jlString        *Class
	jlClass         *Class
	javaStringCache map[string]*Instance
}

func InitVM(config *gj.Config) (*VM, error) {
	var err error
	vm := &VM{
		sysProps:        config.SysProps,
		classCache:      make(map[string]*Class),
		classLock:       &sync.Mutex{},
		javaStringCache: make(map[string]*Instance),
	}
	vm.mainThread = NewThread(vm)

	vm.classPaths, err = gj.InitClassPaths(config.ClassPath)
	if err != nil {
		return nil, err
	}

	classes, err := vm.initializeClasses([]string{
		"java/lang/String",
		"java/lang/System",
		"java/lang/Class",
		"java/lang/ThreadGroup",
	})
	if err != nil {
		return nil, err
	}

	vm.jlString = classes[0]
	vm.jlClass = classes[2]

	if err = vm.initializeMainThread(); err != nil {
		return nil, err
	}

	if err = vm.initializeSystemClass(); err != nil {
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

func (vm *VM) SysProps() map[string]string {
	return vm.sysProps
}

func (vm *VM) JavaLangStringClass() *Class {
	return vm.jlString
}

func (vm *VM) JavaLangClassClass() *Class {
	return vm.jlClass
}

func (vm *VM) FindInitializedClass(name *string, curThread *Thread) (*Class, ClassState, error) {
	class, err := vm.FindClass(name)
	if err != nil {
		return nil, NotInitialized, err
	}

	state, err := class.Initialize(curThread)
	return class, state, err
}

func (vm *VM) JavaString2(thread *Thread, s *string) *Instance {
	if cache, ok := vm.javaStringCache[*s]; ok {
		return cache
	}

	js := GoString(*s).ToJavaString(thread)
	vm.javaStringCache[*s] = js
	return js
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

func (vm *VM) initializeClasses(classNames []string) ([]*Class, error) {
	classes := make([]*Class, len(classNames))
	var err error
	var state ClassState

	for i, className := range classNames {
		classes[i], state, err = vm.FindInitializedClass(&className, vm.mainThread)
		if err != nil {
			return nil, err
		}
		if state == FailedInitialization {
			return nil, fmt.Errorf("class '%s' initialization failed", className)
		}
	}

	return classes, nil
}

func (vm *VM) initializeMainThread() error {
	// Create system thread group.
	tgClassName := "java/lang/ThreadGroup"
	tgClass, err := vm.FindClass(&tgClassName)
	if err != nil {
		return err
	}

	sysTg := NewInstance(tgClass)
	frame := NewFrame(tgClass, tgClass.File().FindMethod("<init>", "()V")).SetLocal(0, sysTg)
	thrown, err := vm.mainThread.Derive().Execute(frame)
	if err != nil {
		return err
	}
	if thrown != nil {
		return fmt.Errorf("failed to construct system thread group. thrown: %+v", thrown)
	}

	// Create main thread group.
	mainStr := "main"
	mainJs, err := vm.JavaString(vm.mainThread, &mainStr)
	if err != nil {
		return err
	}

	mainTg := NewInstance(tgClass)
	frame = NewFrame(tgClass, tgClass.File().FindMethod("<init>", "(Ljava/lang/ThreadGroup;Ljava/lang/String;)V")).
		SetLocals([]interface{}{mainTg, sysTg, mainJs})
	thrown, err = vm.mainThread.Derive().Execute(frame)
	if err != nil {
		return err
	}
	if thrown != nil {
		return fmt.Errorf("failed to construct main thread group. thrown: %+v", thrown)
	}

	// Create main thread.
	tClassName := "java/lang/Thread"
	tClass, err := vm.FindClass(&tClassName)
	if err != nil {
		return err
	}

	mainJThread := NewInstance(tClass)
	threadPriorityField := "priority"
	mainJThread.PutField(&tClassName, &threadPriorityField, 5)
	vm.mainThread.SetJavaThread(mainJThread)

	frame = NewFrame(tClass, tClass.File().FindMethod("<init>", "(Ljava/lang/ThreadGroup;Ljava/lang/String;)V")).
		SetLocals([]interface{}{mainJThread, mainTg, mainJs})
	thrown, err = vm.mainThread.Derive().Execute(frame)
	if err != nil {
		return err
	}
	if thrown != nil {
		return fmt.Errorf("failed to construct main thread. thrown: %+v", thrown)
	}

	return nil
}

func (vm *VM) initializeSystemClass() error {
	sysClassName := "java/lang/System"
	sys, state, err := vm.FindInitializedClass(&sysClassName, vm.mainThread)
	if err != nil {
		return err
	}
	if state == FailedInitialization {
		return fmt.Errorf("failed initialization for java/lang/System")
	}

	frame := NewFrame(sys, sys.File().FindMethod("initializeSystemClass", "()V"))
	thrown, err := vm.mainThread.Derive().Execute(frame)
	if err != nil {
		return err
	}
	if thrown != nil {
		return fmt.Errorf("failed to call java/lang/System.initializeSystemClass. thrown: %+v", thrown)
	}

	return nil
}
