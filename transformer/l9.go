package transformer

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/LeakIX/l9format"
)

type JsonServiceTransformer struct {
	Transformer
	reader      *bufio.Reader
	jsonEncoder *json.Encoder
}

func NewJsonServiceTransformer() TransformerInterface {
	return &JsonServiceTransformer{}
}

func (t *JsonServiceTransformer) Decode(outputTransformer TransformerInterface) (err error) {
	if t.reader == nil {
		t.reader = bufio.NewReaderSize(t.Reader, 1024*1024)
	}
	event := l9format.L9Event{}
	bytes, isPrefix, err := t.reader.ReadLine()
	if err == nil && !isPrefix {
		err = json.Unmarshal(bytes, &event)
		if err != nil {
			return err
		}
		err = outputTransformer.Encode(event)
		if err != nil {
			return err
		}
	} else if isPrefix {
		err = errors.New("line buffer overflow")
	}
	return err
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
