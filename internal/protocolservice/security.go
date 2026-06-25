package protocolservice

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

func contentSecurityPolicy(nonce string) string {
	return "default-src 'none'; img-src 'self'; style-src 'nonce-" + nonce + "'; script-src 'nonce-" + nonce + "'; base-uri 'none'; form-action 'none'; frame-ancestors 'none'"
}

func newNonce() (string, error) {
	var b [16]byte
	if _, err := io.ReadFull(rand.Reader, b[:]); err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(b[:]), nil
}
