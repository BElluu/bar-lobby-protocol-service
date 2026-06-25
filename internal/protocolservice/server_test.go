package protocolservice

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProtocolURLFromRequestBuildsGenericProtocolURL(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "known route",
			path: "/internal/ping",
			want: "barrts://internal/ping",
		},
		{
			name: "generic route",
			path: "/super/akcja",
			want: "barrts://super/akcja",
		},
		{
			name: "preserves raw query",
			path: "/lobby/invite?id=555&name=A%2BB&empty=&token=a%3Db%26c",
			want: "barrts://lobby/invite?id=555&name=A%2BB&empty=&token=a%3Db%26c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			got, err := protocolURLFromRequest(req)
			if err != nil {
				t.Fatalf("protocolURLFromRequest() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("protocolURLFromRequest() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestInvalidPathsRedirectToBAR(t *testing.T) {
	tests := []string{
		"/",
		"/internal",
		"/internal/ping/extra",
		"/internal//ping",
		"/internal/ping/",
		"/internal/<script>",
		"/internal/ping%2Fextra",
		"/http://example.com/open",
		"/internal/ping;rm",
		"/internal/..",
		"/internal/cmd.exe",
	}

	server := NewHandler()
	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()

			server.ServeHTTP(rec, req)

			assertRedirectToFallback(t, rec)
		})
	}
}

func TestProtocolPageEscapesAcceptedQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/lobby/invite?id=555&name=A%26B&token=a%3Db%26c", nil)
	rec := httptest.NewRecorder()

	NewHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "href=\"barrts://lobby/invite?id=555&amp;name=A%26B&amp;token=a%3Db%26c\"") {
		t.Fatal("rendered HTML does not safely escape ampersands in the fallback href")
	}
	if !strings.Contains(body, `window.location.href = "barrts:\/\/lobby\/invite?id=555\u0026name=A%26B\u0026token=a%3Db%26c"`) {
		t.Fatal("rendered HTML does not safely escape the JavaScript protocol URL")
	}
}

func TestUnsafeQueryRedirectsToBAR(t *testing.T) {
	tests := []string{
		`/lobby/invite?name=%22%3E%3Cscript%3Ealert(1)%3C%2Fscript%3E`,
		`/lobby/invite?cmd=$(calc)`,
		`/lobby/invite?cmd=%60calc%60`,
		`/lobby/invite?path=C%3A%5CWindows%5CSystem32`,
	}

	server := NewHandler()
	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()

			server.ServeHTTP(rec, req)

			assertRedirectToFallback(t, rec)
		})
	}
}

func TestSecurityHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/internal/ping", nil)
	rec := httptest.NewRecorder()

	NewHandler().ServeHTTP(rec, req)

	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatal("missing X-Content-Type-Options: nosniff")
	}
	csp := rec.Header().Get("Content-Security-Policy")
	if !strings.Contains(csp, "default-src 'none'") || !strings.Contains(csp, "frame-ancestors 'none'") {
		t.Fatalf("unexpected Content-Security-Policy: %q", csp)
	}
	if !strings.Contains(csp, "img-src 'self'") {
		t.Fatalf("Content-Security-Policy should allow only local images: %q", csp)
	}
	if strings.Contains(csp, "'unsafe-inline'") {
		t.Fatalf("Content-Security-Policy should not allow unsafe-inline: %q", csp)
	}
	if !strings.Contains(rec.Body.String(), ` nonce="`) {
		t.Fatal("rendered HTML does not include CSP nonces")
	}
}

func TestOnlyGetAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/internal/ping", nil)
	rec := httptest.NewRecorder()

	NewHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func assertRedirectToFallback(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()

	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusFound)
	}
	if location := rec.Header().Get("Location"); location != FallbackURL {
		t.Fatalf("Location = %q, want %q", location, FallbackURL)
	}
}
