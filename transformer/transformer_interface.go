package transformer

import (
	"github.com/LeakIX/l9format"
	"io"
)

var Transformers = []TransformerInterface{
	NewJsonServiceTransformer(),
	NewUrlServiceTransformer(),
	NewHostPortTransformer(),
	NewHumanTransformer(),
	NewTbiCoreTransformer(),
	NewNmapTransformer(),
	NewMasscanTransformer(),
}

type TransformerInterface interface {
	Decode() (l9format.L9Event, error)
	Encode(event l9format.L9Event) error
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

func (t *Transformer) Close() {

}

type NoDataError struct {
	s string
}

func (e *NoDataError) Error() string {
	return e.s
}

func NewNoDataError(text string) error {
	return &NoDataError{text}
}