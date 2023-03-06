package eval

import (
	"bytes"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEval(t *testing.T) {
	type tc struct {
		rule       string
		policy     string
		policyPath string
		err        error
	}
	for _, test := range []tc{

		{
			// first statement must be of type *ast.Package
			rule: `
		package policy
		default allow = false
		allow {
		    input.statements[0]["_type"] == "*ast.Package"
		}`,
			policy: `package kubernetes.admission`,
			err:    nil,
		},

		// first statement must be of type *ast.
		{
			rule: `
		package policy
		default allow = false
		allow {
		    input.statements[0]["_type"] == "*ast.package"
		}`,
			policy:     `package kubernetes.admission`,
			policyPath: "some/package.rego",
			err:        errors.New("some/package.rego failed"),
		},

		// Allow use of http.send call only by comparing ext_calls to allowed list
		{
			rule: `
		package policy
		default allow = false
		valid_ext_calls := { "http.send" }
		invalid_set := {x | x := input.ext_calls[_]} 
		invalid := invalid_set - valid_ext_calls
		allow {
			count(invalid) == 0
		}`,
			policy: `
		package play
		default hello = false
		hello {
			response := http.send({
				"method" : "GET",
				"url": "http://localhost:8181/v1/data/example"
			})
		}`,
			policyPath: "ext_calls.rego",
			err:        nil,
		},

		// Deny use of http.send call by comparing ext_calls to allowed list
		{
			rule: `
		package policy
		default allow = false
		valid_ext_calls := { "abc.def" }
		invalid_set := {x | x := input.ext_calls[_]} 
		invalid := invalid_set - valid_ext_calls
		allow {
			count(invalid) == 0
		}`,
			policy: `
		package play
		default hello = false
		hello {
			response := http.send({
				"method" : "GET",
				"url": "http://localhost:8181/v1/data/example"
			})
		}`,
			policyPath: "ext_calls.rego",
			err:        errors.New("ext_calls.rego failed"),
		},
	} {
		assert.Equal(
			t,
			test.err,
			Eval(context.Background(), []byte(test.rule), test.policyPath, bytes.NewBufferString(test.policy)),
		)
	}

}
