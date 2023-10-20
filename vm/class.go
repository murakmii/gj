package vm

import (
	"fmt"
	"github.com/murakmii/gj/class_file"
	"strings"
	"sync"
)

type (
	Class struct {
		file         *class_file.ClassFile
		fields       []interface{}
		totalIFields int
		state        ClassState
		initCond     *sync.Cond
		initBy       *Thread

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
		file:         file,
		fields:       make([]interface{}, len(file.AllFields())-len(file.InstanceFields())),
		totalIFields: -1,
		state:        NotInitialized,
		initCond:     sync.NewCond(&sync.Mutex{}),
		initBy:       nil,

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

func (class *Class) IsInstanceOf(className *string) bool {
	return class.IsSubClassOf(className) || class.Implements(className)
}

func (class *Class) Implements(ifName *string) bool {
	for _, i := range class.file.Interfaces() {
		if *i == *ifName {
			return true
		}
	}
	return class.super != nil && class.super.Implements(ifName)
}

func (class *Class) TotalInstanceFields() int {
	return class.totalIFields
}

func (class *Class) SetStaticField(field *class_file.FieldInfo, value interface{}) {
	class.fields[field.ID()] = value
}

func (class *Class) GetStaticField(field *class_file.FieldInfo) interface{} {
	value := class.fields[field.ID()]
	if value == nil && !field.NullableDefaultValue() {
		class.fields[field.ID()] = field.DefaultValue()
		value = field.DefaultValue()
	}
	return value
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
		if curThread.Equal(class.initBy) {
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

	if _, err := class.initializeFieldID(curThread.VM()); err != nil {
		return err
	}

	// Initialize constant fields
	for _, f := range class.file.StaticFields() {
		if constValAttr, ok := f.ConstantValue(); ok {
			constVal := class.file.ConstantPool().Const(uint16(constValAttr))

			switch cv := constVal.(type) {
			case *string:
				class.fields[f.ID()], err = curThread.VM().JavaString(curThread, cv)
				if err != nil {
					return fmt.Errorf("failed to set default string value of static field: %s", err)
				}
			default:
				class.fields[f.ID()] = cv
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
		unCatchEx, err := curThread.Derive().Execute(NewFrame(class, clinit))
		if err != nil {
			return err
		}

		if unCatchEx != nil {
			state = FailedInitialization

			fmt.Println("------------------------")

			//detail := "detailMessage"
			//detailDesc := "Ljava/lang/String;"
			//fmt.Printf("!!! exception: %s\n\n", JavaStringToGoString(unCatchEx.GetField(&detail, &detailDesc).(*Instance)))

			for i, t := range unCatchEx.VMData().([]string) {
				fmt.Printf("%s%s\n", strings.Repeat(" ", i), t)
			}

			fmt.Println("------------------------")

			return nil
		}
	}

	state = Initialized
	return nil
}

func (class *Class) initializeFieldID(vm *VM) (int, error) {
	if class.totalIFields != -1 {
		return class.totalIFields, nil
	}

	id := 0
	if class.file.SuperClass() != nil {
		super, err := vm.FindClass(class.file.SuperClass())
		if err != nil {
			return -1, err
		}
		id, err = super.initializeFieldID(vm)
		if err != nil {
			return -1, err
		}
	}

	for _, f := range class.file.InstanceFields() {
		f.SetID(id)
		id++
	}

	class.totalIFields = id
	return id, nil
}
