package transformer

import (
	"bufio"
	"errors"
	"fmt"
	"gitlab.nobody.run/tbi/core"
	"io"
	"net"
	"strings"
)

type HostPortTransformer struct{
	Transformer
	scanner *bufio.Scanner
}


func NewHostPortTransformer() TransformerInterface{
	return &HostPortTransformer{}
}

func (t *HostPortTransformer) Name() string {
	return "hostport"
}

func (t *HostPortTransformer) Decode() (hostService core.HostService, err error) {
	if t.scanner == nil {
		t.scanner = bufio.NewScanner(t.Reader)
	}
	if t.scanner.Scan() {
		inputParts := strings.Split(t.scanner.Text(), ":")
		if len(inputParts) < 2 {
			return hostService, errors.New(fmt.Sprintf("couldn't parse %s", t.scanner.Text()))
		}
		hostService.Port = inputParts[len(inputParts)-1]
		host := strings.Trim(strings.TrimSuffix(t.scanner.Text(), ":" + hostService.Port), "[]")
		ip := net.ParseIP(host)
		if ip != nil {
			hostService.Ip = ip.String()
		} else {
			hostService.Hostname = host
		}
	} else {
		return hostService, io.EOF
	}
	return hostService, err
}

func (t *HostPortTransformer) Encode(hostService core.HostService) error {
	if len(hostService.Hostname) < 1 {
		hostService.Hostname = hostService.Ip
	}
	hostPortString := fmt.Sprintf("%s\n", net.JoinHostPort(hostService.Hostname, hostService.Port))
	written, err := io.WriteString(t.Writer, hostPortString)
	if err != nil {
		return err
	}
	if written < len(hostPortString) {
		return errors.New("output_overflow")
	}
	return nil
}