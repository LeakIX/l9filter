package transformer

import (
	"encoding/json"
	"github.com/LeakIX/l9format"
	"strconv"
	"time"
)

type SxTransformer struct {
	Transformer
	jsonEncoder *json.Encoder
	jsonDecoder *json.Decoder
}

func NewSxTransformer() TransformerInterface {
	return &SxTransformer{}
}

func (t *SxTransformer) Decode(outputTransformer TransformerInterface) (err error) {
	if t.jsonDecoder == nil {
		t.jsonDecoder = json.NewDecoder(t.Reader)
	}
	sxLine := SxResult{}
	err = t.jsonDecoder.Decode(&sxLine)
	if err != nil {
		return err
	}
	if len(sxLine.ScanType) < 1 {
		sxLine.ScanType = "arpscan"
	}
	return outputTransformer.Encode(l9format.L9Event{
		EventType:     sxLine.ScanType,
		Ip:            sxLine.Ip,
		Time:          time.Now(),
		EventSource:   "sx-" + sxLine.ScanType,
		EventPipeline: []string{"sx-" + sxLine.ScanType},
		Port:          strconv.Itoa(sxLine.Port),
		Vendor:        sxLine.Vendor,
	})
}

func (t *SxTransformer) Encode(event l9format.L9Event) error {
	if t.jsonEncoder == nil {
		t.jsonEncoder = json.NewEncoder(t.Writer)
	}
	port, err := strconv.Atoi(event.Port)
	if err != nil {
		return err
	}
	return t.jsonEncoder.Encode(SxResult{
		Ip:       event.Ip,
		Mac:      event.Mac,
		Vendor:   event.Vendor,
		Port:     port,
		ScanType: event.EventType,
	})
}

func (t *SxTransformer) Name() string {
	return "sx"
}

type SxResult struct {
	Ip       string `json:"ip"`
	Mac      string `json:"mac"`
	Vendor   string `json:"vendor"`
	Port     int    `json:"port"'`
	ScanType string `json:"scan"`
}
