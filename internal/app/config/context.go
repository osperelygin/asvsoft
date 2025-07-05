package config

import "context"

type configContextKeyType struct {
}

var configContextKey = configContextKeyType{}

// FromContext - returns config from root context
func FromContext(ctx context.Context) *ModuleConfig {
	cfgRaw := ctx.Value(configContextKey)
	cfg, ok := cfgRaw.(*ModuleConfig)

	if ok {
		return cfg
	}

	return nil
}

// WrapContext - wraps config into context, so it can be passed through the application
func WrapContext(ctx context.Context, cfg *ModuleConfig) context.Context {
	return context.WithValue(ctx, configContextKey, cfg)
}
