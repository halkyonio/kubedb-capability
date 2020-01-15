package postgresql

import (
	"fmt"
	"halkyon.io/api/v1beta1"
	"halkyon.io/kubedb-capability/pkg/plugin"
	framework "halkyon.io/operator-framework"
	apps "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubedbv1 "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
)

type postgres struct {
	*framework.BaseDependentResource
}

func (res postgres) Fetch() (runtime.Object, error) {
	panic("should never be called")
}

var _ framework.DependentResource = &postgres{}

func (res postgres) Update(_ runtime.Object) (bool, error) {
	return false, nil
}

func NewPostgres(owner v1beta1.HalkyonResource) *postgres {
	config := framework.NewConfig(postgresGVK)
	config.CheckedForReadiness = true
	config.OwnerStatusField = "PodName" // todo: find a way to compute this as above instead of hardcoding it
	p := &postgres{framework.NewConfiguredBaseDependentResource(owner, config)}
	return p
}

func (res postgres) Name() string {
	return PostgresName(res.Owner())
}

//buildSecret returns the postgres resource
func (res postgres) Build(empty bool) (runtime.Object, error) {
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
			Version:  plugin.GetVersionFrom(c),
			Replicas: plugin.ReplicaNumber(1),
			UpdateStrategy: apps.StatefulSetUpdateStrategy{
				Type: apps.RollingUpdateStatefulSetStrategyType,
			},
			StorageType:       kubedbv1.StorageTypeEphemeral,
			TerminationPolicy: kubedbv1.TerminationPolicyDelete,
		}

		paramsMap := plugin.ParametersAsMap(c.Spec.Parameters)
		if secret := plugin.GetSecretOrNil(paramsMap); secret != nil {
			postgres.Spec.DatabaseSecret = secret
		}
		if dbNameConfig := plugin.GetDatabaseNameConfigOrNil(KubedbPgDatabaseName, paramsMap); dbNameConfig != nil {
			postgres.Spec.PodTemplate = *dbNameConfig
		}
	}
	return postgres, nil
}

func (res postgres) IsReady(underlying runtime.Object) (ready bool, message string) {
	psql := underlying.(*kubedbv1.Postgres)
	ready = psql.Status.Phase == kubedbv1.DatabasePhaseRunning
	if !ready {
		msg := ""
		reason := psql.Status.Reason
		if len(reason) > 0 {
			msg = ": " + reason
		}
		message = fmt.Sprintf("%s is not ready%s", psql.Name, msg)
	}
	return
}

func (res postgres) NameFrom(underlying runtime.Object) string {
	return underlying.(*kubedbv1.Postgres).Name
}
