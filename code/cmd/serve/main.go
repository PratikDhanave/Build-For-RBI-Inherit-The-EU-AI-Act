// Zero-dependency Go HTTP server for previewing the interactive web deck locally
// with proper `application/wasm` MIME headers for the offline Go WebAssembly labs.
//
// Run from repo root via:
//
//	make serve
//	(cd code && go run ./cmd/serve --dir .. --addr :8000)
package main

import (
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
	"path/filepath"
)

func main() {
	dir := flag.String("dir", "..", "directory to serve")
	addr := flag.String("addr", ":8000", "listen address")
	flag.Parse()

	_ = mime.AddExtensionType(".wasm", "application/wasm")
	abs, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatalf("resolve dir: %v", err)
	}

	fmt.Printf("Serving %s on http://localhost%s (press Ctrl+C to stop)\n", abs, *addr)
	if err := http.ListenAndServe(*addr, http.FileServer(http.Dir(abs))); err != nil {
		log.Fatal(err)
	}
}
