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
		if !ok || len(size[0]) < 1 {
			size[0] = "500"
		}

		openAiApiKey, ok := r.URL.Query()["k"]
		if !ok || len(openAiApiKey[0]) < 1 {
			openAiApiKey[0] = ""
		}

		// Call the Search method of the Searcher and encode the results as a JSON response
		intVar, err := strconv.Atoi(size[0])
		results := searcher.Search(query[0], intVar)

		//Verify if results is empty, if it the query can be sent to autocorrect
		if len(results) == 0 && openAiApiKey[0] != "" {
			correctedQuery := searcher.Correct(query[0], openAiApiKey[0])
			results = searcher.Search(correctedQuery, intVar)
		}

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
	sentenceSeparatorsString = ".,?! "
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

// Recover the work title that is before given index by using heuristics.
func (s *Searcher) RecoverWorkTitle(idx int) string {
	// Cuts the complete works until the point of the found index
	workToIndex := s.CompleteWorks[:idx]
	// Split the string at newlines
	linesList := strings.Split(workToIndex, "\r\n")
	// Variable to store if 'Contents' was found
	contentsLineFound := false
	// Loops through the list at the inverse order
	for i := len(linesList) - 1; i >= 0; i-- {
		// If the contents line was found, return the next not empty string
		if contentsLineFound && linesList[i] != "" {
			return linesList[i]
		}
		// Searches for the line 'Contents', that by definition is right after the title
		if linesList[i] == "Contents" {
			contentsLineFound = true
		}
	}

	// If no title was found return "?"
	return "?"
}

type SearchResult struct {
	Text      string
	WorkTitle string
}

// Search takes a query string as a parameter, searches the text using
// the suffix array index, and builds a slice of strings containing the
// surrounding 250 characters of each match found.
func (s *Searcher) Search(query string, querySize int) []SearchResult {
	// Create lowercase version of the query
	lowercaseQuery := strings.ToLower(query)
	// Search the text using the suffix array index.
	idxs := s.SuffixArray.Lookup([]byte(lowercaseQuery), -1)
	// Initialize a results slice to store the found matches.
	results := []SearchResult{}
	// Iterate over the indices of the found matches.
	for _, idx := range idxs {
		workTitle := s.RecoverWorkTitle(idx)
		// Extract a substring around the match (querySize/2 characters before and after).
		halfQuerySize := int(math.Floor(float64(querySize) / 2.0))
		textFound := s.CompleteWorks[idx-halfQuerySize : idx+halfQuerySize]
		// Replace the line breaks from txt to html line breaks, improving readability
		textFoundHtml := strings.Replace(textFound, "\r\n", "<br>", -1)
		// Append at the result array, with the sentences trimmed
		trimmedSentence := TrimSentences(textFoundHtml)
		if trimmedSentence != "" {
			results = append(results, SearchResult{Text: trimmedSentence, WorkTitle: workTitle})
		}
	}
	// Return the results slice.
	return results
}

type OpenAiJson struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		Logprobs     any    `json:"logprobs"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Uses OpenAI GPT-3 API to try to correct a misspeled word or sentence
func (s *Searcher) Correct(query string, apiKey string) string {
	// Set API endpoint
	apiURL := "https://api.openai.com/v1/completions"

	// Set API request headers
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + apiKey,
	}

	// Set API request data
	data := map[string]interface{}{
		"model": "text-davinci-003",
		"prompt": "The following sentece from Shakespeare's work is misspelled. Give me the correct sentence.\n\"" +
			query +
			"\"",
		"temperature":       0.7,
		"max_tokens":        256,
		"top_p":             1,
		"frequency_penalty": 0,
		"presence_penalty":  0,
	}

	// Marshal data to JSON
	payload, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	// Create HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}

	// Set headers on request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send request and get response
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// Decode response JSON
	var response OpenAiJson
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	//Remove quotes and new lines
	responseString := strings.ReplaceAll(response.Choices[0].Text, "\"", "")
	responseString = strings.ReplaceAll(responseString, "\n", "")

	// Return string without the new lines and without double quotes
	return strings.ReplaceAll(responseString, "\n", "")
}
