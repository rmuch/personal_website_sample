package main

import (
	"github.com/gorilla/handlers"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Entry point for the application.
func main() {
	// We're going!
	log.Printf("Application server starting up...")

	log.Printf("Configuring request multiplexer...")

	// Set up request multiplexer.
	m := http.NewServeMux()

	// Set up routes.
	m.HandleFunc("/", rootHandler)
	m.HandleFunc("/home", homeHandler)
	m.HandleFunc("/projects", projectsHandler)
	m.HandleFunc("/about", aboutHandler)

	// Serve stylesheet directory. We can do the same for scripts, etc.
	m.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))

	// Wrap in our log filter, then in the Gorilla log filter.
	// It's not necessary to have both, so we can modify the code here to have one or the other.
	var w http.Handler = NewLogFilter(m)

	// Open log file.
	file, err := os.Create(logRoot + "log-" + time.Now().Format(strings.Replace(time.RFC3339, ":", "-", -1)) + ".log")
	if err != nil {
		log.Print(err)
		return
	}
	defer file.Close()

	w = handlers.CombinedLoggingHandler(io.MultiWriter(file, os.Stdout), w)

	// Set up the server.
	s := &http.Server{
		Addr:    listenInterface,
		Handler: w,
	}

	log.Printf("Routes configured, starting to serve...")

	// Start the server.
	log.Fatal(s.ListenAndServe())
}
