package main

import (
	"github.com/LeakIX/l9filter"
	"github.com/alecthomas/kong"
)

var App struct {
	Transform l9filter.TransformCommand `cmd help:"Takes input and transforms it " default:"1"`
}

func main() {
	ctx := kong.Parse(&App)
	// Call the Run() method of the selected parsed command.
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
