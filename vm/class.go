package vm

import (
	"github.com/murakmii/gj/class_file"
	"sync"
)

type (
	Class struct {
		file     *class_file.ClassFile
		fields   map[string]interface{}
		state    ClassState
		initCond *sync.Cond
		initBy   *Thread

		super      *Class
		interfaces []*Class
	}

	ClassState uint8
)

const (
	NotInitialized ClassState = iota
	Initializing
	Initialized
	FailedInitialization
)

func NewClass(file *class_file.ClassFile) *Class {
	return &Class{
		file:     file,
		fields:   make(map[string]interface{}),
		state:    NotInitialized,
		initCond: sync.NewCond(&sync.Mutex{}),
		initBy:   nil,

		// Set in initialize method
		super:      nil,
		interfaces: nil,
	}
}

func (class *Class) File() *class_file.ClassFile {
	return class.file
}

func (class *Class) Super() *Class {
	return class.super
}

func (class *Class) IsSubClassOf(className *string) bool {
	return class.file.ThisClass() == *className || (class.super != nil && class.super.IsSubClassOf(className))
}

func (class *Class) SetStaticField(name *string, value interface{}) {
	class.fields[*name] = value
}

func (class *Class) GetStaticField(name *string) interface{} {
	return class.fields[*name]
}

// See: https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-5.html#jvms-5.4.3.2
func (class *Class) ResolveField(name, desc string) (*Class, *class_file.FieldInfo) {
	// Step.1 find from this class
	if field := class.file.FindField(name, desc); field != nil {
		return class, field
	}

	// Step.2 find from super interfaces
	for _, ifClass := range class.interfaces {
		if field := ifClass.File().FindField(name, desc); field != nil {
			return ifClass, field
		}
	}

	// Step.3 find from super class
	if class.super != nil {
		return class.super.ResolveField(name, desc)
	}

	return nil, nil
}

// See: https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-5.html#jvms-5.4.3.3
func (class *Class) ResolveMethod(name, desc string) (*Class, *class_file.MethodInfo) {
	if method := class.file.FindMethod(name, desc); method != nil {
		return class, method
	}

	if class.super != nil {
		if super, method := class.super.ResolveMethod(name, desc); super != nil {
			return super, method
		}
	}

	for _, ifClass := range class.interfaces {
		if resolvedIf, method := ifClass.ResolveMethod(name, desc); method != nil {
			return resolvedIf, method
		}
	}

	return nil, nil
}

// Initializes class and return state of class.
// This method implements initialization process of JVM spec
// See: https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-5.html#jvms-5.5
func (class *Class) Initialize(curThread *Thread) (ClassState, error) {
	class.initCond.L.Lock()

	switch class.state {
	case NotInitialized:
		class.state = Initializing
		class.initBy = curThread
		class.initCond.L.Unlock()

		err := class.initialize(curThread)
		if err != nil {
			return NotInitialized, err
		}
		return class.state, nil

	case Initializing:
		if class.initBy == curThread {
			class.initCond.L.Unlock()
			return class.state, nil
		}
		class.initCond.Wait()
		return class.state, nil

	default:
		class.initCond.L.Unlock()
		return class.state, nil
	}
}

func (class *Class) initialize(curThread *Thread) error {
	var state ClassState
	var err error

	defer func() {
		class.initCond.L.Lock()
		class.state = state
		class.initBy = nil
		class.initCond.Broadcast() // Wake up all threads are waiting initialization of this class
		class.initCond.L.Unlock()
	}()

	// Initialize constant fields
	for _, f := range class.file.Fields(class_file.StaticFlag & class_file.FinalFlag) {
		if constValAttr, ok := f.ConstantValue(); ok {
			constVal := class.file.ConstantPool().Const(constValAttr)

			switch cv := constVal.(type) {
			case *string:
			// TODO: set java.lang.String instance
			default:
				class.fields[*f.Name()] = cv
			}
		}
	}

	// Initialize super class
	if class.file.SuperClass() != nil {
		class.super, state, err = curThread.VM().FindInitializedClass(class.file.SuperClass(), curThread)
		if err != nil || state == FailedInitialization {
			return err
		}
	}

	// Initialize interfaces
	var ifClass *Class
	for _, ifName := range class.file.Interfaces() {
		ifClass, state, err = curThread.VM().FindInitializedClass(ifName, curThread)
		if err != nil || state == FailedInitialization {
			return err
		}
		class.interfaces = append(class.interfaces, ifClass)
	}

	// Call clinit
	clinit := class.file.FindMethod("<clinit>", "()V")
	if clinit != nil {
		frame := NewFrame(curThread, class, clinit)
		frameOp, err := frame.Execute()
		if err != nil {
			return err
		}

		if frameOp == ThrowFromFrame {
			state = FailedInitialization
			return nil
		}
	}

	state = Initialized
	return nil
}
