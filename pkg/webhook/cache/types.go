package cache

import "k8s.io/client-go/tools/cache"

type Namespace string

type Group string

type Version string

type Kind string

type RuleName string

type CacheIndex struct {
	ClusterPolicies map[Group]map[Version]map[Kind][]RuleName
	Policies        map[Namespace]map[Group]map[Version]map[Kind][]RuleName
}

type CacheController struct {
	ClusterInformer   cache.SharedIndexInformer
	NamespaceInformer cache.SharedIndexInformer
	CacheIndex        *CacheIndex
}
