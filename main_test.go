package main

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_eventsHandler(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(EventsHandler))
	defer ts.Close()

	// Make the reqeust to the test server
	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Check the response code
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Read the response(s) delimited by double-newlines
	scanner := bufio.NewScanner(res.Body)
	scanner.Split(splitDoubleNewline)
	count := 0
	for scanner.Scan() {
		logrus.Info(scanner.Text())
		count++
	}
	assert.Equal(t, 10, count, "Missing some or all responses")
}
