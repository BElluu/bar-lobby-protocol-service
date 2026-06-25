package protocolservice

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var segmentPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func protocolURLFromRequest(r *http.Request) (string, error) {
	if !strings.HasPrefix(r.URL.Path, "/") {
		return "", errors.New("invalid protocol path")
	}

	segments := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(segments) != 2 || !validSegment(segments[0]) || !validSegment(segments[1]) {
		return "", errors.New("invalid protocol path")
	}
	if !validRawQuery(r.URL.RawQuery) {
		return "", errors.New("invalid protocol query")
	}

	target := "barrts://" + segments[0] + "/" + segments[1]
	if r.URL.RawQuery != "" {
		target += "?" + r.URL.RawQuery
	}

	return target, nil
}

func validSegment(segment string) bool {
	return segmentPattern.MatchString(segment)
}

func validRawQuery(rawQuery string) bool {
	if rawQuery == "" {
		return true
	}

	decoded, err := url.QueryUnescape(rawQuery)
	if err != nil {
		return false
	}

	return !strings.ContainsAny(decoded, "<>\"'`$\\\x00\r\n")
}
