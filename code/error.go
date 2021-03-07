package code

import (
	"errors"
	"fmt"

	"golang.org/x/xerrors"
)

type geekfesError struct {
	code ErrCode
	err  error
}

func (e *geekfesError) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.code, e.err)
}

func Errorf(code ErrCode, format string, v ...interface{}) error {
	if code == OK {
		return nil
	}
	return &geekfesError{
		code: code,
		err:  xerrors.Errorf(format, v...),
	}
}

func From(err error) ErrCode {
	if err == nil {
		return OK
	}
	var e *geekfesError
	if errors.As(err, &e) {
		return e.code
	}
	return Internal
}
