# BAR Lobby Protocol Service

This is a small external HTTP service for BAR Lobby desktop protocol links.

BAR Lobby registers the `barrts://` custom protocol handler on the user's system. Public links can use normal HTTPS URLs, for example:

```text
https://bar.devopsowy.pl/internal/ping
https://bar.devopsowy.pl/lobby/invite?id=555
```

The service turns those requests into protocol links for the desktop app:

```text
/internal/ping        -> barrts://internal/ping
/lobby/invite?id=555 -> barrts://lobby/invite?id=555
/super/akcja         -> barrts://super/akcja
```

Instead of relying on an HTTP redirect to a custom protocol, the service returns an HTML page. The page attempts to open BAR Lobby with JavaScript and also shows a fallback link: `Click here if nothing happens`.

The protocol URL is not displayed to the user.

Requests that do not match `/<handler>/<action>` redirect to `https://www.beyondallreason.info`.

## Running

```sh
go run .
```

By default the service listens on `:47777`.

To use a different address:

```sh
ADDR=:3000 go run .
```

## Deployment With Docker

Build the image:

```sh
docker build -t bar-lobby-protocol-service .
```

Run it locally:

```sh
docker run --rm -p 47777:47777 bar-lobby-protocol-service
```

Run it with a custom listen address inside the container:

```sh
docker run --rm -e ADDR=:3000 -p 3000:3000 bar-lobby-protocol-service
```

For production, run the container behind your reverse proxy and point the public host, for example `https://bar.devopsowy.pl`, to the container port. The service itself listens on HTTP; TLS should normally be terminated by the reverse proxy or platform load balancer.

## Deployment Without Docker

Build a native binary:

```sh
go build -trimpath -ldflags="-s -w" -o bar-lobby-protocol-service .
```

Run it from the project directory so the local `assets/` directory is available:

```sh
./bar-lobby-protocol-service
```

Use a custom listen address if needed:

```sh
ADDR=:3000 ./bar-lobby-protocol-service
```

For deployment, copy the binary and the `assets/` directory together:

```text
bar-lobby-protocol-service
assets/
```

Then run the binary under your process manager, for example systemd, supervisor, or your hosting platform. As with Docker, put it behind a reverse proxy or load balancer for HTTPS.

## Tests

```sh
go test ./...
```

The tests cover:

- generic `/<handler>/<action>` protocol URL generation
- preserving accepted raw query parameters
- redirecting invalid paths
- redirecting suspicious query strings
- HTML and JavaScript escaping of generated links
- basic security headers

## Project Structure

```text
main.go                            process entrypoint
internal/protocolservice/server.go HTTP server, routing, headers, response handling
internal/protocolservice/protocol.go request validation and barrts:// URL generation
internal/protocolservice/security.go CSP and nonce generation
internal/protocolservice/template.go HTML page template
internal/protocolservice/server_test.go behavior and security tests
assets/                            local page images
```

## Routing

Valid URLs must have exactly two path segments:

```text
/<handler>/<action>
```

Both segments currently accept only:

```text
A-Z a-z 0-9 _ -
```

This keeps the service generic while rejecting slashes, dot paths, encoded slash tricks, shell metacharacters, and HTML-related characters in the path.

Any other path redirects to:

```text
https://www.beyondallreason.info
```

## Security Model

The service never opens applications itself and never executes system commands. It only returns HTML containing a `barrts://` URL derived from the validated request.

Security controls:

- only `GET` is accepted
- only `barrts://` links are generated
- path must be exactly two simple segments
- query strings are preserved only when they decode without suspicious characters such as `<`, `>`, quotes, backticks, `$`, backslashes, null bytes, or newlines
- invalid or suspicious requests redirect only to `https://www.beyondallreason.info`
- output is rendered with Go's `html/template`
- fallback `href` uses a trusted `template.URL` only after validation has forced the scheme to `barrts://`
- Content Security Policy uses per-response nonces for inline CSS and JavaScript
- images are served from local files under `/assets/`, not from an external CDN
- `X-Content-Type-Options: nosniff` is set
- the HTTP server has read/write/header/idle timeouts

This service is not an allowlist of BAR Lobby actions. It intentionally accepts any valid two-segment route so new desktop actions can be tested without redeploying this service.

If the action surface becomes stable, adding an allowlist for supported handler/action pairs would be a good next hardening step.
