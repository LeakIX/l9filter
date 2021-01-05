// WIP, transformer for LeakIX legacy format

package transformer

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/LeakIX/l9format"
	"gitlab.nobody.run/tbi/core"
	"strings"
	"time"
)

type TbiCoreTransformer struct {
	Transformer
	reader      *bufio.Reader
	jsonEncoder *json.Encoder
}

func NewTbiCoreTransformer() TransformerInterface {
	return &TbiCoreTransformer{}
}

func (t *TbiCoreTransformer) Decode(outputTransformer TransformerInterface) (err error) {
	if t.reader == nil {
		t.reader = bufio.NewReaderSize(t.Reader, 1024*1024)
	}
	hostServiceLeak := &core.HostServiceLeak{}
	hostService := &core.HostService{}
	bytes, isPrefix, err := t.reader.ReadLine()
	if err == nil && !isPrefix {
		err = json.Unmarshal(bytes, &hostServiceLeak)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bytes, &hostService)
		if err != nil {
			return err
		}
		if len(hostServiceLeak.Plugin) > 0 {
			event, err := t.decodeLeak(hostServiceLeak)
			if err != nil {
				return err
			}
			return outputTransformer.Encode(event)
		} else {
			event, err := t.decodeService(hostService)
			if err != nil {
				return err
			}
			return outputTransformer.Encode(event)
		}
	} else if isPrefix {
		err = errors.New("line buffer overflow")
	}
	return err
}

func (t *TbiCoreTransformer) decodeService(hostService *core.HostService) (l9format.L9Event, error) {
	event := l9format.L9Event{
		EventType:     "service",
		EventSource:   "l9filter-tbi",
		EventPipeline: []string{"l9filter-tbi"},
		Ip:            hostService.Ip,
		Host:          hostService.Hostname,
		Reverse:       hostService.Reverse,
		Port:          hostService.Port,
		Transports:    []string{"tcp"},
		Protocol:      hostService.Type,
		Summary:       hostService.Data,
		Time:          time.Unix(hostService.Date, 0),
		SSL: l9format.L9SSLEvent{
			JARM:        hostService.Certificate.JARM,
			CypherSuite: hostService.Certificate.CypherSuite,
			Version:     hostService.Certificate.Version,
			Certificate: l9format.Certificate{
				CommonName:  hostService.Certificate.CommonName,
				Domains:     hostService.Certificate.Domains,
				Fingerprint: hostService.Certificate.Fingerprint,
				KeyAlgo:     hostService.Certificate.KeyAlgo,
				KeySize:     hostService.Certificate.KeySize,
				IssuerName:  hostService.Certificate.IssuerName,
				NotBefore:   hostService.Certificate.NotBefore,
				NotAfter:    hostService.Certificate.NotAfter,
				Valid:       hostService.Certificate.Valid,
			},
		},
		Service: l9format.L9ServiceEvent{
			Software: l9format.Software{
				Name:            hostService.Software.Name,
				Version:         hostService.Software.Version,
				OperatingSystem: hostService.Software.OperatingSystem,
				Fingerprint:     hostService.Software.Fingerprint,
			},
		},
	}
	event.Http.Headers = make(map[string]string)
	for headerName, headerValues := range hostService.Headers {
		if len(headerValues) < 1 {
			continue
		}
		event.Http.Headers[headerName] = headerValues[0]
	}
	for _, softwareModule := range hostService.Software.Modules {
		event.Service.Software.Modules = append(event.Service.Software.Modules, l9format.SoftwareModule{
			Name:        softwareModule.Name,
			Version:     softwareModule.Version,
			Fingerprint: softwareModule.Fingerprint,
		})
	}
	if strings.HasPrefix(hostService.Type, "http") {
		if hostService.Scheme == "https" {
			event.SSL.Enabled = true
			event.Transports = append(event.Transports, "tls")
		}
		event.Transports = append(event.Transports, "http")
	}
	if len(event.SSL.CypherSuite) > 0 || !event.SSL.Enabled {
		event.SSL.Enabled = true
		event.Transports = append(event.Transports, "tls")
	}
	return event, nil
}

