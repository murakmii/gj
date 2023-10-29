package vm

import (
	"fmt"
	"github.com/murakmii/gojiai"
	"sync"
)

type (
	VM struct {
		sysProps map[string]string

		classPaths        []gojiai.ClassPath
		classCache        map[string]*Class
		specialClassCache []*Class
		classLock         *sync.Mutex

		mainThread *Thread
		executor   *ThreadExecutor

		javaStringCache map[string]*Instance

		nativeMem *NativeMemAllocator
	}
)

func InitVM(config *gojiai.Config) (*VM, error) {
	var err error
	vm := &VM{
		sysProps:          config.SysProps,
		classCache:        make(map[string]*Class),
		specialClassCache: make([]*Class, 256),
		classLock:         &sync.Mutex{},
		executor:          NewThreadExecutor(),
		javaStringCache:   make(map[string]*Instance),
		nativeMem:         CreateNativeMemAllocator(),
	}
	vm.mainThread = NewThread(vm, "main", true, false)

	vm.classPaths, err = gojiai.InitClassPaths(config.ClassPath)
	if err != nil {
		return nil, err
	}

	classes, err := vm.initializeClasses([]string{
		"java/lang/Object",
		"java/lang/String",
		"java/lang/StackTraceElement",
		"java/lang/System",
		"java/lang/Class",
	})
	if err != nil {
		return nil, err
	}

	for _, class := range vm.classCache {
		class.InitJava(vm)
	}

	// Disable native library loading. Return(0xB1) immediately
	classes[3].File().FindMethod("loadLibrary", "(Ljava/lang/String;)V").
		Code().OverrideCode([]byte{0xB1})

	_, err = vm.initializeClasses([]string{"java/lang/ThreadGroup", "java/lang/Thread"})
	if err != nil {
		return nil, err
	}

	if err = vm.initializeMainThread(); err != nil {
		return nil, err
	}

	if err = vm.initializeSystemClass(); err != nil {
		return nil, err
	}

	return vm, nil
}

func (vm *VM) ClassCacheNum() int {
	return len(vm.classCache)
}

func (vm *VM) NativeMem() *NativeMemAllocator {
	return vm.nativeMem
}

func (vm *VM) SysProps() map[string]string {
	return vm.sysProps
}

func (vm *VM) DoneLoadingMinimumClass() bool {
	return vm.specialClassCache[JavaLangClassID] != nil && vm.specialClassCache[JavaLangClassID].State() == Initialized
}

func (vm *VM) SpecialClass(id SpecialClassID) *Class {
	return vm.specialClassCache[id]
}

func (vm *VM) Executor() *ThreadExecutor {
	return vm.executor
}

func (vm *VM) Class(className string, thread *Thread) (*Class, error) {
	class, ok := vm.classCache[className]
	if ok {
		if thread != nil {
			state, err := class.Initialize(thread)
			if err != nil {
				return nil, fmt.Errorf("failed to initialize class '%s': %w", className, err)
			}
			if state == FailedInitialization {
				// TODO: return JavaError
				panic("failed to initialize class: " + className)
			}
		}

		return class, nil
	}

	vm.classLock.Lock()
	if className[0] == '[' {
		class = NewArrayClass(vm, className)

	} else if className == "byte" || className == "char" || className == "double" || className == "float" ||
		className == "int" || className == "long" || className == "short" || className == "boolean" {
		class = NewPrimitiveClass(vm, className)

	} else {
		for _, classPath := range vm.classPaths {
			classFile, err := classPath.SearchClass(className + ".class")
			if err != nil {
				return nil, err
			}
			if classFile != nil {
				class = NewClass(classFile)
			}
		}

		if class == nil {
			vm.classLock.Unlock()
			return nil, fmt.Errorf("class '%s' not found", className)
		}
	}

	vm.classCache[className] = class
	if !class.ID().IsUnknown() {
		vm.specialClassCache[class.ID()] = class
	}

	vm.classLock.Unlock()

	if thread != nil && class.State() == NotInitialized {
		state, err := class.Initialize(thread)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize class '%s': %w", className, err)
		}
		if state == FailedInitialization {
			// TODO: return JavaError
			panic("failed to initialize class: " + className)
		}
	}

	return class, nil
}

func (vm *VM) JavaString(s string) *Instance {
	// TODO: lock
	if cache, ok := vm.javaStringCache[s]; ok {
		return cache
	}

	js := NewString(vm, s)
	vm.javaStringCache[s] = js
	return js
}

func (vm *VM) ExecMain(className string, args []string) error {
	class, err := vm.Class(className, vm.mainThread)
	if err != nil {
		return err
	}

	_, main := class.ResolveMethod("main", "([Ljava/lang/String;)V")
	if main == nil {
		return fmt.Errorf("main method not found in %s", className)
	}

	array, slice := NewArray(vm, "[Ljava/lang/String;", len(args))
	for i := range slice {
		slice[i] = NewString(vm, args[i])
	}

	vm.executor.Start(vm.mainThread, NewFrame(class, main).SetLocal(0, array))
	return nil
}

func (vm *VM) initializeClasses(classNames []string) ([]*Class, error) {
	classes := make([]*Class, len(classNames))
	var err error

	for i, className := range classNames {
		classes[i], err = vm.Class(className, vm.mainThread)
		if err != nil {
			return nil, err
		}
	}

	return classes, nil
}

func (vm *VM) initializeMainThread() error {
	// Create system thread group.
	tgClass, err := vm.Class("java/lang/ThreadGroup", nil)
	if err != nil {
		return err
	}

	sysTg := NewInstance(tgClass)
	frame := NewFrame(tgClass, tgClass.File().FindMethod("<init>", "()V")).SetLocal(0, sysTg)
	err = vm.mainThread.Execute(frame)
	if err != nil {
		return err
	}

	// Create main thread group.
	mainJs := vm.JavaString("main")

	mainTg := NewInstance(tgClass)
	frame = NewFrame(tgClass, tgClass.File().FindMethod("<init>", "(Ljava/lang/ThreadGroup;Ljava/lang/String;)V")).
		SetLocals([]interface{}{mainTg, sysTg, mainJs})
	err = vm.mainThread.Execute(frame)
	if err != nil {
		return err
	}

	// Create main thread.
	tClass, err := vm.Class("java/lang/Thread", nil)
	if err != nil {
		return err
	}

	mainJThread := NewInstance(tClass)
	mainJThread.PutField("priority", "I", int32(5))
	mainJThread.PutField("threadStatus", "I", int32(4)) // RUNNABLE
	mainJThread.ToBeThread(vm.mainThread)

	vm.mainThread.SetJavaThread(mainJThread)

	frame = NewFrame(tClass, tClass.File().FindMethod("<init>", "(Ljava/lang/ThreadGroup;Ljava/lang/String;)V")).
		SetLocals([]interface{}{mainJThread, mainTg, mainJs})
	err = vm.mainThread.Execute(frame)
	if err != nil {
		return err
	}

	return nil
}

func (vm *VM) initializeSystemClass() error {
	sys, err := vm.Class("java/lang/System", vm.mainThread)
	if err != nil {
		return err
	}

	frame := NewFrame(sys, sys.File().FindMethod("initializeSystemClass", "()V"))
	err = vm.mainThread.Execute(frame)
	if err != nil {
		return err
	}

	return nil
}
