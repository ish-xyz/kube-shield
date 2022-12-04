package cache

import (
	"sync"

	"k8s.io/client-go/tools/cache"
)

type Namespace string

type Group string

type Version string

type Kind string

type PolicyName string

type CacheIndex struct {
	sync.Mutex
	Policies map[Namespace]map[Group]map[Version]map[Kind][]PolicyName
}

type Controller struct {
	ClusterInformer   cache.SharedIndexInformer
	NamespaceInformer cache.SharedIndexInformer
	CacheIndex        *CacheIndex
}
