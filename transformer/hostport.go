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

type HostPortTransformer struct {
	Transformer
	scanner *bufio.Scanner
}

func NewHostPortTransformer() TransformerInterface {
	return &HostPortTransformer{}
}

func (t *HostPortTransformer) Name() string {
	return "hostport"
}

func (t *HostPortTransformer) Decode(outputTransformer TransformerInterface) (err error) {
	if t.scanner == nil {
		t.scanner = bufio.NewScanner(t.Reader)
	}
	if t.scanner.Scan() {
		inputParts := strings.Split(t.scanner.Text(), ":")
		if len(inputParts) < 2 {
			return errors.New(fmt.Sprintf("couldn't parse %s", t.scanner.Text()))
		}
		event := l9format.L9Event{
			Port: inputParts[len(inputParts)-1],
		}

		host := strings.Trim(strings.TrimSuffix(t.scanner.Text(), ":"+event.Port), "[]")
		ip := net.ParseIP(host)
		if ip != nil {
			event.Ip = ip.String()
		} else {
			event.Host = host
		}
		return outputTransformer.Encode(event)
	}
	return io.EOF
}

func (t *HostPortTransformer) Encode(event l9format.L9Event) error {
	if len(event.Host) < 1 {
		event.Host = event.Ip
	}
	hostPortString := fmt.Sprintf("%s\n", net.JoinHostPort(event.Host, event.Port))
	written, err := io.WriteString(t.Writer, hostPortString)
	if err != nil {
		return err
	}
	if written < len(hostPortString) {
		return errors.New("output_overflow")
	}
	return nil
}
