package loader

import (
	"context"

	"github.com/RedLabsPlatform/kube-shield/pkg/defaults"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

const (
	defaultNamespace = "default"
)

// Return a Kubernetes Client
func NewLoader(cfg *rest.Config, client *kubernetes.Clientset, dynClient dynamic.Interface) *Loader {

	return &Loader{
		KubeConfig:    cfg,
		KubeClient:    client,
		KubeDynClient: dynClient,
	}
}

// List ClusterPolicies (Cluster level resource)
func (l *Loader) GetClusterPolicies(ctx context.Context) (*unstructured.UnstructuredList, error) {
	return l.getPolicies(ctx, defaults.CLUSTER_POLICY_KIND, "")
}

// List Policies (Namespaced resource)
func (l *Loader) GetPolicies(ctx context.Context, namespace string) (*unstructured.UnstructuredList, error) {
	if namespace == "" {
		namespace = defaults.DEFAULT_NS
	}
	return l.getPolicies(ctx, defaults.POLICY_KIND, namespace)
}

// Private: list Policies and Cluster Policies
func (l *Loader) getPolicies(ctx context.Context, kind string, namespace string) (*unstructured.UnstructuredList, error) {

	var dr dynamic.ResourceInterface

	// Init discovery client and mapper
	dc, err := discovery.NewDiscoveryClientForConfig(l.KubeConfig)
	if err != nil {
		logrus.Debugln(err)
		return nil, err
	}

	// Get GVR
	groupResources, err := restmapper.GetAPIGroupResources(dc)
	if err != nil {
		logrus.Debugln(err)
		return nil, err
	}

	mapper := restmapper.NewDiscoveryRESTMapper(groupResources)
	mapping, err := mapper.RESTMapping(schema.GroupKind{Kind: kind, Group: defaults.API_GROUP})
	if err != nil {
		logrus.Debugln(err)
		return nil, err
	}

	// Init dynamic client
	dr = l.KubeDynClient.Resource(mapping.Resource)
	if namespace != "" {
		dr = l.KubeDynClient.Resource(mapping.Resource).Namespace(namespace)
	}

	retrievedObjects, err := dr.List(ctx, metav1.ListOptions{})
	if retrievedObjects == nil {
		retrievedObjects = &unstructured.UnstructuredList{
			Items: []unstructured.Unstructured{},
		}
	}

	logrus.Debugf("Number of policies retrieved %d", len(retrievedObjects.Items))

	return retrievedObjects, err
}
