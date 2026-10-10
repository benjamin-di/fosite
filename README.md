# <div align="center"><img src="https://github.com/user-attachments/assets/d906ac59-9269-4aa0-8c0d-e9619d0de04c" width="100"/> </br>Pocket ID Fosite</div>

The OAuth 2.0 and OpenID Connect framework behind
[Pocket ID](https://github.com/pocket-id/pocket-id).

This repository is a fork of [Ory Fosite](https://github.com/ory/fosite). Ory no
longer develops the standalone repository and maintains Fosite in
[ory/hydra](https://github.com/ory/hydra/tree/master/fosite) instead. This fork
ports those upstream changes and adds the features Pocket ID needs. It is
maintained for Pocket ID, so its API follows Pocket ID's requirements and can
change between releases.

## Specifications

- [RFC 6749: The OAuth 2.0 Authorization Framework](https://www.rfc-editor.org/rfc/rfc6749)
- [RFC 6819: OAuth 2.0 Threat Model and Security Considerations](https://www.rfc-editor.org/rfc/rfc6819)
- [RFC 7009: OAuth 2.0 Token Revocation](https://www.rfc-editor.org/rfc/rfc7009)
- [RFC 7523: JWT Profile for OAuth 2.0 Client Authentication and Authorization Grants](https://www.rfc-editor.org/rfc/rfc7523)
- [RFC 7636: Proof Key for Code Exchange (PKCE)](https://www.rfc-editor.org/rfc/rfc7636)
- [RFC 7662: OAuth 2.0 Token Introspection](https://www.rfc-editor.org/rfc/rfc7662)
- [RFC 8252: OAuth 2.0 for Native Apps](https://www.rfc-editor.org/rfc/rfc8252)
- [RFC 8628: OAuth 2.0 Device Authorization Grant](https://www.rfc-editor.org/rfc/rfc8628)
- [RFC 8707: Resource Indicators for OAuth 2.0](https://www.rfc-editor.org/rfc/rfc8707)
- [RFC 9068: JWT Profile for OAuth 2.0 Access Tokens](https://www.rfc-editor.org/rfc/rfc9068)
- [RFC 9126: OAuth 2.0 Pushed Authorization Requests](https://www.rfc-editor.org/rfc/rfc9126)
- [RFC 9207: OAuth 2.0 Authorization Server Issuer Identification](https://www.rfc-editor.org/rfc/rfc9207)
- [OAuth Client ID Metadata Document](https://datatracker.ietf.org/doc/draft-ietf-oauth-client-id-metadata-document/)
- [OpenID Connect Core 1.0](https://openid.net/specs/openid-connect-core-1_0.html)

## Usage

```sh
go get github.com/pocket-id/fosite
```

See the [releases](https://github.com/pocket-id/fosite/releases) for the
available versions.

## Development

```sh
go test ./...  # run the tests
make format    # fix license headers, Go imports and Markdown formatting
make gen       # regenerate the mocks in internal/
```

Ory's copyright stays on files inherited from upstream. `make format` marks
inherited files that the fork changes and adds the Pocket ID copyright to new
files.

## License

Apache 2.0, see [LICENSE](LICENSE). Fosite was originally created by
[Ory](https://www.ory.com).
