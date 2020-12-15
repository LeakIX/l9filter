package transformer

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/LeakIX/l9format"
	"io"
)

type HumanTransformer struct{
	Transformer
	scanner *bufio.Scanner
}


func NewHumanTransformer() TransformerInterface{
	return &HumanTransformer{}
}

func (t *HumanTransformer) Name() string {
	return "human"
}

func (t *HumanTransformer) Decode() (event l9format.L9Event, err error) {
	return event, errors.New("can't speak human yet")
}

func (t *HumanTransformer) Encode(event l9format.L9Event) error {
	_, err := io.WriteString(t.Writer, fmt.Sprintf(
		"IP: %s, PORT:%s, PROTO:%s, SSL:%t\n%.1024s\n\n",
		event.Ip, event.Port, event.Protocol, event.SSL.Enabled , event.Summary))
	if err != nil {
		return err
	}
	return nil
}