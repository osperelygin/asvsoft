// Package app ...
package app

import (
	"asvsoft/internal/app/cli"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	// BuildTime   - Время сборки
	BuildTime string
	// BuildCommit -  Коммит из которого был билд
	BuildCommit string
	// BuildBranch -  Ветка из которой был билд
	BuildBranch string
}

func Init(cfg Config) error {
	log.Infof(
		"BuildTime: %s, BuildCommit: %s, BuildBranch: %s",
		cfg.BuildTime, cfg.BuildCommit, cfg.BuildBranch,
	)

	err := cli.RootCmd().Execute()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
