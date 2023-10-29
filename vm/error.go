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
	return &JavaError{
		message:   exception.GetField("detailMessage", "Ljava/lang/String;").(*Instance).AsString(),
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
	err = thread.Execute(NewFrame(constrClass, constr).SetLocals([]interface{}{ex, NewString(thread.VM(), message)}))
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
