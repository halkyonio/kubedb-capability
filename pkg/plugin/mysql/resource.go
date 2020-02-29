package mysql

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"halkyon.io/api/capability/v1beta1"
	"halkyon.io/kubedb-capability/pkg/plugin"
	"halkyon.io/operator-framework"
	"halkyon.io/operator-framework/plugins/capability"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubedbv1 "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
	"strings"
)

var _ capability.PluginResource = &MySQLPluginResource{}
var versionsMapping = make(map[string]string, 11)

type MySQLPluginResource struct {
	capability.QueryingSimplePluginResourceStem
}

func (m MySQLPluginResource) CheckValidity(owner framework.SerializableResource) []string {
	return []string{}
}

func (m MySQLPluginResource) GetDependentResourcesWith(owner framework.SerializableResource) []framework.DependentResource {
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
		if !version.Spec.Deprecated {
			external := version.Spec.Version
			internal, ok := versionsMapping[external]
			if !ok {
				versions = append(versions, external)
				versionsMapping[external] = version.Name
			} else {
				if strings.Compare(internal, version.Name) < 0 {
					versionsMapping[external] = version.Name
				}
			}
		}
	}
	info := capability.TypeInfo{
		Type:     kubedbv1.ResourceKindMySQL,
		Versions: versions,
	}
	return info
}
