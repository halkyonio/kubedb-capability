package postgresql

import (
	core "k8s.io/api/core/v1"
	kubedbv1 "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
)

const (
	Secret = "Secret"
	// KubeDB Postgres const
	KubedbPgDatabaseName = "POSTGRES_DB"
	KubedbPgUser         = "POSTGRES_USER"
	KubedbPgPassword     = "POSTGRES_PASSWORD"
	// Capability const
	DbConfigName = "DB_CONFIG_NAME"
	DbHost       = "DB_HOST"
	DbPort       = "DB_PORT"
	DbName       = "DB_NAME"
	DbUser       = "DB_USER"
	DbPassword   = "DB_PASSWORD"
)

var (
	postgresGVK = kubedbv1.SchemeGroupVersion.WithKind(kubedbv1.ResourceKindPostgres)
	secretGVK   = core.SchemeGroupVersion.WithKind(Secret)
)
