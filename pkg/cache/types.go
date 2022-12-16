package cache

import (
	"sync"

	"k8s.io/client-go/tools/cache"
)

type Verb string

type Namespace string

type Group string

type Resource string

type PolicyName string

type CacheIndex struct {
	sync.Mutex
	Policies map[Verb]map[Namespace]map[Group]map[Resource][]PolicyName
}

type Controller struct {
	ClusterInformer   cache.SharedIndexInformer
	NamespaceInformer cache.SharedIndexInformer
	CacheIndex        *CacheIndex
}
