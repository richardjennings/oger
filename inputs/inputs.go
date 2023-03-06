package inputs

import (
	"encoding/json"
	"github.com/open-policy-agent/opa/ast"
	"reflect"
)

func Inputs(stmts []ast.Statement, path string) (map[string]interface{}, error) {
	r := make(map[string]interface{})
	var err error
	r["statements"], err = statements(stmts)
	r["path"] = path
	r["ext_calls"] = extCalls(stmts)

	return r, err
}

// Statements creates a Rego AST representation to evaluate in policy
func statements(ss []ast.Statement) ([]map[string]interface{}, error) {
	var r []map[string]interface{}
	var err error
	var stmt map[string]interface{}
	for _, v := range ss {
		stmt, err = statement(v)
		stmt["_type"] = reflect.TypeOf(v).String()
		r = append(r, stmt)
	}
	return r, err
}

func statement(s ast.Statement) (map[string]interface{}, error) {
	r := make(map[string]interface{})
	var err error
	var b []byte
	b, err = json.Marshal(s)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &r)
	return r, err
}

func extCalls(ss []ast.Statement) []string {
	var calls []string
	visitor := ast.NewGenericVisitor(func(x interface{}) bool {
		switch x := x.(type) {
		case ast.Call:
			calls = append(calls, x[0].Value.String())
		}
		return false
	})
	for _, v := range ss {
		visitor.Walk(v)
	}

	// find all function declarations in the ast.
	// extCalls are calls to functions not defined in this AST

	return calls
}
