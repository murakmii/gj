package lang

import (
	"fmt"
	"github.com/murakmii/gj/class_file"
	"github.com/murakmii/gj/vm"
	"strings"
)

func ClassDesiredAssertionStatus0(thread *vm.Thread, args []interface{}) error {
	// return false
	thread.CurrentFrame().PushOperand(0)
	return nil
}

func ClassGetPrimitiveClass(thread *vm.Thread, _ []interface{}) error {
	// TODO: generate class instance
	// return null
	thread.CurrentFrame().PushOperand(nil)
	return nil
}

func ClassIsAssignableFrom(thread *vm.Thread, args []interface{}) error {
	thisName := args[0].(*vm.Instance).VMData().(*string)
	argName := args[1].(*vm.Instance).VMData().(*string)

	if (*thisName)[0] != 'L' || (*argName)[0] != 'L' {
		thread.CurrentFrame().PushOperand(0)
		return nil
	}

	argClass, err := thread.VM().FindClass(argName)
	if err != nil {
		return err
	}

	result := 0
	if argClass.IsSubClassOf(thisName) || argClass.Implements(thisName) {
		result = 1
	}

	thread.CurrentFrame().PushOperand(result)
	return nil
}

func ClassIsPrimitive(thread *vm.Thread, args []interface{}) error {
	className := *(args[0].(*vm.Instance).VMData().(*string))

	result := 0
	if className[0] != 'L' && className[0] != '[' {
		result = 1
	}

	thread.CurrentFrame().PushOperand(result)
	return nil
}

func ClassGetName0(thread *vm.Thread, args []interface{}) error {
	javaClass := args[0].(*vm.Instance)

	name := strings.ReplaceAll(*javaClass.VMData().(*string), "/", ".")
	if name[0] == 'L' {
		name = name[1 : len(name)-1]
	}

	nameJS, err := thread.VM().JavaString(thread, &name)
	if err != nil {
		return err
	}

	thread.CurrentFrame().PushOperand(nameJS)
	return nil
}

func ClassForName0(thread *vm.Thread, args []interface{}) error {
	name := vm.JavaStringToGoString(args[0].(*vm.Instance))
	thread.CurrentFrame().PushOperand(thread.VM().JavaClass(&name))
	return nil
}

func ClassGetSuperClass(thread *vm.Thread, args []interface{}) error {
	classInstance := args[0].(*vm.Instance)

	class, state, err := thread.VM().FindInitializedClass(classInstance.VMData().(*string), thread)
	if err != nil {
		return err
	}
	if state == vm.FailedInitialization {
		return fmt.Errorf("failed initialization of class class in Class.getDeclaredFields0")
	}

	if class.Super() == nil {
		thread.CurrentFrame().PushOperand(nil)
		return nil
	}

	superName := class.Super().File().ThisClass()
	thread.CurrentFrame().PushOperand(thread.VM().JavaClass(&superName))
	return nil
}

func ClassGetDeclaredFields0(thread *vm.Thread, args []interface{}) error {
	class := args[0].(*vm.Instance)
	pubOnly := args[1].(int) == 1

	targetClass, state, err := thread.VM().FindInitializedClass(class.VMData().(*string), thread)
	if err != nil {
		return err
	}
	if state == vm.FailedInitialization {
		return fmt.Errorf("failed initialization of class class in Class.getDeclaredFields0")
	}

	var fields []*class_file.FieldInfo
	for _, f := range targetClass.File().AllFields() {
		if !pubOnly || f.AccessFlag().Contain(class_file.PublicFlag) {
			fields = append(fields, f)
		}
	}

	fieldClassName := "java/lang/reflect/Field"
	fieldClass, state, err := thread.VM().FindInitializedClass(&fieldClassName, thread)
	if err != nil {
		return err
	}
	if state == vm.FailedInitialization {
		return fmt.Errorf("failed initialization of class class in Class.getDeclaredFields0")
	}

	_, cstr := fieldClass.ResolveMethod("<init>", "(Ljava/lang/Class;Ljava/lang/String;Ljava/lang/Class;IILjava/lang/String;[B)V")
	ret := vm.NewArray("Ljava/lang/reflect/Field;", len(fields))

	for i, f := range fields {
		fInstance := vm.NewInstance(fieldClass)

		var signature *vm.Instance
		sig, ok := f.Signature()
		if ok {
			signature = vm.GoString(*targetClass.File().ConstantPool().Utf8(uint16(sig))).ToJavaString(thread)
		}

		annotation := vm.NewArray("B", len(f.RawAnnotations()))
		for i, b := range f.RawAnnotations() {
			annotation.Set(i, int(b))
		}

		thrown, err := thread.Derive().Execute(vm.NewFrame(fieldClass, cstr).SetLocals([]interface{}{
			fInstance,
			class,
			thread.VM().JavaString2(thread, f.Name()),
			thread.VM().JavaClass(f.Desc()),
			int(f.AccessFlag()),
			f.ID(),
			signature,
			annotation,
		}))
		if err != nil {
			return err
		}
		if thrown != nil {
			return fmt.Errorf("exception in Class.getDeclaredFields0: %+v", thrown)
		}

		ret.Set(i, fInstance)
	}

	thread.CurrentFrame().PushOperand(ret)
	return nil
}
