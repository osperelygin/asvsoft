// Package app ...
package app

import (
	"asvsoft/internal/app/cli"
	"asvsoft/internal/app/ctxutils"
	"asvsoft/internal/app/ds"
	"context"
)

func Init(appinfo *ds.AppInfo) error {
	ctx := ctxutils.InitStorage(context.Background())
	ctxutils.SaveAppInfo(ctx, appinfo)

	err := cli.RootCmd().ExecuteContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
