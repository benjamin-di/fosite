// Copyright © 2025 Ory Corp
// Modified by contributors to the Pocket ID fork.
// SPDX-License-Identifier: Apache-2.0

package openid

import (
	"context"
	"time"

	"github.com/ory/fosite"
)

func CallGenerateIDToken(ctx context.Context, lifespan time.Duration, fosr fosite.Requester, h *IDTokenHandleHelper) (token string, err error) {
	return h.generateIDToken(ctx, lifespan, fosr)
}
