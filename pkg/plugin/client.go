package plugin

import (
	client "kubedb.dev/apimachinery/client/clientset/versioned"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

var (
	Client = client.NewForConfigOrDie(controllerruntime.GetConfigOrDie()).CatalogV1alpha1()
)
