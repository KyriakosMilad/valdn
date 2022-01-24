# Valdn

_Everything you need to validate data in go._

[![Go Reference](https://pkg.go.dev/badge/github.com/KyriakosMilad/valdn.svg)](https://pkg.go.dev/github.com/KyriakosMilad/valdn)
[![Go Report Card](https://goreportcard.com/badge/github.com/KyriakosMilad/valdn)](https://goreportcard.com/report/github.com/KyriakosMilad/valdn)
[![Build Status](https://app.travis-ci.com/KyriakosMilad/valdn.svg?branch=master)](https://app.travis-ci.com/KyriakosMilad/valdn)
[![Coverage Status](https://coveralls.io/repos/github/KyriakosMilad/valdn/badge.svg?branch=master)](https://coveralls.io/github/KyriakosMilad/valdn?branch=master)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

Valdn is a golang cross-validation library. Validates request, nested JSON, nested struct, nested map, and nested slice. Validates
any other Kind as a non-nested value.

## Table of Contents

<!--ts-->
* [Installation](#installation)
* [Features](#features)
* [Validate Single Value](#validate-single-value)
* [Validate Struct](#validate-struct)
* [Validate Map](#validate-map)
* [Validate Slice](#validate-slice)
* [Validate JSON](#validate-json)
<!--te-->

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
	Name  string `valdn:"required"`
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

## Validate Single Value

Use valdn.Validate() to validate one single value.

valdn.Validate takes three arguments: `name, value, and slice of rules ([]string{...})` and returns `error`

Example:

```go
package main

import (
	"github.com/KyriakosMilad/valdn"
	"log"
)

func main() {
	name := "Narmer"
	err := valdn.Validate("name", name, []string{"required", "kind:string", "minLen:6"})

	if err != nil {
		log.Fatal(err)
	}
}
```

this will output:

```
name's length must be greater than or equal: 6
```

Keep in mind when using valdn.Validate:

- It doesn't validate nested fields.
- If an error is found it will not check the rest of the rules and returns the error.
- It panics if one of the rules is not registered.

## Validate Struct

Use valdn.ValidateStruct() to validate struct.

valdn.ValidateStruct() takes two arguments: `value and rules (valdn.Rules{...})` and returns `valdn.Errors`

Example:

```go
package main

import (
	"github.com/KyriakosMilad/valdn"
	"log"
)

// struct and it's nested fields must be exported, so it can be accessed by valdn
type User struct {
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

this will output:

```
Roles.write does not equal true
```

Validate struct using struct field tag

```go
package main

import (
	"github.com/KyriakosMilad/valdn"
	"log"
)

// struct and it's nested fields must be exported, so it can be accessed by valdn
type User struct {
	Name string `valdn:"required|maxLen:3"`
}

func main() {
	user := User{
		Name: "Ramses",
	}

	rules := valdn.Rules{}

	errors := valdn.ValidateStruct(user, rules)

	if len(errors) > 0 {
		log.Fatal(errors)
	}
}
```

this will output:

```
Name's length must be lower than or equal: 3
```

You can change the TagName and Separator used to identify rules in struct field tag:

`valdn.TagName = "valdn"`

`valdn.Separator = "|"`

Keep in mind when using valdn.ValidateStruct:

- It panics if val is not kind of struct.
- It panics if val is not exported, or it's fields is not exported.
- It panics if val is not a struct.
- If an error is found it will not check the rest of the field's rules and continue to the next field.
- If a parent has error it's nested fields will not be validated.
- It panics if one of the rules is not registered.
- It panics if one of the nested fields is a map and it's type is not map[string]interface{}.
- It panics if one of the nested fields is a slice and it's type is not []interface{}.

## Validate Map

Use valdn.ValidateMap() to validate map.

valdn.ValidateMap() takes two arguments: `value and rules (valdn.Rules{...})` and returns `valdn.Errors`

Example:

```go
package main

import (
	"github.com/KyriakosMilad/valdn"
	"log"
)

func main() {
	egyptianClubsFoundedYear := map[string]interface{}{
		"Zamalek SC":         1811,
		"Al Ahly SC":         1907,
		"Ismailly SC":        1924,
		"Al Masry SC":        1920,
		"Ittihad of Alex SC": 1914,
	}

	// use * to apply rules to all nested fields
	rules := valdn.Rules{"*": {"required", "numerical", "min:1"}, "Zamalek SC": {"equal:1911"}}

	errors := valdn.ValidateMap(egyptianClubsFoundedYear, rules)

	if len(errors) > 0 {
		log.Fatal(errors)
	}
}
```

this will output:

```
Zamalek SC does not equal 1911
```

Keep in mind when using valdn.ValidateMap:

- If an error is found it will not check the rest of the field's rules and continue to the next field.
- If a parent has error it's nested fields will not be validated.
- It panics if one of the rules is not registered.
- It panics if one of the nested fields is a map and it's type is not map[string]interface{}.
- It panics if one of the nested fields is a slice and it's type is not []interface{}.

## Validate Slice

Use valdn.ValidateSlice() to validate slice.

valdn.ValidateSlice() takes two arguments: `value and rules (valdn.Rules{...})` and returns `valdn.Errors`

Example:

```go
package main

import (
	"github.com/KyriakosMilad/valdn"
	"log"
)

func main() {
	letters := []interface{}{
		"a", // 0
		"b", // 1
		"c", // 2
		"d", // 3
	}

	// use * to apply rules to all nested fields
	rules := valdn.Rules{"*": {"required", "kind:string", "len:1"}, "0": {"equal:a"}}

	errors := valdn.ValidateSlice(letters, rules)

	if len(errors) > 0 {
		log.Fatal(errors)
	}
}
```

this will output:

```
0 does not equal a
```

Keep in mind when using valdn.ValidateSlice:

- If an error is found it will not check the rest of the field's rules and continue to the next field.
- If a parent has error it's nested fields will not be validated.
- It panics if one of the rules is not registered.
- It panics if one of the nested fields is a map and it's type is not map[string]interface{}.
- It panics if one of the nested fields is a slice and it's type is not []interface{}.

## Validate JSON

Use valdn.ValidateJSON() to validate JSON.

valdn.ValidateJSON() takes two arguments: `value and rules (valdn.Rules{...})` and returns `valdn.Errors`

Example:

```go
package main

import (
	"github.com/KyriakosMilad/valdn"
	"log"
)

func main() {
	stringJSON := `{"name":11}`

	// use * to apply rules to all nested fields
	rules := valdn.Rules{"type": {"required", "kind:string"}, "value": {"required"}}

	errors := valdn.ValidateJSON(stringJSON, rules)

	if len(errors) > 0 {
		log.Fatal(errors)
	}
}
```

this will output:

```
name must be kind of string
value is required
```

Keep in mind when using valdn.ValidateJSON:

- It panics if val is not JSON.
- If an error is found it will not check the rest of the field's rules and continue to the next field.
- If parent has error it's nested fields will not be validated.
- It panics if one of the rules is not registered.

I'm working on the rest of the documentation.****
