package qa

import (
	"reflect"
	"runtime"
	"sync"

	"github.com/reusee/dscope"
)

type CheckFunc func() []error

var _ dscope.CustomReducer = CheckFunc(nil)

func (_ CheckFunc) Reduce(_ dscope.Scope, vs []reflect.Value) reflect.Value {
	fn := CheckFunc(func() (ret []error) {
		var l sync.Mutex
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
				errs := fn()
				l.Lock()
				ret = append(ret, errs...)
				l.Unlock()
			}()
		}
		for i := 0; i < cap(sem); i++ {
			sem <- struct{}{}
		}
		return
	})
	return reflect.ValueOf(fn)
}

func (_ Def) CheckFunc() CheckFunc {
	return func() []error {
		return nil
	}
}
