package mongodb

import (
	"fmt"
	"github.com/appscode/go/strings"
	"github.com/hashicorp/go-hclog"
	"halkyon.io/api/capability/v1beta1"
	beta1 "halkyon.io/api/v1beta1"
	"halkyon.io/kubedb-capability/pkg/plugin"
	framework "halkyon.io/operator-framework"
	"halkyon.io/operator-framework/plugins/capability"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubedbv1 "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
)

var _ capability.PluginResource = &MongoDBPluginResource{}

type MongoDBPluginResource struct {
	capability.QueryingSimplePluginResourceStem
}

func (m MongoDBPluginResource) GetDependentResourcesWith(owner beta1.HalkyonResource) []framework.DependentResource {
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
		if !version.Spec.Deprecated && !strings.Contains(versions, version.Spec.Version) {
			versions = append(versions, version.Spec.Version)
		}
	}
	info := capability.TypeInfo{
		Type:     kubedbv1.ResourceKindMongoDB,
		Versions: versions,
	}
	return info
}
