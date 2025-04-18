package main

import (
	"log"
	"net/http"
)

func handlerHealthz(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type","text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("200 OK"))
}

func main() {
	const port = "8080"
	const filepathRoot = "."
	
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/healthz", handlerHealthz)
	serveMux.Handle("/app/", http.StripPrefix("/app/",http.FileServer(http.Dir(filepathRoot))))

	server := &http.Server{
		Addr: ":" + port,
		Handler: serveMux,
	}
	
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}