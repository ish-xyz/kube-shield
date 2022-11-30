package cache

import (
	"time"

	"github.com/RedLabsPlatform/kube-shield/pkg/config/defaults"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

func NewCacheController(clientset dynamic.Interface) *CacheController {

	clusterPolicy := schema.GroupVersionResource{
		Group:    defaults.CR_GROUP,
		Version:  defaults.CR_VERSION,
		Resource: defaults.CLUSTER_POLICY_KIND,
	}

	policy := schema.GroupVersionResource{
		Group:    defaults.CR_GROUP,
		Version:  defaults.CR_VERSION,
		Resource: defaults.CLUSTER_POLICY_KIND,
	}

	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(clientset, time.Minute, metav1.NamespaceAll, nil)

	return &CacheController{
		ClusterInformer:   factory.ForResource(clusterPolicy).Informer(),
		NamespaceInformer: factory.ForResource(policy).Informer(),
	}
}

func NewEmptyCache() *Cache {
	return &Cache{
		ClusterPolicies: make(map[Group]map[Version]map[Kind]RuleName),
		Policies:        make(map[Group]map[Version]map[Kind]RuleName),
	}
}

func NewDynamicInformer(clientset dynamic.Interface, grp, ver, res string) cache.SharedIndexInformer {

	gvr := schema.GroupVersionResource{
		Group:    grp,
		Version:  ver,
		Resource: res,
	}
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(clientset, time.Minute, metav1.NamespaceAll, nil)

	return factory.ForResource(gvr).Informer()
}

// OnAdd()  -> add name to Cache
// OnDelete() -> remove from Cache
// Reconcile() -> reconciles cache with resources in the cluster (using informers)
