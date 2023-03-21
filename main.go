package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Initialize a new Searcher struct
	searcher := Searcher{}

	// Load the contents of "completeworks.txt" into the Searcher
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Serve static files from the "./static" directory
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Register the "/search" route with the handleSearch function
	http.HandleFunc("/search", handleSearch(searcher))

	// Get the server's listening port from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	// Start the HTTP server and handle any errors
	fmt.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Define the Searcher struct
type Searcher struct {
	CompleteWorks string
	CompleteWorksLowercase string
	SuffixArray   *suffixarray.Index
}

// handleSearch takes a Searcher as a parameter and returns an HTTP handler function
func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for the presence of the search query in the URL parameters
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		// Call the Search method of the Searcher and encode the results as a JSON response
		results := searcher.Search(query[0])
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		// Set the content type of the response to "application/json" and send the response to the client
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

// Load reads the contents of the specified file, assigns the contents to the
// CompleteWorks field of the Searcher, and creates a new suffix array index
// from the file contents, assigning it to the SuffixArray field.
func (s *Searcher) Load(filename string) error {
	// Read the file contents.
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	// Assign the file contents to the CompleteWorks field.
	s.CompleteWorks = string(dat)
	// Create a lowercase version to make case insensitive searches
	s.CompleteWorksLowercase = strings.ToLower(s.CompleteWorks) 
	// Create a new suffix array index from the file contents and
	// assign it to the SuffixArray field.
	s.SuffixArray = suffixarray.New([]byte(s.CompleteWorksLowercase))
	return nil
}

// Search takes a query string as a parameter, searches the text using
// the suffix array index, and builds a slice of strings containing the
// surrounding 250 characters of each match found.
func (s *Searcher) Search(query string) []string {
	// Create lowercase version of the query
	lowercaseQuery := strings.ToLower(query) 
	// Search the text using the suffix array index.
	idxs := s.SuffixArray.Lookup([]byte(lowercaseQuery), -1)
	// Initialize a results slice to store the found matches.
	results := []string{}
	// Iterate over the indices of the found matches.
	for _, idx := range idxs {
		// Extract a substring around the match (250 characters before and after).
		results = append(results, s.CompleteWorks[idx-250:idx+250])
	}
	// Return the results slice.
	return results
}