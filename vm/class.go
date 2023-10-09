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

		return 0, nil

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

	}()

	// TODO: init constant field

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

	// TODO: call clinit

	state = Initialized
	return nil
}
