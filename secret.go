package main

import (
	"halkyon.io/api/v1beta1"
	framework "halkyon.io/operator-framework"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type secret struct {
	*framework.BaseDependentResource
}

func (res secret) NameFrom(underlying runtime.Object) string {
	return framework.DefaultNameFrom(res, underlying)
}

func (res secret) Fetch(helper *framework.K8SHelper) (runtime.Object, error) {
	return framework.DefaultFetcher(res, helper)
}

func (res secret) IsReady(underlying runtime.Object) (ready bool, message string) {
	return framework.DefaultIsReady(underlying)
}

var _ framework.DependentResource = &secret{}

func (res secret) Update(toUpdate runtime.Object) (bool, error) {
	return false, nil
}

func newSecret(owner v1beta1.HalkyonResource) secret {
	config := framework.NewConfig(secretGVK, owner.GetNamespace())
	config.Watched = false
	return secret{framework.NewConfiguredBaseDependentResource(owner, config)}
}

//buildSecret returns the secret resource
func (res secret) Build(empty bool) (runtime.Object, error) {
	secret := &v1.Secret{}
	if !empty {
		c := ownerAsCapability(res)
		ls := getAppLabels(c.Name)
		paramsMap := parametersAsMap(c.Spec.Parameters)
		secret.ObjectMeta = metav1.ObjectMeta{
			Name:      res.Name(),
			Namespace: c.Namespace,
			Labels:    ls,
		}
		secret.Data = map[string][]byte{
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
		}
	}

	return secret, nil
}

func (res secret) Name() string {
	c := ownerAsCapability(res)
	paramsMap := parametersAsMap(c.Spec.Parameters)
	return SetDefaultSecretNameIfEmpty(c.Name, paramsMap[DbConfigName])
}
