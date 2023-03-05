package eval

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEval(t *testing.T) {
	rule := []byte(`
package policy

default allow = false

# package must be specified as the first statement
allow {
    input.statements[0]["_type"] == "*ast.Package"
}

`)
	policy := bytes.NewBufferString(`
package kubernetes.admission
`)
	assert.Nil(t, Eval(context.Background(), rule, "kubernetes/admission.rego", policy))
}
