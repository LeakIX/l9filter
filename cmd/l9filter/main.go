package main

import (
	"github.com/LeakIX/l9filter/transformer"
	"github.com/alecthomas/kong"
	"io"
	"log"
	"os"
)

var App struct {
	Filter FilterCommand `cmd help:"Takes input, filters output" default:"1"`
}

type FilterCommand struct {
	InputFormat string `required help:"input format" short:"i"`
	OutputFormat string `required help:"output format" short:"o"`
	SourceFile string `help:"Input file, stdin if none" short:"s"`
	TargetFile string `help:"Output file, stdout if none" short:"f"`
	PortFilter string `help:"Filter on port" short:"p"`
	TypeFilter string `help:"Filter on type" short:"t"`
	InputTransformer transformer.TransformerInterface `kong:"-"`
	OutputTransformer transformer.TransformerInterface `kong:"-"`
	LogWriter io.Writer `kong:"-"`
}

func main() {
	ctx := kong.Parse(&App)
	// Call the Run() method of the selected parsed command.
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}

func (cmd *FilterCommand) Run() error {
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
			log.Println("selected input : " + trs.Name())
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
			log.Println("selected input :  " + trs.Name())
			cmd.OutputTransformer = trs
		}
	}
	for {
		hostService, err := cmd.InputTransformer.Decode()
		if  err != nil {
			return err
		}
		if len(cmd.TypeFilter) > 0 {
			if hostService.Type == cmd.TypeFilter {
				continue
			}
		}
		if len(cmd.PortFilter) > 0 {
			if hostService.Port == cmd.PortFilter {
				continue
			}
		}
		err = cmd.OutputTransformer.Encode(hostService)
		if err != nil {
			return err
		}
	}
}
