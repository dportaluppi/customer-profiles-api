package pkg

type ErrIDType struct {
	msg string
}

func (e ErrIDType) Error() string {
	return e.msg
}

func NewErrID(msg string) ErrIDType {
	return ErrIDType{msg: msg}
}

type ErrInvalidType struct {
	msg string
}

func (e ErrInvalidType) Error() string {
	return e.msg
}

func NewErrInvalid(msg string) ErrInvalidType {
	return ErrInvalidType{msg: msg}
}

type ErrNotFoundType struct {
	msg string
}

func (e ErrNotFoundType) Error() string {
	return e.msg
}

func NewErrNotFound(msg string) ErrNotFoundType {
	return ErrNotFoundType{msg: msg}
}

type ErrConflictType struct {
	msg string
}

func (e ErrConflictType) Error() string {
	return e.msg
}

func NewErrConflict(msg string) ErrConflictType {
	return ErrConflictType{msg: msg}
}

type ErrInternalErrorType struct {
	msg string
}

func (e ErrInternalErrorType) Error() string {
	return e.msg
}

func NewErrInternalError(msg string) ErrInternalErrorType {
	return ErrInternalErrorType{msg: msg}
}
