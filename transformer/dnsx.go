package transformer

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/LeakIX/l9format"
	"github.com/projectdiscovery/retryabledns"
	"time"
)

type DnsxTransformer struct {
	Transformer
	reader      *bufio.Reader
	jsonEncoder *json.Encoder
}

func NewDnsxTransformer() TransformerInterface {
	return &DnsxTransformer{}
}

func (t *DnsxTransformer) Decode(outputTransformer TransformerInterface) (err error) {
	if t.reader == nil {
		t.reader = bufio.NewReaderSize(t.Reader, 1024*1024)
	}
	dnsxLine := &retryabledns.DNSData{}
	bytes, isPrefix, err := t.reader.ReadLine()
	if err == nil && !isPrefix {
		err = json.Unmarshal(bytes, &dnsxLine)
		if err != nil {
			return err
		}
		// Output a l9event for each A record
		for _, ip := range dnsxLine.A {
			err = outputTransformer.Encode(l9format.L9Event{
				EventType: "resolve",
				Ip:        ip,
				Host:      dnsxLine.Host,
				Time:      time.Now(),
			})
			if err != nil {
				return err
			}
			// Output a l9event for all CNAMEs
			for _, cname := range dnsxLine.CNAME {
				err = outputTransformer.Encode(l9format.L9Event{
					EventType: "resolve",
					Ip:        ip,
					Host:      cname,
					Time:      time.Now(),
				})
			}
		}
		//Do stuff here
	} else if isPrefix {
		err = errors.New("line buffer overflow")
	}
	return err
}

func (t *DnsxTransformer) Encode(event l9format.L9Event) error {
	if t.jsonEncoder == nil {
		t.jsonEncoder = json.NewEncoder(t.Writer)
	}
	return t.jsonEncoder.Encode(retryabledns.DNSData{
		Host: event.Host,
		A:    []string{event.Ip},
	})
}

func (t *DnsxTransformer) Name() string {
	return "dnsx"
}
