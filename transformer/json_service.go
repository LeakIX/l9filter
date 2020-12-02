package transformer

import (
	"bufio"
	"encoding/json"
	"gitlab.nobody.run/tbi/core"
	"io"
)

type JsonServiceTransformer struct{
	Transformer
	scanner *bufio.Scanner
	jsonEncoder *json.Encoder
}


func NewJsonServiceTransformer() TransformerInterface{
	return &JsonServiceTransformer{}
}

func (t *JsonServiceTransformer) Decode() (hostService core.HostService, err error) {
	if t.scanner == nil {
		t.scanner = bufio.NewScanner(t.Reader)
	}
	if t.scanner.Scan() {
		err = json.Unmarshal(t.scanner.Bytes(), &hostService)
	} else {
		return hostService, io.EOF
	}
	return hostService, err
}

func (t *JsonServiceTransformer) Encode(hostService core.HostService) error {
	if t.jsonEncoder == nil {
		t.jsonEncoder = json.NewEncoder(t.Writer)
	}
	return t.jsonEncoder.Encode(hostService)
}

func (t *JsonServiceTransformer) Name() string {
	return "json"
}