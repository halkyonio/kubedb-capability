package main

import (
	capability "halkyon.io/api/capability/v1beta1"
	framework "halkyon.io/operator-framework"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type secret struct {
	*framework.DependentResourceHelper
}

func (res secret) Update(toUpdate runtime.Object) (bool, error) {
	return false, nil
}

func newSecret(owner framework.Resource) secret {
	resource := framework.NewDependentResource(&v1.Secret{}, owner)
	s := secret{DependentResourceHelper: resource}
	resource.SetDelegate(s)
	return s
}

func (res secret) ownerAsCapability() *capability.Capability {
	return ownerAsCapability(res)
}

//buildSecret returns the secret resource
func (res secret) Build() (runtime.Object, error) {
	c := res.ownerAsCapability()
	ls := getAppLabels(c.Name)
	paramsMap := parametersAsMap(c.Spec.Parameters)
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      res.Name(),
			Namespace: c.Namespace,
			Labels:    ls,
		},
		Data: map[string][]byte{
			KubedbPgUser:         []byte(paramsMap[DbUser]),
			KubedbPgPassword:     []byte(paramsMap[DbPassword]),
			KubedbPgDatabaseName: []byte(SetDefaultDatabaseName(paramsMap[DbName])),
			// TODO : To be reviewed according to the discussion started with issue #75
			// as we will create another secret when a link will be issued
			DbHost:     []byte(SetDefaultDatabaseHost(c.Name, paramsMap[DbHost])),
			DbPort:     []byte(SetDefaultDatabasePort(paramsMap[DbPort])),
			DbName:     []byte(SetDefaultDatabaseName(paramsMap[DbName])),
			DbUser:     []byte((paramsMap[DbUser])),
			DbPassword: []byte(paramsMap[DbPassword]),
		},
	}

	return secret, nil
}

func (res secret) Name() string {
	c := res.ownerAsCapability()
	paramsMap := parametersAsMap(c.Spec.Parameters)
	return SetDefaultSecretNameIfEmpty(c.Name, paramsMap[DbConfigName])
}

func (res secret) ShouldWatch() bool {
	return false
}
