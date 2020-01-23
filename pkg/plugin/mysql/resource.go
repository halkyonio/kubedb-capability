package mysql

import (
	"fmt"
	"github.com/appscode/go/strings"
	"github.com/hashicorp/go-hclog"
	"halkyon.io/api/capability/v1beta1"
	beta1 "halkyon.io/api/v1beta1"
	"halkyon.io/kubedb-capability/pkg/plugin"
	"halkyon.io/operator-framework"
	"halkyon.io/operator-framework/plugins/capability"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubedbv1 "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
)

var _ capability.PluginResource = &MySQLPluginResource{}

type MySQLPluginResource struct {
	capability.QueryingSimplePluginResourceStem
}

func (m MySQLPluginResource) GetDependentResourcesWith(owner beta1.HalkyonResource) []framework.DependentResource {
	mySQL := NewMySQL(owner)
	return []framework.DependentResource{framework.NewOwnedRole(mySQL),
		plugin.NewRoleBinding(mySQL),
		plugin.NewSecret(mySQL),
		mySQL}
}

func NewPluginResource() capability.PluginResource {
	return &MySQLPluginResource{capability.NewQueryingSimplePluginResourceStem(v1beta1.DatabaseCategory, resolver)}
}

func resolver(logger hclog.Logger) capability.TypeInfo {
	list, err := plugin.Client.MySQLVersions().List(v1.ListOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("error retrieving versions: %v", err))
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
	return info
}
