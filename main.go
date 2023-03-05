package main

import (
	"context"
	"github.com/richardjennings/oger/eval"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalln("usage: oger <path/to/rule> </path/to/policies>")
	}

	ctx := context.Background()

	pf, err := os.Open(os.Args[1])
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

}
