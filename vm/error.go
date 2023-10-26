package vm

import "errors"

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
	return nil
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
