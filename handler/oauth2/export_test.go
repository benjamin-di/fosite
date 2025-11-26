// Copyright © 2025 Ory Corp
// Modified by contributors to the Pocket ID fork.
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"time"

	"github.com/ory/fosite"
)

func CallGetExpiresIn(r fosite.Requester, key fosite.TokenType, defaultLifespan time.Duration, now time.Time) time.Duration {
	return getExpiresIn(r, key, defaultLifespan, now)
}

func CallSignature(token string) string {
	return signature(token)
}
