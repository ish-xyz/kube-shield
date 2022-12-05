package engine

import "github.com/RedLabsPlatform/kube-shield/pkg/webhook/cache"

type Engine struct {

	// Dependency injection for the Cache controller
	CacheController *cache.Controller

	// By default the webhook fails when an invalid check is encountered, this behaviour can be changed here
	SkipInvalidChecks bool
}
