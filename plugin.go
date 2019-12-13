package main

import (
	"github.com/hashicorp/go-plugin"
	halkyon "halkyon.io/api/capability/v1beta1"
	plugins "halkyon.io/plugins/capability"
)

const pluginName = "postgresql-capability"

type PostgresPluginResource struct {
	*postgres
	ct halkyon.CapabilityType
	cc halkyon.CapabilityCategory
}

func (p *PostgresPluginResource) GetSupportedCategory() halkyon.CapabilityCategory {
	return p.cc
}

func (p *PostgresPluginResource) GetSupportedType() halkyon.CapabilityType {
	return p.ct
}

func main() {
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
