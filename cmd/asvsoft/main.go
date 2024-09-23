// Package main ...
package main

import (
	"asvsoft/internal/app"
	"asvsoft/internal/app/ds"
	"os"
)

// ldflags
var (
	// BuildTime   - Время сборки
	BuildTime string
	// BuildCommit -  Коммит из которого был билд
	BuildCommit string
	// BuildBranch -  Ветка из которой был билд
	BuildBranch string
)

func main() {
	err := app.Init(&ds.AppInfo{
		BuildTime:   BuildTime,
		BuildCommit: BuildCommit,
		BuildBranch: BuildBranch,
	})
	if err != nil {
		os.Exit(1)
	}
}
