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
	return r, err
}

// Statements creates a Rego AST representation to evaluate in policy
func statements(ss []ast.Statement) ([]map[string]interface{}, error) {
	var r []map[string]interface{}
	var err error
	stmt := make(map[string]interface{})
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
	err = json.Unmarshal(b, &r)
	return r, err
}
