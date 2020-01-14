package mysql

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

var _ framework.DependentResource = &mysql{}
var mysqlGVK = kubedbv1.SchemeGroupVersion.WithKind(kubedbv1.ResourceKindMySQL)

type mysql struct {
	*framework.BaseDependentResource
}

func NewMySQL(owner v1beta1.HalkyonResource) *mysql {
	config := framework.NewConfig(mysqlGVK)
	config.CheckedForReadiness = true
	config.OwnerStatusField = "PodName" // todo: find a way to compute this as above instead of hardcoding it
	p := &mysql{framework.NewConfiguredBaseDependentResource(owner, config)}
	return p
}

func (m *mysql) Name() string {
	return framework.DefaultDependentResourceNameFor(m.Owner())
}

func (m *mysql) NameFrom(underlying runtime.Object) string {
	return framework.DefaultNameFrom(m, underlying)
}

func (m *mysql) Fetch() (runtime.Object, error) {
	panic("should never be called")
}

func (m *mysql) Build(empty bool) (runtime.Object, error) {
	mysql := &kubedbv1.MySQL{}
	if !empty {
		c := plugin.OwnerAsCapability(m)
		ls := plugin.GetAppLabels(c.Name)
		mysql.ObjectMeta = metav1.ObjectMeta{
			Name:      m.Name(),
			Namespace: c.Namespace,
			Labels:    ls,
		}
		mysql.Spec = kubedbv1.MySQLSpec{
			Version:  "8.0-v2",
			Replicas: plugin.ReplicaNumber(1),
			UpdateStrategy: apps.StatefulSetUpdateStrategy{
				Type: apps.RollingUpdateStatefulSetStrategyType,
			},
			StorageType:       kubedbv1.StorageTypeEphemeral,
			TerminationPolicy: kubedbv1.TerminationPolicyDelete,
		}
	}
	return mysql, nil
}

func (m *mysql) Update(_ runtime.Object) (bool, error) {
	return false, nil
}

func (m *mysql) IsReady(underlying runtime.Object) (ready bool, message string) {
	mySQL := underlying.(*kubedbv1.MySQL)
	ready = mySQL.Status.Phase == kubedbv1.DatabasePhaseRunning
	if !ready {
		msg := ""
		reason := mySQL.Status.Reason
		if len(reason) > 0 {
			msg = ": " + reason
		}
		message = fmt.Sprintf("%s is not ready%s", mySQL.Name, msg)
	}
	return
}
