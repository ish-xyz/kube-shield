package loader

import (
	"time"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

type ClusterPolicyLoader struct {
	Informer cache.SharedIndexInformer
}

func NewClusterPolicyLoader(clientset dynamic.Interface) *ClusterPolicyLoader {
	resource := schema.GroupVersionResource{Group: "kube-shield.red-labs.co.uk", Version: "v1", Resource: "ClusterPolicy"}
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(clientset, time.Minute, metav1.NamespaceAll, nil)
	informer := factory.ForResource(resource).Informer()
	clusterPolicyLoader := &ClusterPolicyLoader{
		Informer: informer,
	}
	return clusterPolicyLoader
}

func (l *ClusterPolicyLoader) List() []*unstructured.Unstructured {
	objList := l.Informer.GetStore().List()
	cplList := make([]*unstructured.Unstructured, len(objList))
	for idx, obj := range objList {
		unstr, exists := obj.(*unstructured.Unstructured)
		if !exists {
			logrus.Fatalln("Error converting object to unstructured")
		}
		cplList[idx] = unstr
	}
	return cplList
}
