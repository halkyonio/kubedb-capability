package main

import (
	"halkyon.io/kubedb-capability/pkg/plugin/mongodb"
	"halkyon.io/kubedb-capability/pkg/plugin/mysql"
	"halkyon.io/kubedb-capability/pkg/plugin/postgresql"
	plugins "halkyon.io/plugins/capability"
)

func main() {
	plugins.StartPluginServerFor(postgresql.NewPluginResource(), mysql.NewPluginResource(), mongodb.NewPluginResource())
}
