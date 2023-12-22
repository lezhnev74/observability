package observability

import (
	"log"
	"time"
)

// This code covers listening for metrics sent from apps
// as opposed to actively polling for ones.
//
// The protocol of incoming messages is simple:
//
// |dateTime(4)|value(8)|metricName|...tagStrings...|
//     uint32    int64      string      strings
//
// bytes in numbers are little endian.
// Strings are separated by 0-byte. Only one string is mandatory (the metric name).

// ReadMetrics provides a channel to push received messages.
// Parsing is happening in "scanMessages", so the read can continue to the next message.
func ReadMetrics() chan<- []byte {
	messages := make(chan []byte, 1000)
	go scanMessages(messages)
	return messages
}

func scanMessages(messages <-chan []byte) {
	for m := range messages {
		// scan the protocol
		if len(m) < 13 {
			log.Printf("invalid incoming message: less than 13 bytes")
			continue
		}

		var t uint32

		m := Metric{}
		m.Timestamp = time.UnixMilli(t * 1000)
	}
}
