package plugin

import (
	"halkyon.io/api/v1beta1"
	framework "halkyon.io/operator-framework"
)

func RoleName() string {
	return "scc-privileged-role"
}

func PostgresName(owner v1beta1.HalkyonResource) string {
	return framework.DefaultDependentResourceNameFor(owner)
}

func ServiceAccountName(owner v1beta1.HalkyonResource) string {
	return PostgresName(owner) // todo: fix me
}
