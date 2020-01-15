package postgresql

import (
	"halkyon.io/api/v1beta1"
	"halkyon.io/kubedb-capability/pkg/plugin"
	framework "halkyon.io/operator-framework"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var secretGVK = v1.SchemeGroupVersion.WithKind("Secret")

type secret struct {
	*framework.BaseDependentResource
}

func (res secret) NameFrom(underlying runtime.Object) string {
	return framework.DefaultNameFrom(res, underlying)
}

func (res secret) Fetch() (runtime.Object, error) {
	return framework.DefaultFetcher(res)
}

func (res secret) IsReady(underlying runtime.Object) (ready bool, message string) {
	return framework.DefaultIsReady(underlying)
}

var _ framework.DependentResource = &secret{}

func (res secret) Update(_ runtime.Object) (bool, error) {
	return false, nil
}

func NewSecret(owner v1beta1.HalkyonResource) secret {
	config := framework.NewConfig(secretGVK)
	config.Watched = false
	return secret{framework.NewConfiguredBaseDependentResource(owner, config)}
}

//buildSecret returns the secret resource
func (res secret) Build(empty bool) (runtime.Object, error) {
	secret := &v1.Secret{}
	if !empty {
		c := plugin.OwnerAsCapability(res)
		ls := plugin.GetAppLabels(c.Name)
		paramsMap := plugin.ParametersAsMap(c.Spec.Parameters)
		secret.ObjectMeta = metav1.ObjectMeta{
			Name:      res.Name(),
			Namespace: c.Namespace,
			Labels:    ls,
		}
		secret.Data = map[string][]byte{
			KubedbPgUser:         []byte(paramsMap[plugin.DbUser]),
			KubedbPgPassword:     []byte(paramsMap[plugin.DbPassword]),
			KubedbPgDatabaseName: []byte(plugin.SetDefaultDatabaseName(paramsMap[plugin.DbName])),
			// TODO : To be reviewed according to the discussion started with issue #75
			// as we will create another secret when a link will be issued
			plugin.DbHost:     []byte(plugin.SetDefaultDatabaseHost(c.Name, paramsMap[plugin.DbHost])),
			plugin.DbPort:     []byte(plugin.SetDefaultDatabasePort(paramsMap[plugin.DbPort])),
			plugin.DbName:     []byte(plugin.SetDefaultDatabaseName(paramsMap[plugin.DbName])),
			plugin.DbUser:     []byte((paramsMap[plugin.DbUser])),
			plugin.DbPassword: []byte(paramsMap[plugin.DbPassword]),
		}
	}

	return secret, nil
}

func (res secret) Name() string {
	c := plugin.OwnerAsCapability(res)
	paramsMap := plugin.ParametersAsMap(c.Spec.Parameters)
	return plugin.SetDefaultSecretNameIfEmpty(c.Name, paramsMap[plugin.DbConfigName])
}
