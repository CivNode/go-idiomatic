// Command go-idiomatic runs the pedagogical Go-idiom checks from the
// go-idiomatic module as a standalone linter driven by the standard
// golang.org/x/tools/go/analysis framework.
//
// Install:
//
//	go install github.com/CivNode/go-idiomatic/cmd/go-idiomatic@latest
//
// Run against a package:
//
//	go-idiomatic ./...
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/CivNode/go-idiomatic/rules"
)

func main() {
	singlechecker.Main(rules.Analyzer)
}
