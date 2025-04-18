package main

import (
	"fmt"
	"os"
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()
	server := http.Server{
		Addr: ":8080",
		Handler: serveMux,
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Unable to start server: %v\n", err)
		os.Exit(1)
	}

}