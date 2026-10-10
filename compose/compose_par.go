// Copyright © 2025 Ory Corp
// Modified by contributors to the Pocket ID fork.
// SPDX-License-Identifier: Apache-2.0

package compose

import (
	"github.com/pocket-id/fosite"
	"github.com/pocket-id/fosite/handler/par"
)

// PushedAuthorizeHandlerFactory creates the basic PAR handler
func PushedAuthorizeHandlerFactory(config fosite.Configurator, storage fosite.Storage, _ interface{}) interface{} {
	return &par.PushedAuthorizeHandler{
		Storage: storage.(fosite.PARStorageProvider),
		Config:  config,
	}
}
