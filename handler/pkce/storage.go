// Copyright © 2025 Ory Corp
// Modified by contributors to the Pocket ID fork.
// SPDX-License-Identifier: Apache-2.0

package pkce

import (
	"context"

	"github.com/ory/fosite"
)

type (
	PKCERequestStorage interface {
		GetPKCERequestSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error)
		CreatePKCERequestSession(ctx context.Context, signature string, requester fosite.Requester) error
		DeletePKCERequestSession(ctx context.Context, signature string) error
	}
	PKCERequestStorageProvider interface {
		PKCERequestStorage() PKCERequestStorage
	}
)
