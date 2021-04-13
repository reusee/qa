package qa

import (
	"reflect"
	"runtime"

	"github.com/reusee/dscope"
)

type CheckFunc func()

var _ dscope.Reducer = CheckFunc(nil)

func (_ CheckFunc) Reduce(_ dscope.Scope, vs []reflect.Value) reflect.Value {
	fn := func() {
		sem := make(chan struct{}, runtime.NumCPU())
		for _, v := range vs {
			fn := v.Interface().(CheckFunc)
			if fn == nil {
				continue
			}
			sem <- struct{}{}
			go func() {
				defer func() {
					<-sem
				}()
				fn()
			}()
		}
		for i := 0; i < cap(sem); i++ {
			sem <- struct{}{}
		}
	}
	return reflect.ValueOf(CheckFunc(fn))
}

func (_ Def) CheckFunc() CheckFunc {
	return func() {}
}
