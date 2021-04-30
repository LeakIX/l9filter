package transformer

import (
	"encoding/json"
	"github.com/LeakIX/l9format"
)

type JsonServiceTransformer struct {
	Transformer
	jsonEncoder *json.Encoder
	jsonDecoder *json.Decoder
}

func NewJsonServiceTransformer() TransformerInterface {
	return &JsonServiceTransformer{}
}

func (t *JsonServiceTransformer) Decode(outputTransformer TransformerInterface) (err error) {
	if t.jsonDecoder == nil {
		t.jsonDecoder = json.NewDecoder(t.Reader)
	}
	for {
		event := l9format.L9Event{}
		err = t.jsonDecoder.Decode(&event)
		if err != nil {
			return err
		}
		err = outputTransformer.Encode(event)
		if err != nil {
			return err
		}
	}
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
