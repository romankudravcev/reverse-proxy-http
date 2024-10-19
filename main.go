package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	targetURL = os.Getenv("TARGET_URL")
	port      = os.Getenv("PORT")
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for %s", r.Method, r.URL.Path)

	// Create a new request to the target cluster
	targetReq, err := http.NewRequest(r.Method, targetURL+r.URL.Path, r.Body)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		log.Printf("Error creating request: %v", err)
		return
	}

	// Copy headers
	for name, values := range r.Header {
		for _, value := range values {
			targetReq.Header.Add(name, value)
		}
	}

	// Send the request to the target cluster
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(targetReq)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusBadGateway)
		log.Printf("Error forwarding request: %v", err)
		return
	}
	defer resp.Body.Close()

	// Copy the response headers
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Set the status code and copy the response body
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Error copying response body: %v", err)
	}

	log.Printf("Successfully proxied request to %s with status %d", targetURL+r.URL.Path, resp.StatusCode)
}

func main() {
	if targetURL == "" || port == "" {
		log.Fatal("TARGET_URL and PORT environment variables must be set")
	}

	log.Printf("Starting proxy server on port %s", port)
	log.Printf("Redirecting to: %s", targetURL)
	http.HandleFunc("/", proxyHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
