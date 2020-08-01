package katana

import (
	"fmt"
)

// Error ...
type Error interface {
	No() int
	Error() string
}

// Err ...
type Err struct {
	no  int
	msg string
}

// No ...
func (e *Err) No() int {
	return e.no
}

func (e *Err) Error() string {
	return e.msg
}

func (e *Err) copy() *Err {
	return &Err{
		no:  e.no,
		msg: e.msg,
	}
}

// Concat ...
func (e *Err) Concat(format string, a ...interface{}) *Err {
	err := e.copy()
	err.msg += fmt.Sprintf(format, a...)
	return err
}

// NewErr ...
func newErr(no int, msg string) *Err {
	return &Err{
		no:  no,
		msg: msg,
	}
}

// err var define
var (
	errSuccess    = newErr(0, "success")
	errBadRequest = newErr(-1, "bad request")
	errInternal   = newErr(-1, "internal error")
)
