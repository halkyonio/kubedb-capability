package main

import (
	"halkyon.io/api/v1beta1"
	framework "halkyon.io/operator-framework"
)

func PostgresName(owner v1beta1.HalkyonResource) string {
	return framework.DefaultDependentResourceNameFor(owner)
}

func ServiceAccountName(owner v1beta1.HalkyonResource) string {
	return PostgresName(owner) // todo: fix me
}
