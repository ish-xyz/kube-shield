package loader

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

type ClusterPolicyLoader struct {
	Informer cache.SharedIndexInformer
}

func NewClusterPolicyLoader(clientset dynamic.Interface) *ClusterPolicyLoader {
	resource := schema.GroupVersionResource{Group: "kube-shield.red-labs.co.uk", Version: "v1", Resource: "Policy"}
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(clientset, time.Minute, metav1.NamespaceAll, nil)
	informer := factory.ForResource(resource).Informer()
	clusterPolicyLoader := &ClusterPolicyLoader{
		Informer: informer,
	}
	return clusterPolicyLoader
}

func (l ClusterPolicyLoader) RegisterHandler(handler cache.ResourceEventHandler) {
	l.Informer.AddEventHandler(handler)
}

func (l ClusterPolicyLoader) Run(stopCh <-chan struct{}) {
	l.Informer.Run(stopCh)
}
