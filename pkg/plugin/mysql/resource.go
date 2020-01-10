package mysql

import (
	"halkyon.io/api/capability/v1beta1"
	beta1 "halkyon.io/api/v1beta1"
	"halkyon.io/operator-framework"
	"halkyon.io/plugins/capability"
	kubedbv1 "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
)

var _ capability.PluginResource = &MySQLPluginResource{}

type MySQLPluginResource struct {
	capability.SimplePluginResourceStem
}

func (m MySQLPluginResource) GetDependentResourcesWith(owner beta1.HalkyonResource) []framework.DependentResource {
	return []framework.DependentResource{NewMySQL(owner)}
}

func NewPluginResource() capability.PluginResource {
	return &MySQLPluginResource{capability.NewSimplePluginResourceStem(v1beta1.DatabaseCategory, kubedbv1.ResourceKindMySQL)}
}
