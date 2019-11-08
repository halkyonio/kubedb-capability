package main

import (
	framework "halkyon.io/operator-framework"
)

type role struct {
	framework.Role
}

func newRole(owner framework.Resource) role {
	generic := framework.NewOwnedRole(owner, func() string { return "scc-privileged-role" })
	r := role{Role: generic}
	generic.SetDelegate(r)
	return r
}
