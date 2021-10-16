package qa

import "github.com/reusee/dscope"

func Run(
	pkg string,
	defObjects ...any,
) (ret []error) {
	defs := dscope.Methods(new(Def))
	for _, obj := range defObjects {
		defs = append(defs, dscope.Methods(obj)...)
	}
	dscope.New(defs...).Fork(
		func() Args {
			return []string{pkg}
		},
	).Call(func(
		check CheckFunc,
	) {
		ret = check()
	})
	return
}
