package informer

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

func NewDynamicInformer(clientset dynamic.Interface, gvr schema.GroupVersionResource) cache.SharedIndexInformer {
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(clientset, time.Minute, metav1.NamespaceAll, nil)

	return factory.ForResource(gvr).Informer()
}
