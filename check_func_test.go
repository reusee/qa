package qa

import (
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
		return func() {
			ok = true
		}
	}).Call(func(
		check CheckFunc,
	) {
		check()
		if !ok {
			t.Fatal()
		}
	})
}
