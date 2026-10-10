// Copyright © 2025 Ory Corp
// Modified by contributors to the Pocket ID fork.
// SPDX-License-Identifier: Apache-2.0

package fosite

//go:generate go tool mockgen -package internal -destination internal/access_request.go github.com/pocket-id/fosite AccessRequester
//go:generate go tool mockgen -package internal -destination internal/access_response.go github.com/pocket-id/fosite AccessResponder
//go:generate go tool mockgen -package internal -destination internal/access_token_storage.go github.com/pocket-id/fosite/handler/oauth2 AccessTokenStorage
//go:generate go tool mockgen -package internal -destination internal/access_token_storage_provider.go github.com/pocket-id/fosite/handler/oauth2 AccessTokenStorageProvider
//go:generate go tool mockgen -package internal -destination internal/access_token_strategy.go github.com/pocket-id/fosite/handler/oauth2 AccessTokenStrategy
//go:generate go tool mockgen -package internal -destination internal/access_token_strategy_provider.go github.com/pocket-id/fosite/handler/oauth2 AccessTokenStrategyProvider
//go:generate go tool mockgen -package internal -destination internal/authorize_code_storage.go github.com/pocket-id/fosite/handler/oauth2 AuthorizeCodeStorage
//go:generate go tool mockgen -package internal -destination internal/authorize_code_storage_provider.go github.com/pocket-id/fosite/handler/oauth2 AuthorizeCodeStorageProvider
//go:generate go tool mockgen -package internal -destination internal/authorize_code_strategy.go github.com/pocket-id/fosite/handler/oauth2 AuthorizeCodeStrategy
//go:generate go tool mockgen -package internal -destination internal/authorize_code_strategy_provider.go github.com/pocket-id/fosite/handler/oauth2 AuthorizeCodeStrategyProvider
//go:generate go tool mockgen -package internal -destination internal/authorize_endpoint_handler.go github.com/pocket-id/fosite AuthorizeEndpointHandler
//go:generate go tool mockgen -package internal -destination internal/authorize_endpoint_handlers_provider.go github.com/pocket-id/fosite AuthorizeEndpointHandlersProvider
//go:generate go tool mockgen -package internal -destination internal/authorize_request.go github.com/pocket-id/fosite AuthorizeRequester
//go:generate go tool mockgen -package internal -destination internal/authorize_response.go github.com/pocket-id/fosite AuthorizeResponder
//go:generate go tool mockgen -package internal -destination internal/client.go github.com/pocket-id/fosite Client
//go:generate go tool mockgen -package internal -destination internal/client_manager.go github.com/pocket-id/fosite ClientManager
//go:generate go tool mockgen -package internal -destination internal/oauth2_strategy.go github.com/pocket-id/fosite/handler/oauth2 CoreStrategy
//go:generate go tool mockgen -package internal -destination internal/device_auth_storage.go github.com/pocket-id/fosite/handler/rfc8628 DeviceAuthStorage
//go:generate go tool mockgen -package internal -destination internal/device_auth_storage_provider.go github.com/pocket-id/fosite/handler/rfc8628 DeviceAuthStorageProvider
//go:generate go tool mockgen -package internal -destination internal/device_code_strategy.go github.com/pocket-id/fosite/handler/rfc8628 DeviceCodeStrategy
//go:generate go tool mockgen -package internal -destination internal/device_code_strategy_provider.go github.com/pocket-id/fosite/handler/rfc8628 DeviceCodeStrategyProvider
//go:generate go tool mockgen -package internal -destination internal/device_rate_limit_strategy.go github.com/pocket-id/fosite/handler/rfc8628 DeviceRateLimitStrategy
//go:generate go tool mockgen -package internal -destination internal/device_rate_limit_strategy_provider.go github.com/pocket-id/fosite/handler/rfc8628 DeviceRateLimitStrategyProvider
//go:generate go tool mockgen -package internal -destination internal/hash.go github.com/pocket-id/fosite Hasher
//go:generate go tool mockgen -package internal -destination internal/open_id_connect_token_strategy.go github.com/pocket-id/fosite/handler/openid OpenIDConnectTokenStrategy
//go:generate go tool mockgen -package internal -destination internal/open_id_connect_token_strategy_provider.go github.com/pocket-id/fosite/handler/openid OpenIDConnectTokenStrategyProvider
//go:generate go tool mockgen -package internal -destination internal/open_id_connect_request_storage.go github.com/pocket-id/fosite/handler/openid OpenIDConnectRequestStorage
//go:generate go tool mockgen -package internal -destination internal/open_id_connect_request_storage_provider.go github.com/pocket-id/fosite/handler/openid OpenIDConnectRequestStorageProvider
//go:generate go tool mockgen -package internal -destination internal/par_storage.go github.com/pocket-id/fosite PARStorage
//go:generate go tool mockgen -package internal -destination internal/par_storage_provider.go github.com/pocket-id/fosite PARStorageProvider
//go:generate go tool mockgen -package internal -destination internal/pkce_request_storage.go github.com/pocket-id/fosite/handler/pkce PKCERequestStorage
//go:generate go tool mockgen -package internal -destination internal/pkce_request_storage_provider.go github.com/pocket-id/fosite/handler/pkce PKCERequestStorageProvider
//go:generate go tool mockgen -package internal -destination internal/refresh_token_storage.go github.com/pocket-id/fosite/handler/oauth2 RefreshTokenStorage
//go:generate go tool mockgen -package internal -destination internal/refresh_token_storage_provider.go github.com/pocket-id/fosite/handler/oauth2 RefreshTokenStorageProvider
//go:generate go tool mockgen -package internal -destination internal/refresh_token_strategy.go github.com/pocket-id/fosite/handler/oauth2 RefreshTokenStrategy
//go:generate go tool mockgen -package internal -destination internal/refresh_token_strategy_provider.go github.com/pocket-id/fosite/handler/oauth2 RefreshTokenStrategyProvider
//go:generate go tool mockgen -package internal -destination internal/request.go github.com/pocket-id/fosite Requester
//go:generate go tool mockgen -package internal -destination internal/resource_owner_password_credentials_grant_storage.go github.com/pocket-id/fosite/handler/oauth2 ResourceOwnerPasswordCredentialsGrantStorage
//go:generate go tool mockgen -package internal -destination internal/resource_owner_password_credentials_grant_storage_provider.go github.com/pocket-id/fosite/handler/oauth2 ResourceOwnerPasswordCredentialsGrantStorageProvider
//go:generate go tool mockgen -package internal -destination internal/revocation_handler.go github.com/pocket-id/fosite RevocationHandler
//go:generate go tool mockgen -package internal -destination internal/revocation_handlers_provider.go github.com/pocket-id/fosite RevocationHandlersProvider
//go:generate go tool mockgen -package internal -destination internal/rfc7523_key_storage.go github.com/pocket-id/fosite/handler/rfc7523 RFC7523KeyStorage
//go:generate go tool mockgen -package internal -destination internal/rfc7523_key_storage_provider.go github.com/pocket-id/fosite/handler/rfc7523 RFC7523KeyStorageProvider
//go:generate go tool mockgen -package internal -destination internal/storage.go github.com/pocket-id/fosite Storage
//go:generate go tool mockgen -package internal -destination internal/token_endpoint_handler.go github.com/pocket-id/fosite TokenEndpointHandler
//go:generate go tool mockgen -package internal -destination internal/token_introspector.go github.com/pocket-id/fosite TokenIntrospector
//go:generate go tool mockgen -package internal -destination internal/token_revocation_storage.go github.com/pocket-id/fosite/handler/oauth2 TokenRevocationStorage
//go:generate go tool mockgen -package internal -destination internal/token_revocation_storage_provider.go github.com/pocket-id/fosite/handler/oauth2 TokenRevocationStorageProvider
//go:generate go tool mockgen -package internal -destination internal/transactional.go github.com/pocket-id/fosite Transactional
//go:generate go tool mockgen -package internal -destination internal/user_code_strategy.go github.com/pocket-id/fosite/handler/rfc8628 UserCodeStrategy
//go:generate go tool mockgen -package internal -destination internal/user_code_strategy_provider.go github.com/pocket-id/fosite/handler/rfc8628 UserCodeStrategyProvider
//go:generate go tool mockgen -package internal -destination internal/pushed_authorize_endpoint_handler.go github.com/pocket-id/fosite PushedAuthorizeEndpointHandler
