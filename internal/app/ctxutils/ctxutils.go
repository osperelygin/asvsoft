// Package ctxutils ...
package ctxutils

import (
	"asvsoft/internal/app/ds"
	sensehat "asvsoft/internal/app/sensors/sense-hat"
	serialport "asvsoft/internal/pkg/serial-port"
	"context"
)

// ----------------------------
// -- Storage
// ----------------------------

type key int

const storageKey key = iota

type Storage struct {
	AppInfo *ds.AppInfo
	SrcCfg  *serialport.Config
	DstCfg  *serialport.Config
	ImuCfg  *sensehat.ImuConfig
}

func InitStorage(parent context.Context) context.Context {
	return context.WithValue(parent, storageKey, &Storage{})
}

func GetStorage(ctx context.Context) *Storage {
	return ctx.Value(storageKey).(*Storage)
}

// ----------------------------
// -- AppInfo
// ----------------------------

func SaveAppInfo(parent context.Context, appInfo *ds.AppInfo) {
	GetStorage(parent).AppInfo = appInfo
}

func GetAppInfo(parent context.Context) *ds.AppInfo {
	return GetStorage(parent).AppInfo
}
