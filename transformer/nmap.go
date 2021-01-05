package transformer

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/LeakIX/l9format"
	"io"
	"net"
	"strings"
)

type NmapTransformer struct {
	Transformer
	scanner *bufio.Scanner
}

func NewNmapTransformer() TransformerInterface {
	return &NmapTransformer{}
}

func (t *NmapTransformer) Name() string {
	return "nmap"
}

func (t *NmapTransformer) Decode(outputTransformer TransformerInterface) (err error) {
	if t.scanner == nil {
		t.scanner = bufio.NewScanner(t.Reader)
	}
	if t.scanner.Scan() {
		if strings.HasPrefix(t.scanner.Text(), "#") {
			return NewNoDataError("commented line")
		}
		if t.scanner.Text() == "" {
			return NewNoDataError("empty line")
		}
		inputParts := strings.Fields(t.scanner.Text())
		if len(inputParts) < 5 {
			return errors.New(fmt.Sprintf("couldn't parse %s", t.scanner.Text()))
		}
		if inputParts[3] != "Ports:" {
			return NewNoDataError("other line")
		}
		portParts := strings.Split(inputParts[len(inputParts)-1], "/")
		if len(portParts) < 3 {
			return errors.New(fmt.Sprintf("couldn't parse %s", t.scanner.Text()))
		}
		event := l9format.L9Event{
			Port:     portParts[0],
			Protocol: portParts[2],
			Host:     strings.TrimSuffix(inputParts[1], "[]"),
		}
		ip := net.ParseIP(event.Host)
		if ip != nil {
			event.Ip = ip.String()
		}
		return outputTransformer.Encode(event)
	}
	return io.EOF
}

func (t *NmapTransformer) Encode(event l9format.L9Event) error {
	if len(event.Host) < 1 {
		event.Host = event.Ip
	}
	if len(event.Protocol) < 1 {
		// default to tcp
		event.Protocol = "tcp"
	}
	hostPortString := fmt.Sprintf("Host: %s () Ports: %s/open/%s////\n", event.Host, event.Port, event.Protocol)
	written, err := io.WriteString(t.Writer, hostPortString)
	if err != nil {
		return err
	}
	if written < len(hostPortString) {
		return errors.New("output_overflow")
	}
	return nil
}
