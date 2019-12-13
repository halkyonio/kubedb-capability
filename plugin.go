package main

import (
	"github.com/hashicorp/go-plugin"
	halkyon "halkyon.io/api/capability/v1beta1"
	plugins "halkyon.io/plugins/capability"
	"k8s.io/apimachinery/pkg/runtime"
	kubedbv1 "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
	"os"
	"path/filepath"
)

var _ plugins.PluginResource = &PostgresPluginResource{}

type PostgresPluginResource struct {
	*postgres
	ct halkyon.CapabilityType
	cc halkyon.CapabilityCategory
}

func (p *PostgresPluginResource) Init(scheme *runtime.Scheme) {
	kubedbv1.AddToScheme(scheme)
}

func (p *PostgresPluginResource) GetSupportedCategory() halkyon.CapabilityCategory {
	return p.cc
}

func (p *PostgresPluginResource) GetSupportedType() halkyon.CapabilityType {
	return p.ct
}

func main() {
	pluginName := filepath.Base(os.Args[0])
	p := &PostgresPluginResource{
		postgres: newPostgres(nil),
		ct:       halkyon.PostgresType,
		cc:       halkyon.DatabaseCategory,
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugins.Handshake,
		Plugins:         map[string]plugin.Plugin{pluginName: &plugins.GoPluginPlugin{Delegate: p}},
	})
}
