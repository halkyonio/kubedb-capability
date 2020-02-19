package postgresql

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

var _ capability.PluginResource = &PostgresPluginResource{}
var versionsMapping = make(map[string]string, 11)

func NewPluginResource() capability.PluginResource {
	return &PostgresPluginResource{capability.NewQueryingSimplePluginResourceStem(v1beta1.DatabaseCategory, resolver)}
}

func resolver(logger hclog.Logger) capability.TypeInfo {
	list, err := plugin.Client.PostgresVersions().List(v1.ListOptions{})
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
		Type:     kubedbv1.ResourceKindPostgres,
		Versions: versions,
	}
	return info
}

type PostgresPluginResource struct {
	capability.QueryingSimplePluginResourceStem
}

func (p *PostgresPluginResource) GetDependentResourcesWith(owner framework.SerializableResource) []framework.DependentResource {
	postgres := NewPostgres(owner)
	return []framework.DependentResource{
		framework.NewOwnedRole(postgres),
		plugin.NewRoleBinding(postgres),
		plugin.NewSecret(postgres),
		postgres,
	}
}
