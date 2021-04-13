package qa

import (
	"testing"

	"github.com/reusee/dscope"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
)

func TestAnalysis(t *testing.T) {
	dscope.New(
		dscope.Methods(new(Def))...,
	).Sub(AnalyzersToDefs([]*analysis.Analyzer{
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		printf.Analyzer,
		sortslice.Analyzer,
		stringintconv.Analyzer,
		testinggoroutine.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
	})...).Sub(func() Args {
		return Args{"."}
	}).Call(func(
		check CheckFunc,
	) {
		errs := check()
		for _, err := range errs {
			t.Fatal(err)
		}
	})
}
