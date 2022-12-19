package main

import (
	"github.com/gilcrest/oplint/oplint"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(oplint.Analyzer)
}
