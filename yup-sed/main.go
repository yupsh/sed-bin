package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	gloo "github.com/gloo-foo/framework"
	. "github.com/yupsh/sed"
)

const (
	flagExpression    = "expression"
	flagScriptFile    = "file"
	flagInPlace       = "in-place"
	flagQuiet         = "quiet"
	flagExtendedRegex = "regexp-extended"
)

func main() {
	app := &cli.App{
		Name:  "sed",
		Usage: "stream editor for filtering and transforming text",
		UsageText: `sed [OPTIONS] [SCRIPT] [FILE...]

   sed is a stream editor. A stream editor is used to perform basic text
   transformations on an input stream (a file or input from a pipeline).`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    flagExpression,
				Aliases: []string{"e"},
				Usage:   "add the script to the commands to be executed",
			},
			&cli.StringFlag{
				Name:    flagScriptFile,
				Aliases: []string{"f"},
				Usage:   "add the contents of script-file to the commands to be executed",
			},
			&cli.BoolFlag{
				Name:    flagInPlace,
				Aliases: []string{"i"},
				Usage:   "edit files in place",
			},
			&cli.BoolFlag{
				Name:    flagQuiet,
				Aliases: []string{"n", "silent"},
				Usage:   "suppress automatic printing of pattern space",
			},
			&cli.BoolFlag{
				Name:    flagExtendedRegex,
				Aliases: []string{"r", "E"},
				Usage:   "use extended regular expressions in the script",
			},
		},
		Action: action,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "sed: %v\n", err)
		os.Exit(1)
	}
}

func action(c *cli.Context) error {
	var params []any

	// First arg might be the expression if -e not used
	if !c.IsSet(flagExpression) && c.NArg() > 0 {
		params = append(params, Expression(c.Args().Get(0)))
		// Remaining args are files
		for i := 1; i < c.NArg(); i++ {
			params = append(params, gloo.File(c.Args().Get(i)))
		}
	} else {
		// All args are files
		for i := 0; i < c.NArg(); i++ {
			params = append(params, gloo.File(c.Args().Get(i)))
		}
	}

	// Add flags based on CLI options
	if c.IsSet(flagExpression) {
		params = append(params, Expression(c.String(flagExpression)))
	}
	if c.IsSet(flagScriptFile) {
		params = append(params, ScriptFile(c.String(flagScriptFile)))
	}
	if c.Bool(flagInPlace) {
		params = append(params, InPlace)
	}
	if c.Bool(flagQuiet) {
		params = append(params, Quiet)
	}
	if c.Bool(flagExtendedRegex) {
		params = append(params, ExtendedRegex)
	}

	// Create and execute the sed command
	cmd := Sed(params...)
	return gloo.Run(cmd)
}
