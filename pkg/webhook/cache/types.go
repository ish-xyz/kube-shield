package cache

import "k8s.io/client-go/tools/cache"

type Group string

type Version string

type Kind string

type RuleName string

type Cache struct {
	ClusterPolicies map[Group]map[Version]map[Kind]RuleName
	Policies        map[Group]map[Version]map[Kind]RuleName
}

type CacheController struct {
	ClusterInformer   cache.SharedIndexInformer
	NamespaceInformer cache.SharedIndexInformer
}
