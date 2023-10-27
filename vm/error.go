package vm

import (
	"errors"
	"fmt"
)

type JavaError struct {
	message   string
	exception *Instance
}

var _ error = (*JavaError)(nil)

func UnwrapJavaError(err error) *JavaError {
	if javaErr, ok := err.(*JavaError); ok {
		return javaErr
	}

	next := errors.Unwrap(err)
	if next == nil {
		return nil
	}

	return UnwrapJavaError(next)
}

func NewJavaErr(exception *Instance) error {
	msgName := "detailMessage"
	msgDesc := "Ljava/lang/String;"

	return &JavaError{
		message:   exception.GetField(&msgName, &msgDesc).(*Instance).GetCharArrayField("value"),
		exception: exception,
	}
}

func CreateJavaError(thread *Thread, className, message string) error {
	exClass, err := thread.VM().Class(className, thread)
	if err != nil {
		return err
	}

	constrClass, constr := exClass.ResolveMethod("<init>", "(Ljava/lang/String;)V")
	if constr == nil {
		return fmt.Errorf("failed to resolve exception constructor for %s", className)
	}

	ex := NewInstance(exClass)
	err = thread.Execute(NewFrame(constrClass, constr).SetLocals([]interface{}{
		ex, GoString(message).ToJavaString(thread.VM()),
	}))
	if err != nil {
		return err
	}

	return &JavaError{message: message, exception: ex}
}

func (e *JavaError) Error() string {
	return e.exception.Class().File().ThisClass() + ": " + e.message
}

func (e *JavaError) Exception() *Instance {
	return e.exception
}

/*func (e *JavaError) CreateException(thread *Thread) (*Instance, error) {
	exClass, err := thread.VM().Class(e.exception, thread)
	if err != nil {
		return nil, err
	}

	_, cstr := exClass.ResolveMethod("<init>", "(Ljava/lang/String;)V")
	if cstr == nil {
		return nil, fmt.Errorf("failed to resolve exception constructor for %s", e.exception)
	}

	ex := NewInstance(exClass)
	_, err = thread.Derive().Execute(NewFrame(exClass, cstr).SetLocals([]interface{}{
		ex,
		GoString(e.message).ToJavaString(thread),
	}))
	if err != nil {
		return nil, err
	}

	return ex, nil
}*/
