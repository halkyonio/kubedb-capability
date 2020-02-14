package plugin

import (
	"fmt"
	"github.com/appscode/go/encoding/json/types"
	v1beta12 "halkyon.io/api/capability/v1beta1"
	"halkyon.io/api/v1beta1"
	framework "halkyon.io/operator-framework"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	v12 "kmodules.xyz/offshoot-api/api/v1"
	kubedbv1 "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
	"strings"
)

type KubeDBStatusAware interface {
	framework.DependentResource
	GetDatabasePhase(underlying runtime.Object) kubedbv1.DatabasePhase
	GetReason(underlying runtime.Object) string
}

func GetCondition(dep KubeDBStatusAware, err error, underlying runtime.Object) *v1beta1.DependentCondition {
	return framework.DefaultCustomizedGetConditionFor(dep, err, underlying, func(underlying runtime.Object, cond *v1beta1.DependentCondition) {
		ready := dep.GetDatabasePhase(underlying) == kubedbv1.DatabasePhaseRunning
		if !ready {
			msg := ""
			reason := dep.GetReason(underlying)
			if len(reason) > 0 {
				msg = ": " + reason
			}
			cond.Type = v1beta1.DependentPending
			cond.Message = fmt.Sprintf("%s is not ready%s", dep.Name(), msg)
		} else {
			cond.Type = v1beta1.DependentReady
			cond.Message = fmt.Sprintf("%s is ready", dep.Name())
		}
	})
}

func OwnerAsCapability(res framework.DependentResource) *v1beta12.Capability {
	return res.Owner().(*v1beta12.Capability)
}

// Convert Array of parameters to a Map
func ParametersAsMap(parameters []v1beta1.NameValuePair) map[string]string {
	result := make(map[string]string)
	for _, parameter := range parameters {
		result[parameter.Name] = parameter.Value
	}
	return result
}

func GetVersionFrom(capability *v1beta12.Capability, versionsMapping map[string]string) types.StrYo {
	version, ok := versionsMapping[capability.Spec.Version]
	if !ok {
		version = "Unknown or deprecated version: " + capability.Spec.Version
	}
	return types.StrYo(version)
}

func GetSecretOrDefault(needsSecret NeedsSecret, parameters map[string]string) *v1.SecretVolumeSource {
	if secretName, ok := parameters[DbConfigName]; ok {
		return &v1.SecretVolumeSource{SecretName: secretName}
	} else {
		// generate default secret name
		return &v1.SecretVolumeSource{SecretName: needsSecret.GetSecretName()}
	}
}

func GetDatabaseNameConfigOrNil(envVarName string, parameters map[string]string) *v12.PodTemplateSpec {
	if dbName, ok := parameters[DbName]; ok {
		return &v12.PodTemplateSpec{
			Spec: v12.PodSpec{
				Env: []v1.EnvVar{
					{Name: envVarName, Value: dbName},
				},
			},
		}
	}
	return nil
}

func DefaultSecretNameFor(secretOwner NeedsSecret) string {
	c := secretOwner.Owner().(*v1beta12.Capability)
	paramsMap := ParametersAsMap(c.Spec.Parameters)
	return SetDefaultSecretNameIfEmpty(c.Name, paramsMap[DbConfigName])
}

func SetDefaultSecretNameIfEmpty(capabilityName, paramSecretName string) string {
	if paramSecretName == "" {
		return strings.ToLower(capabilityName) + "-config"
	} else {
		return paramSecretName
	}
}

func SetDefaultDatabaseName(paramDatabaseName string) string {
	if paramDatabaseName == "" {
		return "sample-db"
	} else {
		return paramDatabaseName
	}
}

func SetDefaultDatabaseHost(capabilityHost, paramHost string) string {
	if paramHost == "" {
		return capabilityHost
	} else {
		return paramHost
	}
}

func SetDefaultDatabasePort(paramPort string) string {
	// TODO. Assign port according to the DB type using Enum
	if paramPort == "" {
		return "5432"
	} else {
		return paramPort
	}
}

//getAppLabels returns an string map with the labels which wil be associated to the kubernetes/ocp resource which will be created and managed by this operator
func GetAppLabels(name string) map[string]string {
	return map[string]string{
		"app": name,
	}
}

func ReplicaNumber(num int) *int32 {
	q := int32(num)
	return &q
}
