package main

import (
	"halkyon.io/api/v1beta1"
	framework "halkyon.io/operator-framework"
	authorizv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type roleBinding struct {
	framework.RoleBinding
}

func newRoleBinding(owner v1beta1.HalkyonResource) roleBinding {
	generic := framework.NewOwnedRoleBinding(owner,
		func() string { return "use-scc-privileged" },
		func() string { return roleNamer() },
		func() string { return newPostgres(owner).Name() })
	rb := roleBinding{RoleBinding: generic}
	return rb
}

func (res roleBinding) Build(empty bool) (runtime.Object, error) {
	build, err := res.RoleBinding.Build(empty)
	if err != nil {
		return nil, err
	}

	ser := build.(*authorizv1.RoleBinding)
	if !empty {
		owner := res.Owner()
		ser.Subjects = append(ser.Subjects, authorizv1.Subject{Kind: "ServiceAccount", Name: PostgresName(owner), Namespace: owner.GetNamespace()})
	}

	return ser, nil
}
