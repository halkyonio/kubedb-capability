package postgresql

import (
	"github.com/appscode/go/strings"
	"github.com/hashicorp/go-hclog"
	"halkyon.io/api/capability/v1beta1"
	v1beta12 "halkyon.io/api/v1beta1"
	"halkyon.io/kubedb-capability/pkg/plugin"
	"halkyon.io/operator-framework"
	"halkyon.io/operator-framework/plugins/capability"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubedbv1 "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
)

var _ capability.PluginResource = &PostgresPluginResource{}

func NewPluginResource() capability.PluginResource {
	return &PostgresPluginResource{capability.NewQueryingSimplePluginResourceStem(v1beta1.DatabaseCategory, resolver)}
}

func resolver(logger hclog.Logger) capability.TypeInfo {
	list, err := plugin.Client.PostgresVersions().List(v1.ListOptions{})
	if err != nil {
		logger.Error("error retrieving versions: %v", err)
	}
	versions := make([]string, 0, len(list.Items))
	for _, version := range list.Items {
		if !version.Spec.Deprecated && !strings.Contains(versions, version.Spec.Version) {
			versions = append(versions, version.Spec.Version)
		}
	}
	info := capability.TypeInfo{
		Type:     kubedbv1.ResourceKindPostgres,
		Versions: versions,
	}
	return info
}

type PostgresPluginResource struct {
	capability.QueryingSimplePluginResourceStem
}

func (p *PostgresPluginResource) GetDependentResourcesWith(owner v1beta12.HalkyonResource) []framework.DependentResource {
	postgres := NewPostgres(owner)
	return []framework.DependentResource{
		framework.NewOwnedRole(postgres),
		plugin.NewRoleBinding(postgres),
		plugin.NewSecret(postgres),
		postgres,
	}
}
