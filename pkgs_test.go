package qa

import (
	"testing"

	"github.com/reusee/dscope"
	"golang.org/x/tools/go/packages"
)

func TestPkgs(t *testing.T) {
	dscope.New(
		dscope.Methods(new(Def))...,
	).Fork(func() Args {
		return []string{"."}
	}).Call(func(
		pkgs []*packages.Package,
	) {
		if len(pkgs) == 0 {
			t.Fatal()
		}
	})
}
