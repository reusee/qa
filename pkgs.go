package qa

import (
	"fmt"

	"github.com/reusee/e5"
	"golang.org/x/tools/go/packages"
)

func (_ Def) Packages(
	args Args,
) []*packages.Package {
	// load
	pkgs, err := packages.Load(
		&packages.Config{
			Mode: 0xffffffff,
		},
		args...,
	)
	ce(err)
	if packages.PrintErrors(pkgs) > 0 {
		e5.Throw(fmt.Errorf("package load error"))
	}
	return pkgs
}
