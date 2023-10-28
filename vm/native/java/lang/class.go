package lang

import (
	"fmt"
	"github.com/murakmii/gj/class_file"
	"github.com/murakmii/gj/vm"
	"strings"
)

func ClassDesiredAssertionStatus0(thread *vm.Thread, _ []interface{}) error {
	// return false
	thread.CurrentFrame().PushOperand(int32(0))
	return nil
}

func ClassGetPrimitiveClass(thread *vm.Thread, args []interface{}) error {
	className := args[0].(*vm.Instance).GetCharArrayField("value")
	class, err := thread.VM().Class(className, thread)
	if err != nil {
		return err
	}

	thread.CurrentFrame().PushOperand(class.Java())
	return nil
}

func ClassIsAssignableFrom(thread *vm.Thread, args []interface{}) error {
	thisName := args[0].(*vm.Instance).AsClass().File().ThisClass()
	argClass := args[1].(*vm.Instance).AsClass()

	var result int32
	if argClass.IsSubClassOf(&thisName) || argClass.Implements(&thisName) {
		result = 1
	}

	thread.CurrentFrame().PushOperand(result)
	return nil
}

func ClassIsArray(thread *vm.Thread, args []interface{}) error {
	var ret int32
	if args[0].(*vm.Instance).AsClass().IsArray() {
		ret = 1
	}

	thread.CurrentFrame().PushOperand(ret)
	return nil
}

func ClassIsInterface(thread *vm.Thread, args []interface{}) error {
	class := args[0].(*vm.Instance).AsClass()

	var result int32
	if class.File().AccessFlag().Contain(class_file.InterfaceFlag) {
		result = 1
	}

	thread.CurrentFrame().PushOperand(result)
	return nil
}

func ClassIsInstance(thread *vm.Thread, args []interface{}) error {
	className := args[0].(*vm.Instance).AsClass().File().ThisClass()

	var result int32
	if args[1].(*vm.Instance).Class().IsSubClassOf(&className) {
		result = 1
	}

	thread.CurrentFrame().PushOperand(result)
	return nil
}

func ClassIsPrimitive(thread *vm.Thread, args []interface{}) error {
	className := args[0].(*vm.Instance).AsClass().File().ThisClass()

	var result int32
	if class_file.JavaTypeSignature(className).IsPrimitive() {
		result = 1
	}

	thread.CurrentFrame().PushOperand(result)
	return nil
}

func ClassGetComponentType(thread *vm.Thread, args []interface{}) error {
	className := args[0].(*vm.Instance).AsClass().File().ThisClass()
	component := class_file.FieldType((className)[1:]).Type()

	class, err := thread.VM().Class(component, thread)
	if err != nil {
		return err
	}

	thread.CurrentFrame().PushOperand(class.Java())
	return nil
}

func ClassGetModifiers(thread *vm.Thread, args []interface{}) error {
	class := args[0].(*vm.Instance).AsClass()
	thread.CurrentFrame().PushOperand(int32(class.File().AccessFlag()))
	return nil
}

func ClassGetName0(thread *vm.Thread, args []interface{}) error {
	name := strings.ReplaceAll(args[0].(*vm.Instance).AsClass().File().ThisClass(), "/", ".")
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
	name := strings.ReplaceAll(args[0].(*vm.Instance).GetCharArrayField("value"), ".", "/")

	class, err := thread.VM().Class(name, thread)
	if err != nil {
		return err
	}

	thread.CurrentFrame().PushOperand(class.Java())
	return nil
}

func ClassGetSuperClass(thread *vm.Thread, args []interface{}) error {
	class := args[0].(*vm.Instance).AsClass()
	if class.Super() == nil {
		thread.CurrentFrame().PushOperand(nil)
		return nil
	}

	thread.CurrentFrame().PushOperand(class.Super().Java())
	return nil
}

