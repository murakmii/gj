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
	}
}

func (class *Class) File() *class_file.ClassFile {
	return class.file
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
	defer func() {
		class.initCond.L.Lock()
		class.state = state
		class.initBy = nil
		class.initCond.Broadcast()
		class.initCond.L.Unlock()
	}()

	// Initialize constant fields
	for _, f := range class.file.Fields() {
		if constValAttr, ok := f.ConstantValue(); ok {
			fieldName := class.file.ConstantPool().Utf8(f.Name())
			constVal := class.file.ConstantPool().Const(constValAttr)

			switch cv := constVal.(type) {
			case *string:
			// TODO: set java.lang.String instance
			default:
				class.fields[*fieldName] = cv
			}
		}
	}

	// Initialize super class and interfaces
	for _, d := range class.file.DependencyClasses() {
		dc, err := curThread.VM().FindClass(d)
		if err != nil {
			return err
		}

		state, err := dc.Initialize(curThread)
		if err != nil {
			return err
		}

		if state == FailedInitialization {
			return nil
		}
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
