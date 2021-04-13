package qa

import (
	"fmt"
	"go/types"
	"reflect"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

func AnalyzersToDefs(analyzers []*analysis.Analyzer) (defs []any) {

	type objectFactKey struct {
		obj types.Object
		t   reflect.Type
	}

	type packageFactKey struct {
		pkg *types.Package
		t   reflect.Type
	}

	return []any{
		func(
			pkgs []*packages.Package,
		) CheckFunc {
			return func() (retErrs []error) {

				doneAnalyzers := make(map[*analysis.Analyzer]bool)
				results := make(map[*types.Package]map[*analysis.Analyzer]any)
				objectFacts := make(map[objectFactKey]analysis.Fact)
				pkgFacts := make(map[packageFactKey]analysis.Fact)

				var runAnalyzer func(*analysis.Analyzer)
				runAnalyzer = func(analyzer *analysis.Analyzer) {
					if _, ok := doneAnalyzers[analyzer]; ok {
						return
					}
					doneAnalyzers[analyzer] = true

					for _, req := range analyzer.Requires {
						runAnalyzer(req)
					}

					for _, pkg := range pkgs {

						pass := &analysis.Pass{
							Analyzer: analyzer,

							Fset:         pkg.Fset,
							Files:        pkg.Syntax,
							OtherFiles:   pkg.OtherFiles,
							IgnoredFiles: pkg.IgnoredFiles,
							Pkg:          pkg.Types,
							TypesInfo:    pkg.TypesInfo,
							TypesSizes:   pkg.TypesSizes,

							Report: func(diagnostic analysis.Diagnostic) {
								retErrs = append(retErrs, fmt.Errorf(
									"%s: <%s> %s",
									pkg.Fset.Position(diagnostic.Pos).String(),
									diagnostic.Category,
									diagnostic.Message,
								))
							},

							ResultOf: results[pkg.Types],

							ImportObjectFact: func(obj types.Object, fact analysis.Fact) bool {
								t := reflect.TypeOf(fact)
								key := objectFactKey{
									obj: obj,
									t:   t,
								}
								v, ok := objectFacts[key]
								if !ok {
									return false
								}
								reflect.ValueOf(fact).Elem().Set(reflect.ValueOf(v).Elem())
								return true
							},

							ImportPackageFact: func(pkg *types.Package, fact analysis.Fact) bool {
								t := reflect.TypeOf(fact)
								key := packageFactKey{
									pkg: pkg,
									t:   t,
								}
								v, ok := pkgFacts[key]
								if !ok {
									return false
								}
								reflect.ValueOf(fact).Elem().Set(reflect.ValueOf(v).Elem())
								return true
							},

							ExportObjectFact: func(obj types.Object, fact analysis.Fact) {
								t := reflect.TypeOf(fact)
								key := objectFactKey{
									obj: obj,
									t:   t,
								}
								objectFacts[key] = fact
							},

							ExportPackageFact: func(fact analysis.Fact) {
								t := reflect.TypeOf(fact)
								key := packageFactKey{
									pkg: pkg.Types,
									t:   t,
								}
								pkgFacts[key] = fact
							},

							AllObjectFacts: func() (ret []analysis.ObjectFact) {
								for _, fact := range analyzer.FactTypes {
									t := reflect.TypeOf(fact)
									for key, v := range objectFacts {
										if key.t != t {
											continue
										}
										ret = append(ret, analysis.ObjectFact{
											Object: key.obj,
											Fact:   v,
										})
									}
								}
								return
							},

							AllPackageFacts: func() (ret []analysis.PackageFact) {
								for _, fact := range analyzer.FactTypes {
									t := reflect.TypeOf(fact)
									for key, v := range pkgFacts {
										if key.pkg != pkg.Types {
											continue
										}
										if reflect.TypeOf(v) != t {
											continue
										}
										ret = append(ret, analysis.PackageFact{
											Package: pkg.Types,
											Fact:    v,
										})
									}
								}
								return
							},
						}

						res, err := analyzer.Run(pass)
						ce(err)
						m, ok := results[pkg.Types]
						if !ok {
							m = make(map[*analysis.Analyzer]any)
							results[pkg.Types] = m
						}
						m[analyzer] = res
					}

				}

				for _, analyzer := range analyzers {
					runAnalyzer(analyzer)
				}

				return
			}
		},
	}
}
