package protocolservice

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultAddr = ":47777"
	FallbackURL = "https://www.beyondallreason.info"
)

type pageData struct {
	ProtocolURL  string
	ProtocolHref template.URL
	Nonce        string
}

func NewHTTPServer(addr string) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           NewHandler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func NewHandler() http.Handler {
	assets := http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/assets/") {
			assets.ServeHTTP(w, r)
			return
		}
		handleProtocolRequest(w, r)
	})
}

func handleProtocolRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	protocolURL, err := protocolURLFromRequest(r)
	if err != nil {
		http.Redirect(w, r, FallbackURL, http.StatusFound)
		return
	}

	nonce, err := newNonce()
	if err != nil {
		log.Printf("generate CSP nonce: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", contentSecurityPolicy(nonce))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	err = pageTemplate.Execute(w, pageData{
		ProtocolURL:  protocolURL,
		ProtocolHref: template.URL(protocolURL),
		Nonce:        nonce,
	})
	if err != nil {
		log.Printf("render protocol page: %v", err)
	}
}
