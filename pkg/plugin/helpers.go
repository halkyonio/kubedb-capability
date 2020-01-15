package plugin

import (
	"github.com/appscode/go/encoding/json/types"
	v1beta12 "halkyon.io/api/capability/v1beta1"
	"halkyon.io/api/v1beta1"
	framework "halkyon.io/operator-framework"
	v1 "k8s.io/api/core/v1"
	v12 "kmodules.xyz/offshoot-api/api/v1"
	"strings"
)

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

func GetVersionFrom(capability *v1beta12.Capability) types.StrYo {
	return types.StrYo(capability.Spec.Version)
}

func GetSecretOrNil(parameters map[string]string) *v1.SecretVolumeSource {
	if secretName, ok := parameters[DbConfigName]; ok {
		return &v1.SecretVolumeSource{SecretName: secretName}
	}
	return nil
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
