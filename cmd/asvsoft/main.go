package main

import (
	"asvsoft/internal/app/cli"
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
	err := cli.RootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
