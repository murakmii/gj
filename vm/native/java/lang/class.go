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

func ClassIsInterface(thread *vm.Thread, args []interface{}) error {
	className := args[0].(*vm.Instance).VMData().(*string)

	class, err := thread.VM().FindClass(className)
	if err != nil {
		return err
	}

	result := 0
	if class.File().AccessFlag().Contain(class_file.InterfaceFlag) {
		result = 1
	}

	thread.CurrentFrame().PushOperand(result)
	return nil
}

func ClassIsPrimitive(thread *vm.Thread, args []interface{}) error {
	result := 0
	if class_file.JavaTypeSignature(*(args[0].(*vm.Instance).VMData().(*string))).IsPrimitive() {
		result = 1
	}

	thread.CurrentFrame().PushOperand(result)
	return nil
}

func ClassGetModifiers(thread *vm.Thread, args []interface{}) error {
	class, err := thread.VM().FindClass(args[0].(*vm.Instance).VMData().(*string))
	if err != nil {
		return err
	}

	thread.CurrentFrame().PushOperand(int(class.File().AccessFlag()))
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
	name := strings.ReplaceAll(vm.JavaStringToGoString(args[0].(*vm.Instance)), ".", "/")
	sig := class_file.JavaTypeSignature(name)

	if !sig.IsPrimitive() && !sig.IsArray() {
		_, state, err := thread.VM().FindInitializedClass(&name, thread)
		if err != nil {
			return err
		}
		if state == vm.FailedInitialization {
			return fmt.Errorf("failed initialization of class class in Class.getDeclaredFields0")
		}
	}

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

func ClassGetDeclaredConstructors(thread *vm.Thread, args []interface{}) error {
	class, state, err := thread.VM().FindInitializedClass(args[0].(*vm.Instance).VMData().(*string), thread)
	if err != nil {
		return err
	}
	if state == vm.FailedInitialization {
		return fmt.Errorf("failed initialization of class class in Class.getDeclaredConstrutors0")
	}

	pubOnly := args[1].(int) == 1
	cstrs := make([]*class_file.MethodInfo, 0)
	for _, m := range class.File().AllMethods() {
		if (*m.Name()) == "<init>" && (!pubOnly || m.IsPublic()) {
			cstrs = append(cstrs, m)
		}
	}

	cstrClassName := "java/lang/reflect/Constructor"
	cstrClass, state, err := thread.VM().FindInitializedClass(&cstrClassName, thread)
	if err != nil {
		return err
	}
	if state == vm.FailedInitialization {
		return fmt.Errorf("failed initialization of class class in Class.getDeclaredConstrutors0")
	}

	_, cstr := cstrClass.ResolveMethod("<init>", "(Ljava/lang/Class;[Ljava/lang/Class;[Ljava/lang/Class;IILjava/lang/String;[B[B)V")
	ret := vm.NewArray("Ljava/lang/reflect/Constructor;", len(cstrs))

	for i, c := range cstrs {
		cInstance := vm.NewInstance(cstrClass)

		var signature *vm.Instance
		sig, ok := c.Signature()
		if ok {
			signature = vm.GoString(*class.File().ConstantPool().Utf8(uint16(sig))).ToJavaString(thread)
		}

		params := class_file.ParseDescriptor2([]byte(*c.Descriptor()))
		pArray := vm.NewArray("Ljava/lang/Class", len(params))
		for i, p := range params {
			jts := class_file.JavaTypeSignature(p)
			if jts.IsReference() && !jts.IsArray() {
				_, state, err := thread.VM().FindInitializedClass(&p, thread)
				if err != nil {
					return err
				}
				if state == vm.FailedInitialization {
					return fmt.Errorf("failed initialization of class class in Class.getDeclaredConstrutors0")
				}
			}

			pArray.Set(i, thread.VM().JavaClass(&p))
		}

		exceptions := c.Exceptions()
		eArray := vm.NewArray("Ljava/lang/Class", len(exceptions))
		for i, e := range exceptions {
			eName := class.File().ConstantPool().ClassInfo(e)
			_, state, err := thread.VM().FindInitializedClass(eName, thread)
			if err != nil {
				return err
			}
			if state == vm.FailedInitialization {
				return fmt.Errorf("failed initialization of class class in Class.getDeclaredConstrutors0")
			}

			eArray.Set(i, thread.VM().JavaClass(eName))
		}

		thrown, err := thread.Derive().Execute(vm.NewFrame(cstrClass, cstr).SetLocals([]interface{}{
			cInstance,
			args[0],
			pArray,
			eArray,
			int(c.AccessFlag()),
			0,
			signature,
			vm.ByteSliceToJavaArray(c.RawAnnotations()),
			vm.ByteSliceToJavaArray(c.RawParamAnnotations()),
		}))
		if err != nil {
			return err
		}
		if thrown != nil {
			return fmt.Errorf("exception in Class.getDeclaredConstrutors0: %+v", thrown)
		}

		ret.Set(i, cInstance)
	}

	thread.CurrentFrame().PushOperand(ret)
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
