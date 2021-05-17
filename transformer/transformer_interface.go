package transformer

import (
	"github.com/LeakIX/l9format"
	"io"
)

var Transformers = []TransformerInterface{
	NewJsonServiceTransformer(),
	NewDnsxTransformer(),
	NewUrlServiceTransformer(),
	NewHostPortTransformer(),
	NewHumanTransformer(),
	NewTbiCoreTransformer(),
	NewNmapTransformer(),
	NewMasscanTransformer(),
	NewSxTransformer(),
}

type TransformerInterface interface {
	Decode(outputTransformer TransformerInterface) error
	Encode(event l9format.L9Event) error
	Name() string
	SetReader(reader io.Reader)
	SetWriter(writer io.Writer)
}

var L9Sources []string

type Transformer struct {
	Reader io.Reader
	Writer io.Writer
}

func (t *Transformer) SetReader(reader io.Reader) {
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
