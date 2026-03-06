package types

import "errors"

var (
	ErrorNotFound        = errors.New("not found")
	ErrorAlreadyExist    = errors.New("already exist")
	ErrorAccessDenied    = errors.New("access denied")
	ErrorInvalidArgument = errors.New("invalid incoming data")
	ErrorInternal        = errors.New("internal error")
)
