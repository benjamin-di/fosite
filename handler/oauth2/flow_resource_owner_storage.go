// Copyright © 2025 Ory Corp
// Modified by contributors to the Pocket ID fork.
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"
)

// ResourceOwnerPasswordCredentialsGrantStorage provides storage for the resource owner password credentials grant.
type ResourceOwnerPasswordCredentialsGrantStorage interface {
	Authenticate(ctx context.Context, name string, secret string) (subject string, err error)
}

// ResourceOwnerPasswordCredentialsGrantStorageProvider provides the resource owner password credentials grant storage.
type ResourceOwnerPasswordCredentialsGrantStorageProvider interface {
	ResourceOwnerPasswordCredentialsGrantStorage() ResourceOwnerPasswordCredentialsGrantStorage
}
