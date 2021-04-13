package qa

import (
	"fmt"

	"github.com/reusee/e4"
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
		e4.Throw(fmt.Errorf("package load error"))
	}
	return pkgs
}
