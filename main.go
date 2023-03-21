package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
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
	CompleteWorks          string
	CompleteWorksLowercase string
	SuffixArray            *suffixarray.Index
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

		size, ok := r.URL.Query()["s"]
		if !ok || len(query[0]) < 1 {
			size[0] = "500"
		}

		// Call the Search method of the Searcher and encode the results as a JSON response
		intVar, err := strconv.Atoi(size[0])
		results := searcher.Search(query[0], intVar)

		// Encode response
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err = enc.Encode(results)
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

// String with the chars that are considered to be sentence separators, to identify the beginning
// and the end os sentences
var (
	sentenceSeparatorsString = ".,?!"
)

// Remove the incomplete sentences at the beggining and at the end of the original string
func TrimSentences(fullSentences string) string {
	// Find the index of the first separator at the string
	firstSeparatorIndex := strings.IndexAny(fullSentences, sentenceSeparatorsString)
	// Find the index of the last separator at the string
	lastSeparatorIndex := strings.LastIndexAny(fullSentences, sentenceSeparatorsString)

	// If they are not found or are the same return an empty string
	if firstSeparatorIndex < 0 || lastSeparatorIndex < 0 || firstSeparatorIndex == lastSeparatorIndex {
		return ""
	}

	// Return the string between the separators
	return fullSentences[firstSeparatorIndex+1 : lastSeparatorIndex+1]
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
func (s *Searcher) Search(query string, querySize int) []string {
	// Create lowercase version of the query
	lowercaseQuery := strings.ToLower(query)
	// Search the text using the suffix array index.
	idxs := s.SuffixArray.Lookup([]byte(lowercaseQuery), -1)
	// Initialize a results slice to store the found matches.
	results := []string{}
	// Iterate over the indices of the found matches.
	for _, idx := range idxs {
		// Extract a substring around the match (querySize/2 characters before and after).
		halfQuerySize := int(math.Floor(float64(querySize) / 2.0))
		textFound := s.CompleteWorks[idx-halfQuerySize : idx+halfQuerySize]
		// Replace the line breaks from txt to html line breaks, improving readability
		textFoundHtml := strings.Replace(textFound, "\r\n", "<br>", -1)
		// Append at the result array, with the sentences trimmed
		trimmedSentence := TrimSentences(textFoundHtml)
		if trimmedSentence != "" {
			results = append(results, trimmedSentence)
		}
	}
	// Return the results slice.
	return results
}
