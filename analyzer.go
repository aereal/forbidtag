package forbidtag

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"strings"

	"github.com/aereal/forbidtag/internal/sets"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	flagAllow  = "allow"
	flagForbid = "forbid"
)

var ErrTagListConfliction = errors.New("cannot specify both -allow and -forbid flags")

var Analyzer = NewAnalyzer()

func NewAnalyzer() *analysis.Analyzer {
	a := &analysis.Analyzer{
		Name: "forbidtag",
		Doc:  "forbidtag finds unexpected struct tags",
		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
		},
		Run: run,
	}
	allowedTags := &strSet{}
	forbiddenTags := &strSet{}
	a.Flags.Var(allowedTags, flagAllow, "comma-separated list of allowed struct tag keys")
	a.Flags.Var(forbiddenTags, flagForbid, "comma-separated list of forbidden struct tag keys")
	return a
}

func run(pass *analysis.Pass) (any, error) {
	allowedTags, _ := lookupSetFlag(pass, flagAllow)
	if allowedTags == nil {
		allowedTags = &strSet{}
	}
	forbiddenTags, _ := lookupSetFlag(pass, flagForbid)
	if forbiddenTags == nil {
		forbiddenTags = &strSet{}
	}

	if !allowedTags.raw().IsEmpty() && !forbiddenTags.raw().IsEmpty() {
		return nil, ErrTagListConfliction
	}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector) //nolint:errcheck

	nodeFilter := []ast.Node{
		(*ast.StructType)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		structType := n.(*ast.StructType) //nolint:errcheck
		for _, field := range structType.Fields.List {
			if field.Tag != nil {
				checkStructTag(pass, field.Tag, allowedTags, forbiddenTags)
			}
		}
	})

	return nil, nil
}

func checkStructTag(pass *analysis.Pass, tagNode *ast.BasicLit, allowedTags, forbiddenTags *strSet) {
	if tagNode.Value == "" {
		return
	}

	tagValue := strings.Trim(tagNode.Value, "`")

	tagKeys := extractTagKeys(tagValue)
	cl := &checklist{allowed: allowedTags.raw(), forbidden: forbiddenTags.raw()}

	for _, key := range tagKeys {
		if cl.shouldReport(key) {
			pass.Reportf(tagNode.Pos(), "unexpected struct tag: %s", key)
		}
	}
}

func extractTagKeys(tagValue string) []string {
	var keys []string

	parts := strings.Fields(tagValue)
	for _, part := range parts {
		colonIdx := strings.Index(part, ":")
		if colonIdx > 0 {
			key := part[:colonIdx]
			keys = append(keys, key)
		}
	}

	return keys
}

type checklist struct {
	allowed, forbidden *sets.Set[string]
}

func (cl *checklist) shouldReport(key string) bool {
	if !cl.allowed.IsEmpty() {
		if !cl.allowed.Contains(key) {
			return true
		}
	} else if !cl.forbidden.IsEmpty() {
		if cl.forbidden.Contains(key) {
			return true
		}
	}
	return false
}

type strSet sets.Set[string]

var (
	_ flag.Getter = (*strSet)(nil)
)

func (l *strSet) Get() any { return l }

func (l *strSet) String() string {
	buf := new(bytes.Buffer)
	var seenFirst bool
	for el := range (*sets.Set[string])(l).Values() {
		if seenFirst {
			fmt.Fprint(buf, ",")
		}
		fmt.Fprint(buf, el)
		if !seenFirst {
			seenFirst = true
		}
	}
	return buf.String()
}

func (l *strSet) Set(v string) error {
	if l == nil {
		*l = strSet{} //nolint:govet
	}
	if v == "" {
		return nil
	}
	(*sets.Set[string])(l).Add(v)
	return nil
}

func (l *strSet) raw() *sets.Set[string] { return (*sets.Set[string])(l) }

func lookupSetFlag(pass *analysis.Pass, name string) (*strSet, bool) {
	fv := pass.Analyzer.Flags.Lookup(name)
	if fv == nil {
		return nil, false
	}
	tl, ok := fv.Value.(*strSet)
	return tl, ok
}
