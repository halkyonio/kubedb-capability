package postgresql

import (
	kubedbv1 "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
)

const (
	// KubeDB Postgres const
	KubedbPgDatabaseName = "POSTGRES_DB"
	KubedbPgUser         = "POSTGRES_USER"
	KubedbPgPassword     = "POSTGRES_PASSWORD"
)

var (
	postgresGVK = kubedbv1.SchemeGroupVersion.WithKind(kubedbv1.ResourceKindPostgres)
)
