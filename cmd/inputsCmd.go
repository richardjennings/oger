package cmd

import (
	"encoding/json"
	"github.com/open-policy-agent/opa/ast"
	"github.com/richardjennings/oger/inputs"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var inputsCmd = &cobra.Command{
	Use:  "inputs </path/to/policy>",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("processing %s\n", args[0])
		f, err := os.Open(args[0])
		if err != nil {
			return err
		}
		p := ast.NewParser().WithReader(f)
		s, _, errs := p.Parse()
		if len(errs) != 0 {
			return errs
		}
		in, err := inputs.Inputs(s, args[0])
		if err != nil {
			return err
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "    ")
		if err := enc.Encode(in); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(inputsCmd)
}
