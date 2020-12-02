package transformer

import (
	"gitlab.nobody.run/tbi/core"
	"io"
)

var Transformers = []TransformerInterface{
	NewJsonServiceTransformer(),
	NewUrlServiceTransformer(),
	NewHostPortTransformer(),
	NewHumanTransformer(),
}

type TransformerInterface interface {
	Decode() (core.HostService, error)
	Encode(hostService core.HostService) error
	Name() string
	SetReader(reader io.Reader)
	SetWriter(writer io.Writer)
}

type Transformer struct {
	Reader io.Reader
	Writer io.Writer
}

func (t *Transformer) SetReader(reader io.Reader)  {
	t.Reader = reader
}


func (t *Transformer) SetWriter(writer io.Writer) {
	t.Writer = writer
}