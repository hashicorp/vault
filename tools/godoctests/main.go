package main

import (
	"github.com/hashicorp/vault/tools/godoctests/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
