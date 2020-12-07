// WIP, transformer for LeakIX legacy format

package transformer

import (
	"bufio"
	"encoding/json"
	"github.com/LeakIX/l9format"
	"gitlab.nobody.run/tbi/core"
	"io"
	"time"
)

type TbiCoreTransformer struct{
	Transformer
	scanner *bufio.Scanner
	jsonEncoder *json.Encoder
}


func NewTbiCoreTransformer() TransformerInterface{
	return &TbiCoreTransformer{}
}

func (t *TbiCoreTransformer) Decode() (event l9format.L9Event, err error) {
	if t.scanner == nil {
		t.scanner = bufio.NewScanner(t.Reader)
	}
	hostServiceLeak := &core.HostServiceLeak{}
	hostService := &core.HostService{}
	if t.scanner.Scan() {
		err = json.Unmarshal(t.scanner.Bytes(), &hostServiceLeak)
		if err != nil {
			return event, err
		}
		err = json.Unmarshal(t.scanner.Bytes(), &hostService)
		if err != nil {
			return event, err
		}
	} else {
		return event, io.EOF
	}
	event.Ip = hostService.Ip
	event.Port = hostService.Port
	event.Host = hostService.Hostname
	event.Time = time.Unix(hostService.Date, 0)
	if len(hostServiceLeak.Plugin) > 0 {
		event.EventType = "leak"
		event.Summary = hostServiceLeak.Data
		event.Leak.Dataset.Rows = hostServiceLeak.DatasetLeak.TotalRows
		event.Leak.Dataset.Size = hostServiceLeak.DatasetLeak.TotalSizeByte
	} else {
		event.EventType = "service"
		event.Summary = hostService.Data
		event.Service.Software.Name = hostService.Software.Name
		event.Service.Software.Version = hostService.Software.Version
	}
	return event, err
}

func (t *TbiCoreTransformer) Encode(event l9format.L9Event) error {
	if t.jsonEncoder == nil {
		t.jsonEncoder = json.NewEncoder(t.Writer)
	}
	return t.jsonEncoder.Encode(event)
}

func (t *TbiCoreTransformer) Name() string {
	return "tbicore"
}