package cache

import (
	"time"

	"github.com/RedLabsPlatform/kube-shield/pkg/config/defaults"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	kcache "k8s.io/client-go/tools/cache"
)

func NewEmptyCacheIndex() *CacheIndex {
	return &CacheIndex{
		ClusterPolicies: make(map[Group]map[Version]map[Kind]RuleName),
		Policies:        make(map[Group]map[Version]map[Kind]RuleName),
	}
}

func NewCacheController(clientset dynamic.Interface, c *CacheIndex) *CacheController {

	clusterPolicy := schema.GroupVersionResource{
		Group:    defaults.CR_GROUP,
		Version:  defaults.CR_VERSION,
		Resource: defaults.POLICY_KIND,
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
		CacheIndex:        c,
	}
}

func (c *CacheController) Run(stopCh <-chan struct{}) {

	// Register handlers
	c.ClusterInformer.AddEventHandler(kcache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAdd,
		DeleteFunc: c.onDelete,
	})

	c.NamespaceInformer.AddEventHandler(kcache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAdd,
		DeleteFunc: c.onDelete,
	})

	//go c.ClusterInformer.Run(stopCh)
	go c.NamespaceInformer.Run(stopCh)

	<-stopCh
}

// Reconcile() -> reconciles cache with resources in the cluster (using informers)
/*
	should list every cluster policy and every policy in each namespace and add it to the index
*/
