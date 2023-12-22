package observability

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

	errCount := 0
	reportError := func(err error) {
		if errCount > 2 {
			return
		}
		if errCount <= 2 {
			log.Print(err)
		}
		if errCount == 2 {
			log.Printf("error reporting paused until a good message arrives")
		}

		errCount++
	}

	for m := range messages {
		// scan the protocol
		if len(m) < 13 {
			reportError(fmt.Errorf("invalid incoming message: less than 13 bytes"))
			continue
		}

		b := bytes.NewBuffer(m)

		var t uint32
		err := binary.Read(b, binary.LittleEndian, &t)
		if err != nil {
			reportError(fmt.Errorf("invalid incoming message: unable to read first 4 bytes as a timestamp"))
			continue
		}

		var v int64
		err = binary.Read(b, binary.LittleEndian, &v)
		if err != nil {
			reportError(fmt.Errorf("invalid incoming message: unable to read 8 bytes as a value"))
			continue
		}

		ss := bytes.Split(b.Bytes(), []byte{'\x00'}) // null-separated
		if len(ss) == 0 {
			reportError(fmt.Errorf("invalid incoming message: unable to read metric name"))
			continue
		}

		m := Metric{}
		m.Timestamp = time.UnixMilli(int64(t) * 1000)
		m.Value = v
		m.Metric = string(ss[0])

		if len(ss) != 0 {
			m.Tag1, ss = string(ss[0]), ss[1:]
		}
		if len(ss) != 0 {
			m.Tag2, ss = string(ss[0]), ss[1:]
		}
		if len(ss) != 0 {
			m.Tag3, ss = string(ss[0]), ss[1:]
		}

		errCount = 0
	}
}
