package oplint_test

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/gilcrest/oplint/oplint"
)

// This is the directory where our test fixtures are.
const fixtureDirName = "testdata"

var fixtureDir = filepath.Join("/Users/gilcrest/go_modules/oplint", "/", fixtureDirName)

func TestAll(t *testing.T) {
	analysistest.Run(t, fixtureDir, oplint.Analyzer, "")
}
