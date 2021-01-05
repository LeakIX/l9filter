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

func (t *UrlServiceTransformer) Decode(outputTransformer TransformerInterface) (err error) {
	if t.scanner == nil {
		t.scanner = bufio.NewScanner(t.Reader)
	}
	if t.scanner.Scan() {
		parsedUrl, err := url.Parse(t.scanner.Text())
		if err != nil {
			return err
		}
		event := l9format.L9Event{
			Protocol: parsedUrl.Scheme,
			Port:     parsedUrl.Port(),
			Host:     parsedUrl.Hostname(),
			Http: l9format.L9HttpEvent{
				Url:  parsedUrl.RequestURI(),
				Root: parsedUrl.Path,
			},
			Transports: []string{"tcp", "http"},
		}
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
			return errors.New("no_port")
		}
		return outputTransformer.Encode(event)
	}
	return io.EOF
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
