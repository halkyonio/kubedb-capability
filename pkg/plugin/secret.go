package plugin

import (
	framework "halkyon.io/operator-framework"
)

func NewSecret(owner framework.NeedsSecret) framework.Secret {
	config := framework.NewDefaultSecretConfig()
	config.Updated = false
	config.Watched = false
	return framework.NewSecret(owner, config)
}
