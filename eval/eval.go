package eval

import (
	"context"
	"fmt"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/richardjennings/oger/inputs"
	"io"
	"log"
)

func Eval(ctx context.Context, rule []byte, path string, policy io.Reader) error {
	p := ast.NewParser().WithReader(policy)
	s, _, errs := p.Parse()
	if len(errs) > 0 {
		return fmt.Errorf("errors: %s", errs)
	}
	in, err := inputs.Inputs(s, path)
	if err != nil {
		return err
	}
	r := rego.New(
		rego.Query("data.policy.allow"),
		rego.Module("policy", string(rule)),
		rego.Input(in),
	)

	result, err := r.Eval(ctx)
	if err != nil {
		return err
	}

	if result.Allowed() {
		log.Println("ok")
		return nil
	}

	return fmt.Errorf("%s failed", path)
}
