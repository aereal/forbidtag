package main

import (
	"github.com/aereal/forbidtag"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(forbidtag.Analyzer)
}
