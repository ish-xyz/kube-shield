package cache

import (
	"strings"
	"time"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	kcache "k8s.io/client-go/tools/cache"
)

func NewCacheController(clientset dynamic.Interface, c *CacheIndex) *Controller {

	clusterPolicy := schema.GroupVersionResource{
		Group:    v1.GroupVersion.Group,
		Version:  v1.GroupVersion.Version,
		Resource: v1.ClusterPolicyKind,
	}

	policy := schema.GroupVersionResource{
		Group:    v1.GroupVersion.Group,
		Version:  v1.GroupVersion.Version,
		Resource: v1.PolicyKind,
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
	// c.ClusterInformer.AddEventHandler(kcache.ResourceEventHandlerFuncs{
	// 	AddFunc:    c.onClusterPolicyAdd,
	// 	UpdateFunc: c.onClusterPolicyUpdate,
	// 	DeleteFunc: c.onClusterPolicyDelete,
	// })

	c.NamespaceInformer.AddEventHandler(kcache.ResourceEventHandlerFuncs{
		AddFunc:    c.onPolicyAdd,
		UpdateFunc: c.onPolicyUpdate,
		DeleteFunc: c.onPolicyDelete,
	})

	//go c.ClusterInformer.Run(clusterPolStopCh)
	c.NamespaceInformer.Run(polStopCh)

	//	for {
	//		select {
	//		case err := <-clusterPolStopCh:
	//			logrus.Fatalf("cluster policy informer failed %v", err)
	//		case err := <-polStopCh:
	//			logrus.Fatalf("namespace policy informer failed %v", err)
	//		}
	//	}
}

// Extract API Group & Version from apiVersion field
// the core api group, which is usually equal to "" (empty) is set as "_core" here
func (c *Controller) GetGV(info string) (Group, Version) {

	gv := strings.Split(info, "/")
	group := "_core"
	version := info
	if len(gv) > 1 {
		group = gv[0]
		version = gv[1]
	}

	return Group(group), Version(version)
}