func (t *TbiCoreTransformer) decodeLeak(hostServiceLeak *core.HostServiceLeak) (l9format.L9Event, error) {
	event := l9format.L9Event{
		EventType:     "leak",
		EventSource:   "l9filter-tbi",
		EventPipeline: []string{hostServiceLeak.Plugin, "l9filter-tbi"},
		Ip:            hostServiceLeak.Ip,
		Host:          hostServiceLeak.Hostname,
		Reverse:       hostServiceLeak.Reverse,
		Port:          hostServiceLeak.Port,
		Transports:    []string{"tcp"},
		Protocol:      hostServiceLeak.Type,
		Summary:       hostServiceLeak.Data,
		Time:          time.Unix(hostServiceLeak.Date, 0),
		Leak: l9format.L9LeakEvent{
			Data: hostServiceLeak.Data,
			Dataset: l9format.DatasetSummary{
				Rows:        hostServiceLeak.DatasetLeak.TotalRows,
				Size:        hostServiceLeak.DatasetLeak.TotalSizeByte,
				Collections: hostServiceLeak.DatasetLeak.TotalCollections,
				Infected:    hostServiceLeak.DatasetLeak.Infected,
			},
		},
	}
	if strings.HasPrefix(hostServiceLeak.Type, "http") || hostServiceLeak.Type == "web" {
		event.Protocol = "http"
		if hostServiceLeak.Scheme == "https" {
			event.SSL.Enabled = true
			event.Transports = append(event.Transports, "tls")
		}
		event.Transports = append(event.Transports, "http")
	}
	return event, nil
}

func (t *TbiCoreTransformer) Encode(event l9format.L9Event) error {
	if t.jsonEncoder == nil {
		t.jsonEncoder = json.NewEncoder(t.Writer)
	}
	switch event.EventType {
	case "leak":
		return t.encodeLeak(event)
	default:
		return t.encodeService(event)
	}
}

func (t *TbiCoreTransformer) encodeService(event l9format.L9Event) error {
	hostService := &core.HostService{
		Ip:   event.Ip,
		Port: event.Port,
		Type: event.Protocol,
		Credentials: []*core.HostServiceCredentials{{
			NoAuth:   event.Service.Credentials.NoAuth,
			Username: event.Service.Credentials.Username,
			Password: event.Service.Credentials.Password,
			Key:      event.Service.Credentials.Key,
			Raw:      event.Service.Credentials.Raw,
		}},
		Software: core.Software{
			Name:            event.Service.Software.Name,
			Version:         event.Service.Software.Version,
			OperatingSystem: event.Service.Software.OperatingSystem,
			Fingerprint:     event.Service.Software.Fingerprint,
		},
		Date:     event.Time.Unix(),
		Data:     event.Summary,
		Scheme:   event.Protocol,
		Hostname: event.Host,
		Reverse:  event.Reverse,
		Certificate: core.HostCertificate{
			JARM:        event.SSL.JARM,
			CommonName:  event.SSL.Certificate.CommonName,
			Domains:     event.SSL.Certificate.Domains,
			Fingerprint: event.SSL.Certificate.Fingerprint,
			KeyAlgo:     event.SSL.Certificate.KeyAlgo,
			KeySize:     event.SSL.Certificate.KeySize,
			CypherSuite: event.SSL.CypherSuite,
			Version:     event.SSL.Version,
			IssuerName:  event.SSL.Certificate.IssuerName,
			NotBefore:   event.SSL.Certificate.NotBefore,
			NotAfter:    event.SSL.Certificate.NotAfter,
			Valid:       event.SSL.Certificate.Valid,
		},
	}
	hostService.Headers = make(map[string][]string)
	for headerName, headerValue := range event.Http.Headers {
		hostService.Headers[headerName] = []string{headerValue}
	}
	for _, softwareModule := range event.Service.Software.Modules {
		hostService.Software.Modules = append(hostService.Software.Modules, &core.SoftwareModule{
			Name:        softwareModule.Name,
			Version:     softwareModule.Version,
			Fingerprint: softwareModule.Fingerprint,
		})
	}
	return t.jsonEncoder.Encode(hostService)
}

func (t *TbiCoreTransformer) encodeLeak(event l9format.L9Event) error {
	hostServiceLeak := &core.HostServiceLeak{
		Ip:       event.Ip,
		Port:     event.Port,
		Type:     event.Protocol,
		Date:     event.Time.Unix(),
		Data:     event.Leak.Data,
		Plugin:   event.EventSource,
		Hostname: event.Host,
		Reverse:  event.Reverse,
		Scheme:   event.Protocol,
		DatasetLeak: core.DatasetLeak{
			TotalRows:        event.Leak.Dataset.Rows,
			TotalSizeByte:    event.Leak.Dataset.Size,
			TotalCollections: event.Leak.Dataset.Collections,
			Infected:         event.Leak.Dataset.Infected,
		},
	}
	return t.jsonEncoder.Encode(hostServiceLeak)
}

func (t *TbiCoreTransformer) Name() string {
	return "tbicore"
}
