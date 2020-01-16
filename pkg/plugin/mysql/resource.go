package mysql

import (
	"github.com/appscode/go/strings"
	"halkyon.io/api/capability/v1beta1"
	beta1 "halkyon.io/api/v1beta1"
	"halkyon.io/kubedb-capability/pkg/plugin"
	"halkyon.io/operator-framework"
	"halkyon.io/plugins/capability"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	list, err := plugin.Client.MySQLVersions().List(v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	versions := make([]string, 0, len(list.Items))
	for _, version := range list.Items {
		if !version.Spec.Deprecated && !strings.Contains(versions, version.Spec.Version) {
			versions = append(versions, version.Spec.Version)
		}
	}
	info := capability.TypeInfo{
		Type:     kubedbv1.ResourceKindMySQL,
		Versions: versions,
	}
	return &MySQLPluginResource{capability.NewSimplePluginResourceStem(v1beta1.DatabaseCategory, info)}
}
