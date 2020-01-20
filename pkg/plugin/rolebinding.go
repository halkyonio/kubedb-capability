package plugin

import (
	framework "halkyon.io/operator-framework"
	authorizv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type roleBinding struct {
	framework.RoleBinding
}

func NewRoleBinding(owner framework.NeedsRoleBinding) roleBinding {
	generic := framework.NewOwnedRoleBinding(owner)
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
		ser.Subjects = append(ser.Subjects, authorizv1.Subject{Kind: "ServiceAccount", Name: res.Delegate.GetServiceAccountName(), Namespace: owner.GetNamespace()})
	}

	return ser, nil
}
