package main

import (
	"halkyon.io/kubedb-capability/pkg/plugin/mongodb"
	"halkyon.io/kubedb-capability/pkg/plugin/mysql"
	"halkyon.io/kubedb-capability/pkg/plugin/postgresql"
	plugins "halkyon.io/operator-framework/plugins/capability"
	"log"
	"os"
	"path/filepath"
)

func main() {
	pluginName := filepath.Base(os.Args[0])
	f, err := os.OpenFile(pluginName+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)

	plugins.StartPluginServerFor(postgresql.NewPluginResource(), mysql.NewPluginResource(), mongodb.NewPluginResource())
}
