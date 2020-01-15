package plugin

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	client "kubedb.dev/apimachinery/client/clientset/versioned"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

var (
	NotDeprecated = v1.ListOptions{FieldSelector: fields.OneTermNotEqualSelector("Deprecated", "true").String()}
	Client        = client.NewForConfigOrDie(controllerruntime.GetConfigOrDie()).CatalogV1alpha1()
)
