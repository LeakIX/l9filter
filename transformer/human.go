package transformer

import (
	"bufio"
	"errors"
	"fmt"
	"gitlab.nobody.run/tbi/core"
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

func (t *HumanTransformer) Decode() (hostService core.HostService, err error) {
	return hostService, errors.New("can't speak human yet")
}

func (t *HumanTransformer) Encode(hostService core.HostService) error {
	_, err := io.WriteString(t.Writer, fmt.Sprintf(
		"IP: %s, PORT:%s, TYPE:%s, SCHEME:%s\n%s\n\n",
		hostService.Ip, hostService.Port, hostService.Type, hostService.Scheme, hostService.Data))
	if err != nil {
		return err
	}
	return nil
}