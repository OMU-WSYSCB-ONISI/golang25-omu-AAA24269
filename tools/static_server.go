package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":8000", "listen address")
	dir := flag.String("dir", "public", "directory to serve")
	flag.Parse()

	fs := http.FileServer(http.Dir(*dir))
	http.Handle("/", fs)

	log.Printf("Serving %s on http://%s", *dir, *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
