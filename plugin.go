package main

import (
	"encoding/gob"
	"github.com/hashicorp/go-plugin"
	halkyon "halkyon.io/api/capability/v1beta1"
	"halkyon.io/api/v1beta1"
	"halkyon.io/operator-framework"
	plugins "halkyon.io/plugins/capability"
	v12 "k8s.io/api/core/v1"
	kubedbv1 "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
	"os"
	"path/filepath"
)

var _ plugins.PluginResource = &PostgresPluginResource{}

var (
	postgresGVK = kubedbv1.SchemeGroupVersion.WithKind(kubedbv1.ResourceKindPostgres)
	secretGVK   = v12.SchemeGroupVersion.WithKind("Secret")
)

type PostgresPluginResource struct {
	ct halkyon.CapabilityType
	cc halkyon.CapabilityCategory
}

func (p *PostgresPluginResource) GetDependentResourcesWith(owner v1beta1.HalkyonResource) []framework.DependentResource {
	return []framework.DependentResource{
		framework.NewOwnedRole(owner, roleNamer),
		newRoleBinding(owner),
		newSecret(owner),
		newPostgres(owner),
	}
}

func roleNamer() string {
	return "scc-privileged-role"
}

func (p *PostgresPluginResource) GetSupportedCategory() halkyon.CapabilityCategory {
	return p.cc
}

func (p *PostgresPluginResource) GetSupportedType() halkyon.CapabilityType {
	return p.ct
}

func main() {
	gob.Register(kubedbv1.Postgres{})
	gob.Register(kubedbv1.PostgresList{})
	pluginName := filepath.Base(os.Args[0])
	p := &PostgresPluginResource{
		ct: halkyon.PostgresType,
		cc: halkyon.DatabaseCategory,
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugins.Handshake,
		Plugins:         map[string]plugin.Plugin{pluginName: &plugins.GoPluginPlugin{Delegate: p}},
	})
}
