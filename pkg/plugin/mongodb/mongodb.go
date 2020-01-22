package mongodb

import (
	"halkyon.io/api/v1beta1"
	"halkyon.io/kubedb-capability/pkg/plugin"
	framework "halkyon.io/operator-framework"
	apps "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubedbv1 "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
)

const (
	// see https://kubedb.com/docs/v0.13.0-rc.0/concepts/databases/mongodb/#spec-podtemplate-spec-env for db name var name
	dbNameVarName     = "MONGO_INITDB_DATABASE"
	dbUserVarName     = "MONGO_INITDB_ROOT_USERNAME"
	dbPasswordVarName = "MONGO_INITDB_ROOT_PASSWORD"
)

var _ framework.DependentResource = &mongodb{}
var gvk = kubedbv1.SchemeGroupVersion.WithKind(kubedbv1.ResourceKindMongoDB)

type mongodb struct {
	*framework.BaseDependentResource
}

func NewMongoDB(owner v1beta1.HalkyonResource) *mongodb {
	config := framework.NewConfig(gvk)
	config.CheckedForReadiness = true
	config.OwnerStatusField = "PodName" // todo: find a way to compute this as above instead of hardcoding it
	p := &mongodb{framework.NewConfiguredBaseDependentResource(owner, config)}
	return p
}

func (m *mongodb) Name() string {
	return framework.DefaultDependentResourceNameFor(m.Owner())
}

func (m *mongodb) NameFrom(underlying runtime.Object) string {
	return framework.DefaultNameFrom(m, underlying)
}

func (m *mongodb) Fetch() (runtime.Object, error) {
	panic("should never be called")
}

func (m *mongodb) Build(empty bool) (runtime.Object, error) {
	mongo := &kubedbv1.MongoDB{}
	if !empty {
		c := plugin.OwnerAsCapability(m)
		ls := plugin.GetAppLabels(c.Name)
		mongo.ObjectMeta = metav1.ObjectMeta{
			Name:      m.Name(),
			Namespace: c.Namespace,
			Labels:    ls,
		}
		mongo.Spec = kubedbv1.MongoDBSpec{
			Version:  plugin.GetVersionFrom(c),
			Replicas: plugin.ReplicaNumber(1),
			UpdateStrategy: apps.StatefulSetUpdateStrategy{
				Type: apps.RollingUpdateStatefulSetStrategyType,
			},
			StorageType:       kubedbv1.StorageTypeEphemeral,
			TerminationPolicy: kubedbv1.TerminationPolicyDelete,
		}

		paramsMap := plugin.ParametersAsMap(c.Spec.Parameters)
		if secret := plugin.GetSecretOrDefault(m, paramsMap); secret != nil {
			mongo.Spec.DatabaseSecret = secret
		}
		if dbNameConfig := plugin.GetDatabaseNameConfigOrNil("MONGO_INITDB_DATABASE", paramsMap); dbNameConfig != nil {
			mongo.Spec.PodTemplate = dbNameConfig
		}
	}
	return mongo, nil
}

func (m *mongodb) Update(toUpdate runtime.Object) (bool, error) {
	return false, nil
}

func (m *mongodb) IsReady(underlying runtime.Object) (ready bool, message string) {
	panic("implement me")
}

func (m *mongodb) GetRoleBindingName() string {
	return "use-scc-privileged"
}

func (m *mongodb) GetAssociatedRoleName() string {
	return m.GetRoleName()
}

func (m *mongodb) GetServiceAccountName() string {
	return m.Name()
}

func (m *mongodb) GetRoleName() string {
	return "scc-privileged-role"
}

func (m *mongodb) GetDataMap() map[string][]byte {
	c := plugin.OwnerAsCapability(m)
	paramsMap := plugin.ParametersAsMap(c.Spec.Parameters)
	return map[string][]byte{
		dbUserVarName:     []byte(paramsMap[plugin.DbUser]),
		dbPasswordVarName: []byte(paramsMap[plugin.DbPassword]),
		dbNameVarName:     []byte(plugin.SetDefaultDatabaseName(paramsMap[plugin.DbName])),
		// TODO : To be reviewed according to the discussion started with issue #75
		// as we will create another secret when a link will be issued
		plugin.DbHost:     []byte(plugin.SetDefaultDatabaseHost(c.Name, paramsMap[plugin.DbHost])),
		plugin.DbPort:     []byte(plugin.SetDefaultDatabasePort(paramsMap[plugin.DbPort])),
		plugin.DbName:     []byte(plugin.SetDefaultDatabaseName(paramsMap[plugin.DbName])),
		plugin.DbUser:     []byte((paramsMap[plugin.DbUser])),
		plugin.DbPassword: []byte(paramsMap[plugin.DbPassword]),
	}
}

func (m *mongodb) GetSecretName() string {
	return plugin.DefaultSecretNameFor(m)
}
