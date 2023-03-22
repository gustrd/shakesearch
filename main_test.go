package main

import (
	"encoding/json"
	"index/suffixarray"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// Define the SearchHandler struct that implements the http.Handler interface

type SearchHandler struct {
	Searcher Searcher
}

func (h *SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleSearch(h.Searcher)(w, r)
}

func Test_E2E_HandleSearch(t *testing.T) {
	// Initialize a new Searcher struct
	searcher := Searcher{}

	// Load the contents of "completeworks.txt" into the Searcher
	err := searcher.Load("completeworks.txt")
	if err != nil {
		t.Fatal(err)
	}

	// Create a test server using the SearchHandler
	addr := "localhost:3001"
	server := &http.Server{Addr: addr, Handler: &SearchHandler{Searcher: searcher}}
	var testfinish = false;
	go func() {
		if err := server.ListenAndServe(); err != nil && !testfinish {
			t.Errorf("server error: %v", err)
		}
	}()
	defer server.Close()

	// Send a search request to the server
	url := "http://" + addr + "/search?q=Luke&s=50&k=&mw=on"
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Verify that the response is successful
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %v", resp.StatusCode)
	}

	// Decode the response JSON
	var response SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that the response has the expected query, message, and results
	expectedQuery := "Luke"
	expectedMessage := "The query was 'Luke'. The search returned a total of 3 results."
	expectedResults := []SearchResult{
		{Text: "will presently to Saint Luke's; there,<br>    at", Play: "MACBETH - ACT III"},
		{Text: "old priest at Saint Luke’s church is at", Play: "THE TAMING OF THE SHREW - ACT IV"},
		{Text: "me to go to Saint Luke’s to bid<br>the", Play: "THE TAMING OF THE SHREW - ACT IV"},
	}
	if response.Query != expectedQuery {
		t.Errorf("unexpected query: %v", response.Query)
	}
	if response.Message != expectedMessage {
		t.Errorf("unexpected message: %v", response.Message)
	}
	if !reflect.DeepEqual(response.Results, expectedResults) {
		t.Errorf("unexpected results: %v", response.Results)
	}

	testfinish = true
}

func TestTrimSentences(t *testing.T) {
	// Set up test variables
	fullSentences := "entence. This is the second complete sentence. This is the third complete sentence. Thi"
	expectedResult := "This is the second complete sentence. This is the third complete sentence."

	// Call the function
	result := TrimSentences(fullSentences)

	// Check the result
	if result != expectedResult {
		t.Errorf("Expected result to be %q but got %q", expectedResult, result)
	}
}

func TestLoad(t *testing.T) {
	// Set up test variables
	searcher := &Searcher{}
	filename := "completeworks.txt"

	// Call the function
	err := searcher.Load(filename)

	// Check the result
	if err != nil {
		t.Errorf("Expected error to be nil but got %v", err)
	}
	if searcher.CompleteWorks == "" {
		t.Error("Expected CompleteWorks to be non-empty but it is empty")
	}
	if searcher.SuffixArray == nil {
		t.Error("Expected SuffixArray to be non-nil but it is nil")
	}
}

func TestRecoverWorkTitle(t *testing.T) {
	// Set up test variables
	searcher := &Searcher{
		CompleteWorks: "THE SONNETS\r\n\r\nContents\r\n\r\nby William Shakespeare\r\n\r\nI\r\n\r\nFrom fairest creatures we desire increase,\r\nThat thereby beauty's rose might never die,\r\nBut as the riper should by time decease,\r\nHis tender heir might bear his memory:",
	}
	idx := 28

	// Call the function
	result := searcher.RecoverWorkTitle(idx)

	// Check the result
	expectedResult := "THE SONNETS"
	if result != expectedResult {
		t.Errorf("Expected result to be %q but got %q", expectedResult, result)
	}
}

func TestRecoverMatchAct(t *testing.T) {
	// Set up test variables
	testCompleteWorks := "ACT I\r\nSCENE I. Elsinore. A platform before the Castle.\r\n\r\nFRANCISCO at his post. Enter to him BERNARDO\r\n\r\nBERNARDO\r\nWho's there?\r\n\r\nFRANCISCO\r\nNay, answer me: stand, and unfold yourself.\r\n\r\nBERNARDO\r\nLong live the king!\r\n\r\nFRANCISCO\r\nBernardo?\r\n\r\nBERNARDO\r\nHe.\r\n\r\nFRANCISCO\r\nYou come most carefully upon your hour.\r\n\r\nBERNARDO\r\n'Tis now struck twelve; get thee to bed, Francisco.\r\n\r\nFRANCISCO\r\nFor this relief much thanks: 'tis bitter cold,\r\nAnd I am sick at heart."
	searcher := &Searcher{
		CompleteWorks: testCompleteWorks,
	}
	idx := 35

	// Call the function
	result := searcher.RecoverMatchAct(idx)

	// Check the result
	expectedResult := "ACT I"
	if result != expectedResult {
		t.Errorf("Expected result to be %q but got %q", expectedResult, result)
	}
}

func TestSearchMatchWholeWordFalse(t *testing.T) {
	// Set up test variables
	testCompleteWorks := "To be or not to be, that is the question:\r\nWhether 'tis nobler in the mind to suffer\r\nThe slings and arrows of outrageous fortune,\r\nOr to take arms against a sea of troubles,\r\nAnd by opposing end them?"
	completeWorksLowercase := strings.ToLower(testCompleteWorks)
	searcher := &Searcher{
		SuffixArray:   suffixarray.New([]byte(completeWorksLowercase)),
		CompleteWorks: testCompleteWorks,
	}
	query := "be"
	querySize := 10
	useMatchWholeWord := false

	// Call the function
	result := searcher.Search(query, querySize, useMatchWholeWord)

	// Check the result
	expectedResult := []SearchResult{
		{
			Text: "be",
			Play: "?",
		},
		{
			Text: "to be,",
			Play: "?",
		},
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected result to be %v but got %v", expectedResult, result)
	}
}

func TestSearchMatchWholeWordTrue(t *testing.T) {
	// Set up test variables
	testCompleteWorks := "To be or not to be, that is the question:\r\nWhether 'tis nobler in the mind to suffer\r\nThe slings and arrows of outrageous fortune,\r\nOr to take arms against a sea of troubles,\r\nAnd by opposing end them?"
	completeWorksLowercase := strings.ToLower(testCompleteWorks)
	searcher := &Searcher{
		SuffixArray:   suffixarray.New([]byte(completeWorksLowercase)),
		CompleteWorks: testCompleteWorks,
	}
	query := "questio"
	querySize := 10
	useMatchWholeWord := true

	// Call the function
	result := searcher.Search(query, querySize, useMatchWholeWord)

	// Check the result
	expectedResult := []SearchResult{}
	
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected result to be %v but got %v", expectedResult, result)
	}
}