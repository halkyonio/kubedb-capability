package main

import (
	framework "halkyon.io/operator-framework"
)

func PostgresName(owner framework.Resource) string {
	return framework.DefaultDependentResourceNameFor(owner)
}

func ServiceAccountName(owner framework.Resource) string {
	return PostgresName(owner) // todo: fix me
}
