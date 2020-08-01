package katana

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type errA struct{}

func (e *errA) No() int {
	return 0
}

func (e *errA) Error() string {
	return "A"
}

type errB struct{}

func (e *errB) No() int {
	return 1
}

func (e *errB) Error() string {
	return "B"
}

func f(Error) {}

func TestError(t *testing.T) {
	as := assert.New(t)
	as.True(false)

	a := &errA{}
	b := &errB{}
	var c Error
	var d interface{}
	c = a
	d = a

	fmt.Println("--->", d.(Error).Error())

	fmt.Println(reflect.TypeOf(a), reflect.TypeOf(b), c, d.(Error))
	fmt.Println(reflect.TypeOf(a) == reflect.TypeOf(new(Error)))

	x, ok := reflect.New(reflect.TypeOf(a)).Elem().Interface().(Error)
	fmt.Println(x, ok)

	fmt.Println(reflect.TypeOf(x))

	funcType := reflect.TypeOf(f)

	var e Error
	fmt.Println(funcType.In(0), reflect.TypeOf(e))

	fmt.Println(funcType.In(0).Kind())

	fmt.Println("1 >>>", reflect.TypeOf(Error(a)))

	is := reflect.TypeOf(a).Implements(funcType.In(0))
	fmt.Println("2 >>>", is)
}
