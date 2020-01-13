package postgresql

import (
	"fmt"
	"github.com/appscode/go/encoding/json/types"
	"halkyon.io/api/v1beta1"
	"halkyon.io/kubedb-capability/pkg/plugin"
	framework "halkyon.io/operator-framework"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ofst "kmodules.xyz/offshoot-api/api/v1"
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
		paramsMap := plugin.ParametersAsMap(c.Spec.Parameters)
		postgres.ObjectMeta = metav1.ObjectMeta{
			Name:      res.Name(),
			Namespace: c.Namespace,
			Labels:    ls,
		}
		postgres.Spec = kubedbv1.PostgresSpec{
			Version:  SetDefaultDatabaseVersionIfEmpty(c.Spec.Version),
			Replicas: plugin.ReplicaNumber(1),
			UpdateStrategy: apps.StatefulSetUpdateStrategy{
				Type: apps.RollingUpdateStatefulSetStrategyType,
			},
			DatabaseSecret: &core.SecretVolumeSource{
				SecretName: plugin.SetDefaultSecretNameIfEmpty(c.Name, paramsMap[DbConfigName]),
			},
			StorageType:       kubedbv1.StorageTypeEphemeral,
			TerminationPolicy: kubedbv1.TerminationPolicyDelete,
			PodTemplate: ofst.PodTemplateSpec{
				Spec: ofst.PodSpec{
					Env: []core.EnvVar{
						{Name: KubedbPgDatabaseName, Value: plugin.SetDefaultDatabaseName(paramsMap[DbName])},
					},
				},
			},
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

func SetDefaultDatabaseVersionIfEmpty(version string) types.StrYo {
	if version == "10.6-v2" {
		return types.StrYo("10.6")
	} else {
		// Map DB Version with the KubeDB Version
		switch version {
		case "9":
			return types.StrYo("9.6-v4")
		case "10":
			return types.StrYo("10.6-v2")
		case "11":
			return types.StrYo("11.2")
		default:
			return types.StrYo("10.6-v2")
		}
	}
}

/*
	// https://github.com/kubernetes/client-go/tree/master/examples/dynamic-create-update-delete-deployment
	// Approach to create dynamically tyhe object without type imported
	postgresRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	postgres := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
		},
	}
	// Create Postgres DB
		fmt.Println("Creating Postgres DB ...")
		result, err := client.Resource(postgresRes).Namespace(namespace).Create(postgres, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}
*/
