package main

import (
	"k8s.io/apimachinery/pkg/runtime"
	kubedbv1 "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
)

var Initializer = initializer{}

type initializer struct{}

func (initializer) Init(scheme *runtime.Scheme) {
	scheme.AddKnownTypes(kubedbv1.SchemeGroupVersion, &kubedbv1.Postgres{}, &kubedbv1.PostgresList{})
}
