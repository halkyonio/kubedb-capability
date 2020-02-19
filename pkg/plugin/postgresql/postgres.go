package postgresql

import (
	"halkyon.io/api/v1beta1"
	"halkyon.io/kubedb-capability/pkg/plugin"
	framework "halkyon.io/operator-framework"
	"halkyon.io/operator-framework/util"
	apps "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubedbv1 "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
)

const (
	// KubeDB Postgres const
	dbNameVarName     = "POSTGRES_DB"
	dbUserVarName     = "POSTGRES_USER"
	dbPasswordVarName = "POSTGRES_PASSWORD"
)

var (
	postgresGVK = kubedbv1.SchemeGroupVersion.WithKind(kubedbv1.ResourceKindPostgres)
)

type postgres struct {
	*framework.BaseDependentResource
}

func (res postgres) Fetch() (runtime.Object, error) {
	panic("should never be called")
}

var _ framework.DependentResource = &postgres{}

func (res postgres) Update(toUpdate runtime.Object) (bool, runtime.Object, error) {
	return false, toUpdate, nil
}

func NewPostgres(owner v1beta1.HalkyonResource) *postgres {
	config := framework.NewConfig(postgresGVK)
	config.CheckedForReadiness = true
	p := &postgres{framework.NewConfiguredBaseDependentResource(owner, config)}
	return p
}

func (res postgres) Name() string {
	return framework.DefaultDependentResourceNameFor(res.Owner())
}

//buildSecret returns the postgres resource
func (res *postgres) Build(empty bool) (runtime.Object, error) {
	postgres := &kubedbv1.Postgres{}
	if !empty {
		c := plugin.OwnerAsCapability(res)
		ls := plugin.GetAppLabels(c.Name)
		postgres.ObjectMeta = metav1.ObjectMeta{
			Name:      res.Name(),
			Namespace: c.Namespace,
			Labels:    ls,
		}
		postgres.Spec = kubedbv1.PostgresSpec{
			Version:  plugin.GetVersionFrom(c, versionsMapping),
			Replicas: plugin.ReplicaNumber(1),
			UpdateStrategy: apps.StatefulSetUpdateStrategy{
				Type: apps.RollingUpdateStatefulSetStrategyType,
			},
			StorageType:       kubedbv1.StorageTypeEphemeral,
			TerminationPolicy: kubedbv1.TerminationPolicyDelete,
		}

		paramsMap := util.ParametersAsMap(c.Spec.Parameters)
		if secret := plugin.GetSecretOrDefault(res, paramsMap); secret != nil {
			postgres.Spec.DatabaseSecret = secret
		}
		if dbNameConfig := plugin.GetDatabaseNameConfigOrNil(dbNameVarName, paramsMap); dbNameConfig != nil {
			postgres.Spec.PodTemplate = *dbNameConfig
		}
	}
	return postgres, nil
}

func (res *postgres) GetDatabasePhase(underlying runtime.Object) kubedbv1.DatabasePhase {
	return statusOf(underlying).Phase
}

func statusOf(underlying runtime.Object) kubedbv1.PostgresStatus {
	return underlying.(*kubedbv1.Postgres).Status
}

func (res *postgres) GetReason(underlying runtime.Object) string {
	return statusOf(underlying).Reason
}

func (res *postgres) GetCondition(underlying runtime.Object, err error) *v1beta1.DependentCondition {
	condition := plugin.GetCondition(res, err, underlying)
	return condition
}

func (res postgres) NameFrom(underlying runtime.Object) string {
	return underlying.(*kubedbv1.Postgres).Name
}

func (res postgres) GetRoleBindingName() string {
	return "use-scc-privileged"
}

func (res postgres) GetAssociatedRoleName() string {
	return res.GetRoleName()
}

func (res postgres) GetServiceAccountName() string {
	return res.Name()
}

func (res postgres) GetRoleName() string {
	return "scc-privileged-role"
}

func (res *postgres) GetDataMap() map[string][]byte {
	c := plugin.OwnerAsCapability(res)
	paramsMap := util.ParametersAsMap(c.Spec.Parameters)
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

func (res *postgres) GetSecretName() string {
	return plugin.DefaultSecretNameFor(res)
}
