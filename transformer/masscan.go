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

type MasscanTransformer struct {
	Transformer
	scanner *bufio.Scanner
}

func NewMasscanTransformer() TransformerInterface {
	return &MasscanTransformer{}
}

func (t *MasscanTransformer) Name() string {
	return "masscan"
}

func (t *MasscanTransformer) Decode() (event l9format.L9Event, err error) {
	if t.scanner == nil {
		t.scanner = bufio.NewScanner(t.Reader)
	}
	if t.scanner.Scan() {
		if strings.HasPrefix(t.scanner.Text(), "#") {
			return event, NewNoDataError("commented line")
		}
		inputParts := strings.Fields(t.scanner.Text())
		if len(inputParts) < 6 {
			return event, errors.New(fmt.Sprintf("couldn't parse %s", t.scanner.Text()))
		}
		portParts := strings.Split(inputParts[len(inputParts)-3], "/")
		if len(portParts) < 2 {
			return event, errors.New(fmt.Sprintf("couldn't parse %s", t.scanner.Text()))
		}
		event.Port = portParts[0]
		event.Protocol = portParts[1]
		event.Host = strings.TrimSuffix(inputParts[len(inputParts)-1], "[]")
		ip := net.ParseIP(event.Host)
		if ip != nil {
			event.Ip = ip.String()
		}
	} else {
		return event, io.EOF
	}
	return event, err
}

func (t *MasscanTransformer) Encode(event l9format.L9Event) error {
	if len(event.Host) < 1 {
		event.Host = event.Ip
	}
	if len(event.Protocol) < 1 {
		// default to tcp
		event.Protocol = "tcp"
	}
	hostPortString := fmt.Sprintf("Discovered open port %s/%s on %s\n", event.Port, event.Protocol, event.Host)
	written, err := io.WriteString(t.Writer, hostPortString)
	if err != nil {
		return err
	}
	if written < len(hostPortString) {
		return errors.New("output_overflow")
	}
	return nil
}
