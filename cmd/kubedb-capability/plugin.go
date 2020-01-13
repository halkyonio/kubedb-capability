package main

import (
	"encoding/gob"
	"github.com/hashicorp/go-plugin"
	"halkyon.io/kubedb-capability/pkg/plugin/mysql"
	"halkyon.io/kubedb-capability/pkg/plugin/postgresql"
	plugins "halkyon.io/plugins/capability"
	kubedbv1 "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
	"os"
	"path/filepath"
)

func main() {
	gob.Register(kubedbv1.Postgres{})
	gob.Register(kubedbv1.PostgresList{})
	pluginName := filepath.Base(os.Args[0])
	p, err := plugins.NewAggregatePluginResource(postgresql.NewPluginResource(), mysql.NewPluginResource())
	if err != nil {
		panic(err)
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugins.Handshake,
		Plugins:         map[string]plugin.Plugin{pluginName: &plugins.GoPluginPlugin{Delegate: p}},
	})
}
