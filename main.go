package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	http.HandleFunc("/events", EventsHandler)
	http.ListenAndServe(":8080", nil)
}

// The Event struct represents the messages we are sending over the SSE channel.
// They will be JSON encoded.
type Event struct {
	ID      uint   `json:"id"`      // Sequentual identifier of the message
	Message string `json:"message"` // Texual message
}

// splitDoubleNewline is a bufio.SplitFunc that splits input by double newlines.
func splitDoubleNewline(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Look for the double newline delimiter
	delimiter := []byte("\n\n")
	index := bytes.Index(data, delimiter)

	if index >= 0 {
		// Found the delimiter, return the token up to the delimiter
		return index + len(delimiter), data[:index], nil
	}

	// If we're at EOF and have remaining data, return it as the last token
	if atEOF && len(data) > 0 {
		return len(data), data, nil
	}

	// Request more data
	return 0, nil, nil
}

func EventsHandler(w http.ResponseWriter, r *http.Request) {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set CORS headers to allow all origins. You may want to restrict this to
	// specific origins in a production environment.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Sending events
	var i uint
	for i = 0; i < 10; i++ {
		e := Event{ID: i, Message: fmt.Sprintf("Event %d", i)}
		msg, err := json.Marshal(e)
		if err != nil {
			logrus.WithField("event", e).Error("Failed to marshal json")
		} else {
			logrus.WithFields(logrus.Fields{"event": e, "message": string(msg)}).Info("Writing event")
			fmt.Fprintf(w, "%s\n\n", msg)
			w.(http.Flusher).Flush()
		}
		time.Sleep(1 * time.Second)
	}
}
