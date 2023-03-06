package cmd

import (
	"context"
	"github.com/richardjennings/oger/eval"
	"github.com/spf13/cobra"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var checkCmd = &cobra.Command{
	Use:  "check <path/to/rule> </path/to/policies>",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		pf, err := os.Open(args[0])
		if err != nil {
			log.Fatalln(err)
		}

		pc, err := io.ReadAll(pf)
		if err != nil {
			log.Fatalln(err)
		}

		if err := filepath.Walk(os.Args[2], func(path string, info fs.FileInfo, err error) error {
			if filepath.Ext(path) == ".rego" {
				log.Printf("processing %s\n", path)
				f, err := os.Open(path)
				if err != nil {
					return err
				}
				return eval.Eval(ctx, pc, path, f)
			}
			return nil
		}); err != nil {
			log.Fatalln(err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
