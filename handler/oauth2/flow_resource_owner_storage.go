// Copyright © 2025 Ory Corp
// Modified by contributors to the Pocket ID fork.
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"
)

type ResourceOwnerPasswordCredentialsGrantStorage interface {
	Authenticate(ctx context.Context, name string, secret string) (subject string, err error)
	AccessTokenStorageProvider
	RefreshTokenStorageProvider
}
