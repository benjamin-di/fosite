# Pocket ID Fosite fork changes

This repository is the [Pocket ID](https://github.com/pocket-id/pocket-id) fork
of [Ory Fosite](https://github.com/ory/fosite).

The fork keeps custom changes as a linear series of commits on top of upstream
`master`. Use the following command to list the custom commits:

```sh
git log --oneline <upstream-master>..main
```

The module path remains `github.com/ory/fosite`. Pocket ID uses a `replace`
directive to select this fork without changing import paths.

## Commits on `main`

### 1. `chore: upgrade deps`

Updates the Go version and project dependencies.

### 2. `chore: update copyright`

Restores the copyright headers inherited from Ory, marks inherited files changed
by the Pocket ID fork, and attributes files added by the fork to Pocket ID. Adds
fork-aware license header tooling and runs it as part of `make format`.

### 3. `ci/cd: adapt actions for Pocket ID use`

Removes GitHub Actions workflows that depend on Ory infrastructure. These
workflows include the stale issue bot, label synchronization, the closed
reference notifier, and the OIDC conformance workflows. Updates `licenses.yml`
to remove automatic commits from the Ory bot.

### 4. `feat: replace ory with pocket-id prefix`

Changes the vendor prefixes for opaque tokens:

- Access tokens use `pocket_id_at_`.
- Refresh tokens use `pocket_id_rt_`.
- Authorization codes use `pocket_id_ac_`.
- RFC 8628 device codes use `pocket_id_dc_`.

### 5. `feat: validate OIDC prompt in authorize requests`

Validates the OIDC `prompt` parameter in `NewAuthorizeRequest` and pushed
authorization requests. For requests with the `openid` scope, Fosite rejects
unknown prompt values and rejects `none` when it is combined with another value.
The default allowed values are `login`, `none`, `consent`, and `select_account`.

Prompt validation occurs after redirect URI validation. Fosite redirects the
error when the redirect URI is valid.

### 6. `fix: redirect unsupported request object errors`

Validates the redirect URI before returning errors from early OIDC request
object processing. Fosite redirects these errors when the redirect URI is valid.
It returns a JSON error when the redirect URI is not valid.

This applies to errors such as unsupported `request` parameters, invalid
`request_uri` values, and malformed request objects.

### 7. `feat: add configurable redirect URI matcher`

Adds the `fosite.RedirectURIMatcher` function type, the
`RedirectURIMatcherProvider` interface, and the `Config.RedirectURIMatcher`
option. The default is `MatchRedirectURIWithClientRedirectURIs`.

Fosite stores the matcher on each `AuthorizeRequest`. It restores the matcher
when a pushed authorization request continues.

### 8. `feat: include wrong redirect uri in error message`

Adds the rejected `redirect_uri` value to the `invalid_request` hint for an
unregistered redirect URI. Pocket ID displays this hint on its error page.

### 9. `feat: add RFC 8707 resource indicator support`

Adds support for RFC 8707 resource indicators:

- Adds `fosite.ErrInvalidTarget` with HTTP status 400.
- Validates the `resource` parameter at the authorization, pushed authorization,
  token, and device authorization endpoints.
- Requires one absolute URI without a fragment or whitespace.
- Adds `Requester.GetResource()`.
- Adds the `IsValidScopeToken` helper for the RFC 6749 scope-token grammar.

Pocket ID resolves resource identifiers to audiences and scopes.

### 10. `feat: ignore the "audience" request parameter`

Stops parsing the non-standard `audience` request parameter at the
authorization, token, and device authorization endpoints. Removes the
`fosite.GetAudiences` helper and the `validateAudience` check.

Fosite continues to support audiences set through `SetRequestedAudience` and
`GrantAudience`. Audience matching also remains available for individual flows.

### 11. `feat: add option to ignore unknown scopes`

Adds `Config.IgnoreUnknownScopes` and the exported
`fosite.FilterRequestedScopes` helper. When the option is enabled, Fosite
removes unregistered scopes from authorization, pushed authorization, and device
authorization requests instead of returning `invalid_scope`.

The default is `false`. Token endpoint grants and the refresh token flow remain
strict.

### 12. `feat: support unsigned request objects for all clients`

Accepts unsigned request objects with `alg: none` from all clients. Signed
request objects still require an `OpenIDConnectClient` with a JSON Web Key Set.
The `request_uri` parameter still requires a registered request URI.

Adds `Config.SupportedRequestObjectSigningAlgorithms` and the
`SupportedRequestObjectSigningAlgorithmsProvider` interface. Fosite returns
`invalid_request_object` when a request object uses an algorithm that is not in
the configured list.

### 13. `feat: add OAuth Client ID Metadata Document support (#1)`

Adds parsing and validation for OAuth Client ID Metadata Documents. Adds a
client resolver that uses registered clients first and resolves eligible client
identifier URLs when no registered client exists.

Adds an HTTP fetcher with response size limits, timeouts, cache controls, and
SSRF protections. The fetcher blocks private and special-use IP addresses and
refuses redirects by default. The same guarded HTTP client loads referenced JSON
Web Key Sets and request objects.

Adds interfaces for client caching, materialization, and authorization server
policy. Authorization, introspection, and client authentication use the resolved
client.

### 14. `feat: add RFC 9068 JWTStrategy`

Adds `RFC9068JWTStrategy` for RFC 9068 access tokens. The strategy sets the JWT
type to `at+jwt`, sets the `client_id` claim from the authenticated client, and
sets the `scope` claim from granted scopes.

The strategy requires `iss`, `sub`, `aud`, `exp`, `iat`, `jti`, and `client_id`.
Refresh tokens and authorization codes remain opaque. Adds
`compose.NewOAuth2RFC9068JWTStrategy` to construct the strategy.

### 15. `feat: add RFC 9207 authorization response issuer`

Adds `Config.AuthorizationResponseIssuer` and the
`AuthorizationResponseIssuerProvider` interface. When configured, Fosite adds
the RFC 9207 `iss` parameter to successful authorization responses and
redirectable authorization error responses.

Fosite writes the parameter for the default, query, fragment, and `form_post`
response modes. An empty configuration value disables the parameter. The
configured issuer replaces an existing `iss` response parameter.
