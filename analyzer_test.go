package forbidtag_test

import (
	"testing"

	"github.com/aereal/forbidtag"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), forbidtag.NewAnalyzer(), "a")
}

func TestAnalyzerForbidMode(t *testing.T) {
	analyzer := forbidtag.NewAnalyzer()
	_ = analyzer.Flags.Set("allow", "")
	_ = analyzer.Flags.Set("forbid", "db")
	_ = analyzer.Flags.Set("forbid", "yaml")
	analysistest.Run(t, analysistest.TestData(), analyzer, "forbid_mode")
}

func TestAnalyzerAllowMode(t *testing.T) {
	analyzer := forbidtag.NewAnalyzer()
	_ = analyzer.Flags.Set("allow", "json")
	_ = analyzer.Flags.Set("forbid", "")
	analysistest.Run(t, analysistest.TestData(), analyzer, "allow_mode")
}
