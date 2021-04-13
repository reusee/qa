package qa

import (
	"io"
	"testing"

	"github.com/reusee/dscope"
)

func TestCheckFunc(t *testing.T) {
	ok := false
	dscope.New(
		dscope.Methods(new(Def))...,
	).Sub(func() Args {
		return []string{"."}
	}, func() CheckFunc {
		return func() []error {
			ok = true
			return []error{io.EOF}
		}
	}).Call(func(
		check CheckFunc,
	) {
		errs := check()
		if !ok {
			t.Fatal()
		}
		if len(errs) != 1 {
			t.Fatal()
		}
		if errs[0] != io.EOF {
			t.Fatal()
		}
	})
}
