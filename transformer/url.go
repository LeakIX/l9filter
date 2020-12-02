package transformer

import (
	"bufio"
	"errors"
	"fmt"
	"gitlab.nobody.run/tbi/core"
	"io"
	"net"
	"net/url"
	"strings"
)

type UrlServiceTransformer struct{
	Transformer
	scanner *bufio.Scanner
}


func NewUrlServiceTransformer() TransformerInterface{
	return &UrlServiceTransformer{}
}

func (t *UrlServiceTransformer) Name() string {
	return "url"
}

func (t *UrlServiceTransformer) Decode() (hostService core.HostService, err error) {
	if t.scanner == nil {
		t.scanner = bufio.NewScanner(t.Reader)
	}
	if t.scanner.Scan() {
		parsedUrl, err := url.Parse(t.scanner.Text())
		if err != nil {
			return hostService, err
		}
		hostService.Scheme = parsedUrl.Scheme
		hostService.Type = parsedUrl.Scheme
		hostService.Port = parsedUrl.Port()
		hostService.Hostname = parsedUrl.Hostname()
		ip := net.ParseIP(parsedUrl.Hostname())
		if ip != nil {
			hostService.Ip = ip.String()
		}
		if hostService.Scheme == "https" {
			hostService.Type = "http"
		}
		if hostService.Port == "" {
			hostService.Port = SchemeDefaultPorts(parsedUrl.Scheme)
		}
		if hostService.Port == "" {
			// Couldn't get a port
			return hostService, errors.New("no_port")
		}
	} else {
		return hostService, io.EOF
	}
	return hostService, err
}

func (t *UrlServiceTransformer) Encode(hostService core.HostService) error {
	if len(hostService.Hostname) < 1 {
		hostService.Hostname = hostService.Ip
	}
	// best guess
	if len(hostService.Scheme) < 1 {
		hostService.Scheme = "http"
		if strings.HasSuffix(hostService.Port, "443") {
			hostService.Scheme = "https"
		}
	}
	urlString := fmt.Sprintf("%s://%s\n", hostService.Scheme, net.JoinHostPort(hostService.Hostname, hostService.Port))
	written, err := io.WriteString(t.Writer, urlString)
	if err != nil {
		return err
	}
	if written < len(urlString) {
		return errors.New("output_overflow")
	}
	return nil
}

var schemeDefaultPorts = map[string]string{
	"http":   "80",
	"https":  "443",
	"socks5": "1080",
	"ws":     "80",
	"wss":    "443",
	"ftp":	  "21",
	"mysql":  "3306",
	"ssh":    "22",
}

// SchemeDefaultPorts returns the default port for scheme s.
// If no default port is defined for scheme s then returns -1.
func SchemeDefaultPorts(s string) string {
	defaultPort, ok := schemeDefaultPorts[s]
	if !ok {
		return ""
	}
	return defaultPort
}