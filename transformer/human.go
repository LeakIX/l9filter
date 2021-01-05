package transformer

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/LeakIX/l9format"
	"io"
)

type HumanTransformer struct {
	Transformer
	scanner *bufio.Scanner
}

func NewHumanTransformer() TransformerInterface {
	return &HumanTransformer{}
}

func (t *HumanTransformer) Name() string {
	return "human"
}

func (t *HumanTransformer) Decode(outputTransformer TransformerInterface) (err error) {
	return errors.New("can't speak human yet")
}

func (t *HumanTransformer) Encode(event l9format.L9Event) error {
	_, err := io.WriteString(t.Writer, fmt.Sprintf(
		"Found %s at %s (%s) on port %s PROTO:%s SSL:%t\n%.1024s\n\n",
		event.EventType, event.Ip, event.Host, event.Port, event.Protocol, event.SSL.Enabled, event.Summary))
	if err != nil {
		return err
	}
	return nil
}
