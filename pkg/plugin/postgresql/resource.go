package postgresql

import (
	"github.com/appscode/go/strings"
	"halkyon.io/api/capability/v1beta1"
	v1beta12 "halkyon.io/api/v1beta1"
	"halkyon.io/kubedb-capability/pkg/plugin"
	"halkyon.io/operator-framework"
	"halkyon.io/plugins/capability"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubedbv1 "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
)

var _ capability.PluginResource = &PostgresPluginResource{}

func NewPluginResource() capability.PluginResource {
	list, err := plugin.Client.PostgresVersions().List(v1.ListOptions{})
	versions := []string{}
	if err == nil {
		versions = make([]string, 0, len(list.Items))
		for _, version := range list.Items {
			if !version.Spec.Deprecated && !strings.Contains(versions, version.Spec.Version) {
				versions = append(versions, version.Spec.Version)
			}
		}
	}
	info := capability.TypeInfo{
		Type:     kubedbv1.ResourceKindPostgres,
		Versions: versions,
	}
	return &PostgresPluginResource{capability.NewSimplePluginResourceStem(v1beta1.DatabaseCategory, info)}
}

type PostgresPluginResource struct {
	capability.SimplePluginResourceStem
}

func (p *PostgresPluginResource) GetDependentResourcesWith(owner v1beta12.HalkyonResource) []framework.DependentResource {
	return []framework.DependentResource{
		framework.NewOwnedRole(owner, RoleName),
		NewRoleBinding(owner),
		NewSecret(owner),
		NewPostgres(owner),
	}
}
