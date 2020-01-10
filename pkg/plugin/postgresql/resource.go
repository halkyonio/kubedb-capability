package postgresql

import (
	"halkyon.io/api/capability/v1beta1"
	v1beta12 "halkyon.io/api/v1beta1"
	"halkyon.io/operator-framework"
	"halkyon.io/plugins/capability"
)

var _ capability.PluginResource = &PostgresPluginResource{}

func NewPluginResource() capability.PluginResource {
	return &PostgresPluginResource{capability.NewSimplePluginResourceStem(v1beta1.DatabaseCategory, v1beta1.PostgresType)}
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
