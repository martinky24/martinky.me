package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Static files with security headers
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", securityHeaders(http.StripPrefix("/static/", fs)))

	// Health check endpoint
	http.HandleFunc("/health", healthHandler)

	// Template routes with security headers
	http.HandleFunc("/", securityHeaders(http.HandlerFunc(serveTemplate)).ServeHTTP)

	port := os.Getenv("MARTINKY_ME_PORT")
	if port == "" {
		slog.Info("using default port", "port", "8090")
		port = "8090"
	}

	slog.Info("server starting", "port", port, "url", "http://localhost:"+port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

// securityHeaders adds security headers to all responses
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")
		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Enable browser XSS protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		// Content Security Policy
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline' https://stackpath.bootstrapcdn.com https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' data:;")
		// Referrer policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		next.ServeHTTP(w, r)
	})
}

// healthHandler returns OK for health checks
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	slog.Info("request", "method", r.Method, "path", r.URL.Path, "remote", r.RemoteAddr)
	path := checkExt(r.URL.Path)
	if r.URL.Path == "/" {
		path = "/index.html"
	}

	partials := filepath.Join("templates", "partials.html")
	fp := filepath.Join("templates", filepath.Clean(path))

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(fp, partials)
	if err != nil {
		slog.Error("template parse error", "error", err, "path", fp)
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, path, nil)
	if err != nil {
		slog.Error("template execute error", "error", err, "path", path)
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
	}
}

func checkExt(url string) string {
	suffix := ".html"
	if strings.HasSuffix(url, suffix) {
		return url
	}
	return url + suffix
}
