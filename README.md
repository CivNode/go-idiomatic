# go-idiomatic

A pedagogical Go-idiom linter. Rules beyond `gopls` and `golangci-lint`, aimed at the layer of feedback you want when learning Go: "this compiles, but it isn't how Go is usually written." Part of the CivNode Training semantic engine.

Apache-2.0 licensed. No runtime dependencies outside `golang.org/x/tools/go/analysis`.

## Install

Library:

```
go get github.com/CivNode/go-idiomatic@latest
```

Standalone linter:

```
go install github.com/CivNode/go-idiomatic/cmd/go-idiomatic@latest
go-idiomatic ./...
```

## Library usage

```go
import goidiomatic "github.com/CivNode/go-idiomatic"

score, findings, err := goidiomatic.Score(src, goidiomatic.DefaultRules())
```

`Score` returns an integer on [0, 100] plus every finding. The score starts at 100 and deducts per severity: 15 for Error, 7 for Warn, 2 for Info, floored at 0.

## Rules (v0.1.0)

| ID | Severity | What it catches |
| --- | --- | --- |
| `prefer-range-int` | info | `for i := 0; i < len(xs); i++` where `for i := range xs` would do |
| `errors-is-as` | warn | `err.Error() == "..."` or `err == someErr`; use `errors.Is` / `errors.As` |
| `any-over-empty-interface` | info | `interface{}` on a Go 1.18+ target; suggests `any` |
| `context-first-arg` | warn | functions where `context.Context` is not the first argument |
| `no-sleep-for-coordination` | error | `time.Sleep` inside a function that also uses channels or `context.Context` |

Each rule ships with golden fixtures under `rules/testdata/src/<fixture>/`.

## Adding a rule

1. Drop a new file in `rules/`, implementing the `Rule` interface.
2. Add it to `rules.All`.
3. Add `testdata/src/<name>_ok/ok.go` and `testdata/src/<name>_bad/bad.go` fixtures, using `// want ...` comments on lines that must flag.
4. Write a test with `analysistest.Run` plus a few in-memory `rules.Run` unit tests for the edge cases.
5. `make lint test` must stay clean.

## Status

v0.1.0 ships the five rules above with >= 85% coverage on the rules package. Future tiers extend the ruleset and add fixers.
