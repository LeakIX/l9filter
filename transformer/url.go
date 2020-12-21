package transformer

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/LeakIX/l9format"
	"io"
	"net"
	"net/url"
)

type UrlServiceTransformer struct {
	Transformer
	scanner *bufio.Scanner
}

func NewUrlServiceTransformer() TransformerInterface {
	return &UrlServiceTransformer{}
}

func (t *UrlServiceTransformer) Name() string {
	return "url"
}

func (t *UrlServiceTransformer) Decode() (event l9format.L9Event, err error) {
	if t.scanner == nil {
		t.scanner = bufio.NewScanner(t.Reader)
	}
	if t.scanner.Scan() {
		parsedUrl, err := url.Parse(t.scanner.Text())
		if err != nil {
			return event, err
		}
		event.Protocol = parsedUrl.Scheme
		event.Port = parsedUrl.Port()
		event.Host = parsedUrl.Hostname()
		event.Http.Url = parsedUrl.RequestURI()
		event.Http.Root = parsedUrl.Path
		event.Transports = []string{"tcp", "http"}
		ip := net.ParseIP(parsedUrl.Hostname())
		if ip != nil {
			event.Ip = ip.String()
		}
		if event.Protocol == "https" {
			event.Transports = []string{"tcp", "tls", "http"}
			event.SSL.Enabled = true
		}
		if event.Port == "" {
			event.Port = SchemeDefaultPorts(parsedUrl.Scheme)
		}
		if event.Port == "" {
			// Couldn't get a port
			return event, errors.New("no_port")
		}
	} else {
		return event, io.EOF
	}
	return event, err
}

func (t *UrlServiceTransformer) Encode(event l9format.L9Event) error {
	if len(event.Host) < 1 {
		event.Host = event.Ip
	}
	scheme := "http"
	// best guess
	if event.SSL.Enabled {
		scheme = "https"
	}
	urlString := fmt.Sprintf("%s://%s\n", scheme, net.JoinHostPort(event.Host, event.Port))
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
	"ftp":    "21",
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
