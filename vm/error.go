package vm

import "fmt"

type JavaError struct {
	original  error
	message   string
	exception string
}

var _ error = (*JavaError)(nil)

func NewJavaError(exception, message string) *JavaError {
	return &JavaError{message: message, exception: exception}
}

func (e *JavaError) Error() string {
	return e.exception + ": " + e.message
}

func (e *JavaError) CreateException(thread *Thread) (*Instance, error) {
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
}
