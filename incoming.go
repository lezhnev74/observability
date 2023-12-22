package observability

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

// ParseMetric provides a channel to push received messages.
// Parsing is happening in "scanMessages", so the read can continue to the next message.
func ParseMetric(payload []byte) (m Metric, err error) {
	// scan the protocol
	if len(payload) < 13 {
		err = fmt.Errorf("invalid incoming message: less than 13 bytes")
		return
	}

	b := bytes.NewBuffer(payload)

	var t uint32
	err = binary.Read(b, binary.LittleEndian, &t)
	if err != nil {
		err = fmt.Errorf("invalid incoming message: unable to read first 4 bytes as a timestamp")
		return
	}

	var v int64
	err = binary.Read(b, binary.LittleEndian, &v)
	if err != nil {
		err = fmt.Errorf("invalid incoming message: unable to read 8 bytes as a value")
		return
	}

	ss := bytes.Split(b.Bytes(), []byte{'\x00'}) // null-separated
	if len(ss) == 0 {
		err = fmt.Errorf("invalid incoming message: unable to read metric name")
		return
	}

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

	return
}
