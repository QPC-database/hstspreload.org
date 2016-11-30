package main

import (
	"net"
	"net/http"
)

type hstsServer struct{}

func (server hstsServer) Handle(pattern string, handler http.Handler) {
	server.HandleFunc(pattern, handler.ServeHTTP)
}

func (hstsServer) HandleFunc(pattern string, handlerFunc http.HandlerFunc) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if hsts(w, r) {
			handlerFunc(w, r)
		}
	})
}

func isLocalhost(hostport string) bool {
	host, _, err := net.SplitHostPort(hostport)
	return err == nil && host == "localhost"
}

// `cont` indicates whether the callee should continue further processing.
func hsts(w http.ResponseWriter, r *http.Request) (cont bool) {

	switch {
	case (r.TLS != nil), maybeAppspotHTTPS(r):
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		return true
	case isLocalhost(r.Host):
		return true
	default:
		// The redirect below causes problems with Managed VMs/Flexible Environments.
		// In a standalone server we'd handle the redirect here, but we let app.yaml
		// handle it for now.

		// u := fmt.Sprintf("https://%s%s", r.Host, r.URL.Path)
		// http.Redirect(w, r, u, http.StatusMovedPermanently)
		return false
	}
}

// Note: This can be spoofed when not run on App Engine/Flexible Environment.
func maybeAppspotHTTPS(r *http.Request) bool {
	return r.Header.Get("X-Appengine-Https") == "on"
}
