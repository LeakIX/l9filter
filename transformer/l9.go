package transformer

import (
	"bufio"
	"encoding/json"
	"github.com/LeakIX/l9format"
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

func (t *JsonServiceTransformer) Decode() (event l9format.L9Event, err error) {
	if t.scanner == nil {
		t.scanner = bufio.NewScanner(t.Reader)
	}
	if t.scanner.Scan() {
		err = json.Unmarshal(t.scanner.Bytes(), &event)
	} else {
		return event, io.EOF
	}
	return event, err
}

func (t *JsonServiceTransformer) Encode(event l9format.L9Event) error {
	if t.jsonEncoder == nil {
		t.jsonEncoder = json.NewEncoder(t.Writer)
	}
	return t.jsonEncoder.Encode(event)
}

func (t *JsonServiceTransformer) Name() string {
	return "l9"
}