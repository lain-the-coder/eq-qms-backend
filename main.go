package main

import "net/http"

func WelcomeHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to EA QMS!"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", WelcomeHome)
	server := &http.Server{
		Addr:    ":1304",
		Handler: mux,
	}
	server.ListenAndServe()
}
