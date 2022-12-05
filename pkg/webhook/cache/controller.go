package cache

import (
	"time"

	"github.com/RedLabsPlatform/kube-shield/pkg/config/defaults"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	kcache "k8s.io/client-go/tools/cache"
)

func NewEmptyCacheIndex() *CacheIndex {

	return &CacheIndex{
		Policies: make(map[Namespace]map[Group]map[Version]map[Kind][]PolicyName),
	}
}

func NewCacheController(clientset dynamic.Interface, c *CacheIndex) *Controller {

	clusterPolicy := schema.GroupVersionResource{
		Group:    defaults.CR_GROUP,
		Version:  defaults.CR_VERSION,
		Resource: defaults.CLUSTER_POLICY_KIND,
	}

	policy := schema.GroupVersionResource{
		Group:    defaults.CR_GROUP,
		Version:  defaults.CR_VERSION,
		Resource: defaults.POLICY_KIND,
	}

	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(clientset, time.Minute, metav1.NamespaceAll, nil)

	return &Controller{
		ClusterInformer:   factory.ForResource(clusterPolicy).Informer(),
		NamespaceInformer: factory.ForResource(policy).Informer(),
		CacheIndex:        c,
	}
}

// Registers handlers for ClusterPolicies and Policies informers and start them
func (c *Controller) Run(polStopCh <-chan struct{}, clusterPolStopCh <-chan struct{}) {

	// Register handlers
	c.ClusterInformer.AddEventHandler(kcache.ResourceEventHandlerFuncs{
		AddFunc:    c.onClusterPolicyAdd,
		UpdateFunc: c.onClusterPolicyUpdate,
		DeleteFunc: c.onClusterPolicyDelete,
	})

	c.NamespaceInformer.AddEventHandler(kcache.ResourceEventHandlerFuncs{
		AddFunc:    c.onPolicyAdd,
		UpdateFunc: c.onPolicyUpdate,
		DeleteFunc: c.onPolicyDelete,
	})

	go c.ClusterInformer.Run(clusterPolStopCh)
	go c.NamespaceInformer.Run(polStopCh)

	for {
		select {
		case err := <-clusterPolStopCh:
			logrus.Fatal("cluster policy informer failed %v", err)
		case err := <-polStopCh:
			logrus.Fatal("namespace policy informer failed %v", err)
		}
	}
}
