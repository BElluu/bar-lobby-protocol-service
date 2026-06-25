package main

import (
	"log"
	"os"

	"bar-lobby-protocol-service/internal/protocolservice"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = protocolservice.DefaultAddr
	}

	log.Printf("BAR Lobby Protocol Service listening on %s", addr)
	if err := protocolservice.NewHTTPServer(addr).ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
