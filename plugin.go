package main

import (
	"github.com/natefinch/pie"
	halkyon "halkyon.io/api/capability/v1beta1"
	plugins "halkyon.io/plugins/capability"
	"log"
	"net/rpc/jsonrpc"
)

func main() {
	p := pie.NewProvider()
	if err := p.RegisterName("postgresql-capability", plugins.NewPluginServer(halkyon.DatabaseCategory, halkyon.PostgresType, newPostgres(nil))); err != nil {
		log.Fatalf("failed to register Plugin: %s", err)
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
}
