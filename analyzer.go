package forbidtag

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var Analyzer = &analysis.Analyzer{
	Name: "forbidtag",
	Doc:  "forbidtag finds unexpected struct tags",
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	Run: run,
}

func run(pass *analysis.Pass) (any, error) {
	panic("TODO")
}
