package mongodb

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"halkyon.io/api/capability/v1beta1"
	"halkyon.io/kubedb-capability/pkg/plugin"
	framework "halkyon.io/operator-framework"
	"halkyon.io/operator-framework/plugins/capability"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubedbv1 "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
	"strings"
)

var _ capability.PluginResource = &MongoDBPluginResource{}
var versionsMapping = make(map[string]string, 11)

type MongoDBPluginResource struct {
	capability.QueryingSimplePluginResourceStem
}

func (m MongoDBPluginResource) CheckValidity(owner framework.SerializableResource) []string {
	return []string{}
}

func (m MongoDBPluginResource) GetDependentResourcesWith(owner framework.SerializableResource) []framework.DependentResource {
	mongoDB := NewMongoDB(owner)
	return []framework.DependentResource{framework.NewOwnedRole(mongoDB),
		plugin.NewRoleBinding(mongoDB),
		plugin.NewSecret(mongoDB),
		mongoDB}
}

func NewPluginResource() capability.PluginResource {
	return &MongoDBPluginResource{capability.NewQueryingSimplePluginResourceStem(v1beta1.DatabaseCategory, resolver)}
}

func resolver(logger hclog.Logger) capability.TypeInfo {
	list, err := plugin.Client.MongoDBVersions().List(v1.ListOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("error retrieving versions: %v", err))
	}
	versions := make([]string, 0, len(list.Items))
	for _, version := range list.Items {
		if !version.Spec.Deprecated {
			external := version.Spec.Version
			internal, ok := versionsMapping[external]
			if !ok {
				versions = append(versions, external)
				versionsMapping[external] = version.Name
			} else {
				if strings.Compare(internal, version.Name) < 0 {
					versionsMapping[external] = version.Name
				}
			}
		}
	}
	info := capability.TypeInfo{
		Type:     kubedbv1.ResourceKindMongoDB,
		Versions: versions,
	}
	return info
}