func ClassGetDeclaredConstructors(thread *vm.Thread, args []interface{}) error {
	class := args[0].(*vm.Instance).AsClass()
	pubOnly := args[1].(int32) == 1

	cstrs := make([]*class_file.MethodInfo, 0)
	for _, m := range class.File().AllMethods() {
		if (*m.Name()) == "<init>" && (!pubOnly || m.IsPublic()) {
			cstrs = append(cstrs, m)
		}
	}

	cstrClass, err := thread.VM().Class("java/lang/reflect/Constructor", thread)
	if err != nil {
		return err
	}

	_, cstr := cstrClass.ResolveMethod("<init>", "(Ljava/lang/Class;[Ljava/lang/Class;[Ljava/lang/Class;IILjava/lang/String;[B[B)V")
	ret, retSlice := vm.NewArray(thread.VM(), "[Ljava/lang/reflect/Constructor;", len(cstrs))

	for i, c := range cstrs {
		cInstance := vm.NewInstance(cstrClass)

		var signature *vm.Instance
		sig, ok := c.Signature()
		if ok {
			signature = vm.GoString(*class.File().ConstantPool().Utf8(uint16(sig))).ToJavaString(thread.VM())
		}

		params := c.Descriptor().Params()
		pArray, pSlice := vm.NewArray(thread.VM(), "[Ljava/lang/Class;", len(params))
		for i, p := range params {
			class, err := thread.VM().Class(p.Type(), thread)
			if err != nil {
				return err
			}

			pSlice[i] = class.Java()
		}

		exceptions := c.Exceptions()
		eArray, eSlice := vm.NewArray(thread.VM(), "[Ljava/lang/Class;", len(exceptions))
		for i, e := range exceptions {
			eName := class.File().ConstantPool().ClassInfo(e)
			eClass, err := thread.VM().Class(*eName, thread)
			if err != nil {
				return err
			}

			eSlice[i] = eClass.Java()
		}

		err = thread.Execute(vm.NewFrame(cstrClass, cstr).SetLocals([]interface{}{
			cInstance,
			args[0],
			pArray,
			eArray,
			int32(c.AccessFlag()),
			int32(c.ID()),
			signature,
			vm.ByteSliceToJavaArray(thread.VM(), c.RawAnnotations()),
			vm.ByteSliceToJavaArray(thread.VM(), c.RawParamAnnotations()),
		}))
		if err != nil {
			return err
		}

		retSlice[i] = cInstance
	}

	thread.CurrentFrame().PushOperand(ret)
	return nil
}

func ClassGetDeclaredFields0(thread *vm.Thread, args []interface{}) error {
	class := args[0].(*vm.Instance)
	targetClass := class.AsClass()
	pubOnly := args[1].(int32) == 1

	var fields []*class_file.FieldInfo
	for _, f := range targetClass.File().AllFields() {
		if !pubOnly || f.AccessFlag().Contain(class_file.PublicFlag) {
			fields = append(fields, f)
		}
	}

	fieldClass, err := thread.VM().Class("java/lang/reflect/Field", thread)
	if err != nil {
		return err
	}

	_, cstr := fieldClass.ResolveMethod("<init>", "(Ljava/lang/Class;Ljava/lang/String;Ljava/lang/Class;IILjava/lang/String;[B)V")
	ret, retSlice := vm.NewArray(thread.VM(), "[Ljava/lang/reflect/Field;", len(fields))

	for i, f := range fields {
		fInstance := vm.NewInstance(fieldClass)

		var signature *vm.Instance
		sig, ok := f.Signature()
		if ok {
			signature = vm.GoString(*targetClass.File().ConstantPool().Utf8(uint16(sig))).ToJavaString(thread.VM())
		}

		descClass, err := thread.VM().Class(f.Descriptor().Type(), thread)
		if err != nil {
			return err
		}

		err = thread.Execute(vm.NewFrame(fieldClass, cstr).SetLocals([]interface{}{
			fInstance,
			class,
			thread.VM().JavaString2(thread, f.Name()),
			descClass.Java(),
			int32(f.AccessFlag()),
			int32(f.ID()),
			signature,
			vm.ByteSliceToJavaArray(thread.VM(), f.RawAnnotations()),
		}))
		if err != nil {
			return err
		}

		retSlice[i] = fInstance
	}

	thread.CurrentFrame().PushOperand(ret)
	return nil
}

func ClassGetEnclosingMethod0(thread *vm.Thread, args []interface{}) error {
	class := args[0].(*vm.Instance).AsClass()
	enc := class.File().EnclosingMethod()

	if enc == nil {
		thread.CurrentFrame().PushOperand(nil)
		return nil
	}

	return fmt.Errorf("Class.getEnclosingMethod0 has NOT been implemented")
}

func ClassGetDeclaringClass0(thread *vm.Thread, args []interface{}) error {
	class := args[0].(*vm.Instance).AsClass()
	if len(class.File().InnerClassesAttr()) == 0 {
		thread.CurrentFrame().PushOperand(nil)
		return nil
	}

	return fmt.Errorf("Class.getDeclaringClass0 has NOT been implemented")
}
