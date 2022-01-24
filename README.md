# Valdn
 _Everything you need to validate data in go._

[![Go Reference](https://pkg.go.dev/badge/github.com/KyriakosMilad/valdn.svg)](https://pkg.go.dev/github.com/KyriakosMilad/valdn)
[![Go Report Card](https://goreportcard.com/badge/github.com/KyriakosMilad/valdn)](https://goreportcard.com/report/github.com/KyriakosMilad/valdn)
[![Build Status](https://app.travis-ci.com/KyriakosMilad/valdn.svg?branch=master)](https://app.travis-ci.com/KyriakosMilad/valdn)
[![Coverage Status](https://coveralls.io/repos/github/KyriakosMilad/valdn/badge.svg?branch=master)](https://coveralls.io/github/KyriakosMilad/valdn?branch=master)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

Valdn is a golang cross-validation library. Validates request, nested JSON, nested struct, nested map, and nested slice. Validates any other Kind as a non-nested value.

## Features

- Validate all kinds.
- Validate request (application/json, multipart/form-data, application/x-www-form-urlencoded) + URL params.
- Validate nested json.
- Validate nested map.
- Validate nested slice.
- Validate nested struct.
- Support using rules in struct field tag.
- +45 rules ready to use.
- +35 checker functions ready to use.
- Add custom rule.
- Add custom validation message.

## Installation

```sh
go get "github.com/KyriakosMilad/valdn"
```

## Quick-Start

Validate single value:
```go
package main

import (
	"github.com/KyriakosMilad/valdn"
	"log"
)

func main() {
	name := "valdn"
	err := valdn.Validate("name", name, []string{"required", "kind:string", "minLen:6"})

	if err != nil {
		log.Fatal(err)
	}
}
```

output:

```
name's length must be greater than or equal: 6
```

Validate nested value:
```go
package main

import (
	"github.com/KyriakosMilad/valdn"
	"log"
)

type User struct {
	Name string `valdn:"required"`
	Roles map[string]interface{}
}

func main() {
        user := User{
            Roles: map[string]interface{}{"read": true, "write": false},
        }

        rules := valdn.Rules{"Roles": {"required", "len:2"}, "Roles.write": {"equal:true"}}

        errors := valdn.ValidateStruct(user, rules)

        if len(errors) > 0 {
            log.Fatal(errors)
        }
}
```

output:

```
Name is required
Roles.write does not equal true
```

I'm working on the rest of the documentation.****
