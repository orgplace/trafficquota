package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	tqclient "github.com/orgplace/trafficquota/client"
)

type handler struct {
	client tqclient.Client
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	if ok, err := h.allowed(r); err != nil {
		log.Print(err)
		status = http.StatusInternalServerError
	} else if !ok {
		status = http.StatusTooManyRequests
	}

	w.WriteHeader(status)
	fmt.Fprintln(w, http.StatusText(status))
}

func (h *handler) allowed(r *http.Request) (bool, error) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false, err
	}

	// Take token with the remote host as a key
	return h.client.Take(host)
}

func main() {
	c, err := tqclient.NewInsecureClient("localhost:3895")
	if err != nil {
		panic(c)
	}

	const listen = ":8080"
	log.Printf("Serving HTTP on %s", listen)
	if err := http.ListenAndServe(
		listen,
		LogHTTP(&handler{client: c}),
	); err != nil {
		panic(err)
	}
}
