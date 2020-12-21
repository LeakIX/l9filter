package l9filter

import (
	"errors"
	"github.com/LeakIX/l9filter/transformer"
	"io"
	"log"
	"os"
)

type TransformCommand struct {
	InputFormat       string                           `required help:"input format" short:"i"`
	OutputFormat      string                           `required help:"output format" short:"o"`
	SourceFile        string                           `help:"Input file, stdin if none" short:"s"`
	TargetFile        string                           `help:"Output file, stdout if none" short:"f"`
	PortFilter        string                           `help:"Filter on port" short:"p"`
	TypeFilter        string                           `help:"Filter on type" short:"t"`
	InputTransformer  transformer.TransformerInterface `kong:"-"`
	OutputTransformer transformer.TransformerInterface `kong:"-"`
	LogWriter         io.Writer                        `kong:"-"`
}

func (cmd *TransformCommand) Run() error {
	for _, trs := range transformer.Transformers {
		if cmd.InputFormat == trs.Name() {
			trs.SetReader(os.Stdin)
			if len(cmd.SourceFile) > 0 {
				inputFile, err := os.Open(cmd.SourceFile)
				if err != nil {
					return err
				}
				trs.SetReader(inputFile)
			}
			log.Println("selected input: " + trs.Name())
			cmd.InputTransformer = trs
		}
		if cmd.OutputFormat == trs.Name() {
			trs.SetWriter(os.Stdout)
			if len(cmd.TargetFile) > 0 {
				_, err := os.Stat(cmd.TargetFile)
				if err == nil {
					return os.ErrExist
				}
				outputFile, err := os.Create(cmd.TargetFile)
				if err != nil {
					return err
				}
				trs.SetWriter(outputFile)
			}
			log.Println("selected output: " + trs.Name())
			cmd.OutputTransformer = trs
		}
	}
	if cmd.InputTransformer == nil {
		return errors.New("input format doesn't exists")
	}
	if cmd.OutputTransformer == nil {
		return errors.New("output format doesn't exists")
	}
	for {
		event, err := cmd.InputTransformer.Decode()
		if err != nil {
			if _, isNoDataError := err.(*transformer.NoDataError); isNoDataError {
				// Happens when we meet comments ect ... we can safely skip
				continue
			}
			if err.Error() == "EOF" {
				// We done with input, close output and leave
				os.Exit(0)
			}
			return err
		}
		err = cmd.OutputTransformer.Encode(event)
		if err != nil {
			return err
		}
	}
}
