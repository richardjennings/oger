# Oger

## What
Evaluating Static Analysis For Rego using Rego

## How
Rego Policy AST properties are added as inputs which may then be evaluated by a rule policy.

The inputs generated for a particular Rego policy can be viewed with `oger inputs <path to policy.rego>`,

For the following Rego
```
package play

default hello = false

hello {
    response := http.send({
        "method" : "GET",
        "url": "http://localhost:8181/v1/data/example"
    })
}
```
The following inputs are generated (abbreviated):

```
{
    "ext_calls": [
        "http.send"
    ],
    "path": "/Users/richardjennings/tmp/policy/test.rego",
    "statements": [
        {
            "_type": "*ast.Package",
            "path": [
                {
                    "type": "var",
                    "value": "data"
                },
                {
                    "type": "string",
                    "value": "play"
                }
            ]
        },
        {
...
```

A rule can be evaluated against this Rego policy, the following allows only an abc.def function call for example:
```
package policy
default allow = false
valid_ext_calls := { "abc.def" }
invalid_set := {x | x := input.ext_calls[_]} 
invalid := invalid_set - valid_ext_calls
allow {
    count(invalid) == 0
}
```

Executing using `oger rule.rego /policy/` returns `/tmp/rule.rego failed`.

