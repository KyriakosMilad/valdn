# Valdn

[![Go Reference](https://pkg.go.dev/badge/github.com/KyriakosMilad/valdn.svg)](https://pkg.go.dev/github.com/KyriakosMilad/valdn)
[![Go Report Card](https://goreportcard.com/badge/github.com/KyriakosMilad/valdn)](https://goreportcard.com/report/github.com/KyriakosMilad/valdn)
[![Build Status](https://app.travis-ci.com/KyriakosMilad/valdn.svg?branch=master)](https://app.travis-ci.com/KyriakosMilad/valdn)
[![Coverage Status](https://coveralls.io/repos/github/KyriakosMilad/valdn/badge.svg?branch=master)](https://coveralls.io/github/KyriakosMilad/valdn?branch=master)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

Valdn is a golang rich validation library. Validates request, nested JSON, nested struct, nested map, nested array, and nested
slice.

## Features

- Support all kinds.
- Support all types (even custom types).
- Validate request (application/json, multipart/form-data, application/x-www-form-urlencoded) + URL params.
- Validate nested JSON.
- Validate nested map.
- Validate nested array.
- Validate nested slice.
- Validate nested struct.
- Support using rules in struct field tag.
- +45 rules ready to use.
- +35 validation functions ready to use.
- Add custom rule.
- Add custom validation message.

## Table of Contents

<!--ts-->

* [Quick-Start](#quick-start)
* [Installation](#installation)
* [Validate single value](#validate-single-value)
* [Validate Collection](#validate-collection)
    * [Validate Struct](#validate-struct)
    * [Validate Map](#validate-map)
    * [Validate Array/Slice](#validate-arrayslice)
* [Validate JSON](#validate-json)
* [Validate Request](#validate-request)
* [Change error messages](#change-error-messages)
* [Add custom rules](#add-custom-rules)
* [Validation rules](#validation-rules)
* [Validation functions](#validation-functions)
* [Contributing](#contributing)
* [License](#license)

<!--te-->

## Quick-Start

### Validate request example:

```go
package main

import (
	"github.com/KyriakosMilad/valdn"
	"net/http"
	"encoding/json"
	"fmt"
)

func main() {
	http.HandleFunc("/", test)
	http.ListenAndServe(":8080", nil)
}

func test(w http.ResponseWriter, r *http.Request) {
	rules := valdn.Rules{}
	errs := valdn.ValidateRequest(r, rules)
	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(errs)
	} else {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Query().Get("name"))
	}
}
```

output:

```
Hello, [whatever value you will pass]
```

let's add rules

change ```rules := valdn.Rules{}``` to ```rules := valdn.Rules{"name": {"required"}}```

now if you try to pass empty name value it will fail and output

```json
{
  "name": "name is required"
}
```

let's add more rules

add ```"kind:string"``` to rules to be like this ```valdn.Rules{"name": {"required", "kind:string"}}```

now if you try to pass non-string value it will fail and output

```json
{
  "name": "name must be kind of string"
}
```

how about the lenght? let's add rules to make sure we get the value we need

add ```"minLen:3" and "maxLen:21"``` to rules to be like
this ```valdn.Rules{"name": {"required", "kind:string", "minLen:3", "maxLen:21"}}```

now if you try to pass value lower than 3 letters or greater than 21 letters it will fail and output

```json
{
  "name": "name's length must be greater than or equal: 3"
}
```

or

```json
{
  "name": "name's length must be lower than or equal: 21"
}
```

note: you can replace ```minLen:3``` and ```maxLen:21``` rules with ```lenBetween:3,21``` rule

quick and simple right? [check all the rules you can use](#validation-rules) or continue
to [discover more about valdn](#installation)

## Installation

```sh
go get "github.com/KyriakosMilad/valdn"
```

and then import it

```go
import "github.com/KyriakosMilad/valdn"
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
	err := valdn.Validate("name", name, []string{"required", "kind:string", "maxLen:5"})

	if err != nil {
		log.Fatal(err)
	}
}
```

this will output:

```
name's length must be lower than or equal: 5
```

Keep in mind when using valdn.Validate:

- It doesn't validate nested fields.
- If an error is found it will not check the rest of the rules and returns the error.
- It panics if one of the rules is not registered.

## Validate Collection

Use valdn.ValidateCollection() to validate [struct](#validate-struct), [map](#validate-map), [slice](#validate-arrayslice)
and [array](#validate-arrayslice).

valdn.ValidateCollection() takes two arguments: `value and rules (valdn.Rules{...})` and returns `valdn.Errors`

Keep in mind when using valdn.ValidateCollection:

- It panics if val is not kind of struct, map, slice or array.
- Unexported struct fields will be ignored.
- If an error is found it will not check the rest of the field's rules and continue to the next field.
- If a parent has error it's nested fields will not be validated.
- It panics if one of the rules is not registered.
- You can use * to apply rules to all direct nested fields, example:

  ``valdn.Rules{"*": "required", "Parent.*": "minLen:5"}``

### Validate Struct

Use [valdn.ValidateCollection()](#validate-collection) to validate struct.

Example:

```go
package main

import (
	"github.com/KyriakosMilad/valdn"
	"log"
)

type User struct {
	Permissions map[string]interface{}
}

func main() {
	user := User{
		Permissions: map[string]interface{}{"read": true, "write": false},
	}

	rules := valdn.Rules{"Permissions": {"required", "len:2"}, "Permissions.write": {"equal:true"}}

	errors := valdn.ValidateCollection(user, rules)

	if len(errors) > 0 {
		log.Fatal(errors)
	}
}
```

this will output:

```
Permissions.write does not equal true
```

Validate struct using struct field tag

```go
package main

import (
	"github.com/KyriakosMilad/valdn"
	"log"
)

type User struct {
	Name string `valdn:"required|maxLen:3"`
}

func main() {
	user := User{
		Name: "Ramses",
	}

	rules := valdn.Rules{}

	errors := valdn.ValidateCollection(user, rules)

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

### Validate Map

Use [valdn.ValidateCollection()](#validate-collection) to validate map.

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

	// use * to apply rules to all direct nested fields
	rules := valdn.Rules{"*": {"required", "int", "min:1"}, "Zamalek SC": {"equal:1911"}}

	errors := valdn.ValidateCollection(egyptianClubsFoundedYear, rules)

	if len(errors) > 0 {
		log.Fatal(errors)
	}
}
```

this will output:

```
Zamalek SC does not equal 1911
```

### Validate Array/Slice

Use [valdn.ValidateCollection()](#validate-collection) to validate array/slice.

Example:

```go
package main

import (
	"github.com/KyriakosMilad/valdn"
	"log"
)

func main() {
	letters := []interface{}{
		"z", // 0
		"b", // 1
		"c", // 2
		"d", // 3
	}

	// use * to apply rules to all direct nested fields
	rules := valdn.Rules{"*": {"required", "kind:string", "len:1"}, "0": {"equal:a"}}

	errors := valdn.ValidateCollection(letters, rules)

	if len(errors) > 0 {
		log.Fatal(errors)
	}
}
```

this will output:

```
0 does not equal a
```

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

	rules := valdn.Rules{"name": {"required", "kind:string"}, "value": {"required"}}

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

## Validate Request

Use valdn.ValidateRequest() to validate all Request types (application/json, multipart/form-data,
application/x-www-form-urlencoded) + URL params.

valdn.ValidateRequest() takes two arguments: `*http.Request and rules (valdn.Rules{...})` and returns `valdn.Errors`

Example:

```go
package main

import (
	"github.com/KyriakosMilad/valdn"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

func main() {
	// Create fake request for example only
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("lang=go")) // set request values: lang = go
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")            // set request header type to application/x-www-form-urlencoded

	rules := valdn.Rules{"lang": {"required", "minLen:3"}, "value": {"required"}}

	errors := valdn.ValidateRequest(r, rules)

	if len(errors) > 0 {
		log.Fatal(errors)
	}
}
```

this will output:

```
lang's length must be greater than or equal: 3
value is required
```

Keep in mind when using valdn.ValidateRequest:

- It panics if body is not compatible with header content type.
- It panics if one of the rules is not registered.
- If name has many values it will be treated as slice.
- If name has values in URL params and request body, they will be merged into one slice with that name.
- If an error is found it will not check the rest of the field's rules and continue to the next field.

## Change error messages

Use valdn.SetErrMsg() to set custom error message for a specific rule.

You can use provided parameters to dynamically set error messages:

- [name]: filed name
- [val]: field value
- [ruleVal]: rule value (rule has value like `min:value` takes float or integer as value)

Example:

```go
package main

import (
	"github.com/KyriakosMilad/valdn"
	"log"
)

func main() {
	valdn.SetErrMsg("min", "[name]'s value is [val], [name] must be greater than [ruleVal]")

	m := map[string]interface{}{"age": 15}

	rules := valdn.Rules{"age": {"min:17"}}

	errors := valdn.ValidateCollection(m, rules)

	if len(errors) > 0 {
		log.Fatal(errors)
	}
}
```

this will output:

```
age's value is 15, age must be greater than 17
```

you can use ```valdn.GetErrMsg``` to get error message

```valdn.GetErrMsg``` takes three parameters ruleName (string), ruleValue (string), fieldName (string), fieldValue (interface{})
and returns the error message

Keep in mind when using valdn.SetErrMsg and valdn.GetErrMsg:

- It panics if rule does not exist.

## Add custom rules

Use valdn.AddRule() to add custom rule. valdn.AddRule() takes three
arguments: `name (string) and ruleFunction (valdn.RuleFunc) and error message`

Example:

```go
package main

import (
	"github.com/KyriakosMilad/valdn"
	"log"
	"fmt"
	"errors"
	"reflect"
)

// Create rule to check if value starts with specific part
func startsWithRule(name string, val interface{}, ruleVal string) error {
	if v, ok := val.(string); ok {
		for i := 0; i < len(ruleVal); i++ {
			if v[i] != ruleVal[i] {
				return errors.New(valdn.GetErrMsg("startsWith", ruleVal, name, val))
			}
		}
	} else {
		panic(fmt.Errorf("startsWithRule expects value to be string, got: %v", reflect.TypeOf(val).Kind()))
	}

	return nil
}

func main() {
	valdn.AddRule("startsWith", startsWithRule, "[name] must start with '[ruleVal]'") // rule name, rule function, validation error message

	s := []interface{}{
		"newcustomrule", // 0
	}

	rules := valdn.Rules{"0": {"startsWith:test"}}

	errs := valdn.ValidateCollection(s, rules)

	if len(errs) > 0 {
		log.Fatal(errs)
	}
}
```

this will output:

```
0 must start with 'test'
```

## Validation rules

| ruleName        | ruleVal                           | Example                                                                      | Description                                                                                                                                                                                                                                                                                                                                                                         |
|-----------------|-----------------------------------|------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| required        | -                                 | required                                                                     | requiredRule checks if val exists, and it's not empty. <br /> It returns error if val is not exist or empty.                                                                                                                                                                                                                                                                        |
| kind            | string                            | kind:map                                                                     | kindRule checks if val's kind equals ruleVal. <br /> It returns error if val's kind does not equal ruleVal.                                                                                                                                                                                                                                                                         |
| notKind         | string                            | notKind:string                                                               | notKindRule checks if val's kind doesn't equal ruleVal. <br /> It returns error if val's kind equals ruleVal.                                                                                                                                                                                                                                                                       |
| kindIn          | string,string,...                 | kind:uint,uint8,uint16                                                       | kindInRule checks if val's kind is one of ruleVal[]. <br /> It returns error if val's kind is not one of ruleVal[].                                                                                                                                                                                                                                                                 |
| kindNotIn       | string,string,...                 | kind:float32,float64                                                         | kindNotInRule checks if val's kind is not one of ruleVal[]. <br /> It returns error if val's kind is one of ruleVal[].                                                                                                                                                                                                                                                              |
| type            | string                            | type:map[string]interface {}                                                 | typeRule checks if val's type equals ruleVal. <br /> It returns error if val's type does not equal ruleVal.                                                                                                                                                                                                                                                                         |
| notType         | string                            | notType:map[interface {}]string                                              | notTypeRule checks if val's type doesn't equal ruleVal. <br /> It returns error if val's type equals ruleVal.                                                                                                                                                                                                                                                                       |
| typeIn          | string                            | typeIn:map[string]int,[]int                                                  | typeInRule checks if val's type is one of ruleVal[]. <br /> It returns error if val's type is not one of ruleVal[].                                                                                                                                                                                                                                                                 |
| typeNotIn       | string                            | typeNotIn:map[int][]interface {},[4]string                                   | typeNotInRule checks if val's type is not one of ruleVal[]. <br /> It returns error if val's type is one of ruleVal[].                                                                                                                                                                                                                                                              |
| equal           | string                            | equal:[1 4 57 109]                                                           | kindRule checks if val's kind equals ruleVal. <br /> It returns error if val's kind does not equal ruleVal.                                                                                                                                                                                                                                                                         |
| int             | -                                 | int                                                                          | intRule checks if val is integer. <br /> It returns error if val is not an integer.                                                                                                                                                                                                                                                                                                 |
| uint            | -                                 | uint                                                                         | uintRule checks if val is unsigned integer. <br /> It returns error if val is not an unsigned integer.                                                                                                                                                                                                                                                                              |
| complex         | -                                 | complex                                                                      | complexRule checks if val is complex number. <br /> It returns error if val is not a complex number.                                                                                                                                                                                                                                                                                |
| float           | -                                 | float                                                                        | floatRule checks if val is float. <br /> It returns error if val is not a float.                                                                                                                                                                                                                                                                                                    |
| ufloat          | -                                 | ufloat                                                                       | ufloatRule checks if val is unsigned float. <br /> It returns error if val is not an unsigned float.                                                                                                                                                                                                                                                                                |
| numeric         | -                                 | numeric                                                                      | numericRule checks if val is numeric. <br /> It returns error if val is not a numeric.                                                                                                                                                                                                                                                                                              |
| between         | integer or float,integer or float | between:18,99                                                                | betweenRule checks if val is between min (ruleVal[0]) and max (ruleVal[1]). <br /> It panics if val is not an integer or a float. <br /> It panics if min or max is not set. <br /> It panics if min is not an integer or a float. <br /> It panics if max is not an integer or a float. <br /> It returns error if val is not between min and max.                                 |
| min             | integer or float                  | min:5                                                                        | minRule checks if val is lower than ruleVal. <br /> It panics if val is not an integer or a float. <br /> It panics if ruleVal is empty. <br /> It panics if ruleVal is not an integer or a float. <br /> It returns error if val is lower than ruleVal.                                                                                                                            |
| max             | integer or float                  | max:5                                                                        | maxRule checks if val is greater than ruleVal. <br /> It panics if val is not an integer or a float. <br /> It panics if ruleVal is empty. <br /> It panics if ruleVal is not an integer or a float. <br /> It returns error if val is greater than ruleVal.                                                                                                                        |
| in              | string,string,...                 | in:55,somestring,[1 2 3]                                                     | inRule checks if val equals one of ruleVal[] items. <br /> It returns error if val doesn't equal any item in ruleVal[].                                                                                                                                                                                                                                                             |
| notIn           | string,string,...                 | notIn:5.4,text,[1 2 3]                                                       | notInRule checks if val doesn't equal any item in ruleVal[]. <br /> It returns error if val equals one of ruleVal[] items.                                                                                                                                                                                                                                                          |
| len             | integer                           | len:7                                                                        | lenRule checks if val's length equals ruleVal. <br /> It panics if val is not array, slice, map, string, integer or float.  <br /> It returns error if val's length doesn't equal ruleVal.                                                                                                                                                                                          |
| minLen          | integer                           | minLen:3                                                                     | minLenRule checks if val's length is greater than or equal ruleVal or not. <br /> It panics if val is not array, slice, map, string, integer or float.  <br /> It returns error if val's length is lower than ruleVal.                                                                                                                                                              |
| maxLen          | integer                           | maxLen:14                                                                    | maxLenRule checks if val's length is lower than or equal ruleVal or not. <br /> It panics if val is not array, slice, map, string, integer or float.  <br /> It returns error if val's length is greater than ruleVal.                                                                                                                                                              |
| lenBetween      | integer,integer                   | lenBetween:14,19                                                             | lenBetweenRule checks if val's length is between ruleVal[0] and ruleVal[1] or not. <br /> It panics if val is not array, slice, map, string, integer or float. <br /> It panics if min or max is not set. <br /> It panics if min is not an integer. <br /> It panics if max is not an integer. <br /> It returns error if val's length is not between ruleVal[0] and ruleVal[1].   |
| lenIn           | integer,integer,...               | lenIn:1,44,190                                                               | lenInRule checks if val's length equals one of ruleVal[] items. <br /> It panics if val is not array, slice, map, string, integer or float. <br />  It panics if one of ruleVal items is not an integer. <br /> It returns error if val's length doesn't equal any item in ruleVal[].                                                                                               |
| lenNotIn        | integer,integer,...               | lenNotIn:7,389,512                                                           | lenNotInRule checks if val's length doesn't equal any item in ruleVal[]. <br /> It panics if val is not array, slice, map, string, integer or float. <br />  It panics if one of ruleVal items is not an integer. <br /> It returns error if val's length equals any item in ruleVal[].                                                                                             |
| regex           | string                            | regex:^[A-Za-z][A-Za-z0-9_]{7,29}$                                           | regexRule checks if val matches ruleVal regular expression. <br /> It panics if val is not a string. <br />  It panics if ruleVal is not a valid regular expression. <br /> It returns error if val doesn't match ruleVal regular expression.                                                                                                                                       |
| notRegex        | string                            | notRegex:^[A-Za-z][A-Za-z0-9_]{7,29}$                                        | notRegexRule checks if val doesn't match ruleVal regular expression. <br /> It panics if val is not a string. <br />  It panics if ruleVal is not a valid regular expression. <br /> It returns error if val matches ruleVal regular expression.                                                                                                                                    |
| email           | -                                 | email                                                                        | emailRule checks if val is a valid email address. <br /> It panics if val is not a string. <br />  It returns error if val is not a valid email address.                                                                                                                                                                                                                            |
| json            | -                                 | json                                                                         | jsonRule checks if val is a valid json. <br /> It panics if val is not a string. <br />  It returns error if val is not a valid json.                                                                                                                                                                                                                                               |
| ipv4            | -                                 | ipv4                                                                         | ipv4Rule checks if val is a valid IPv4. <br /> It panics if val is not a string. <br />  It returns error if val is not a valid IPv4.                                                                                                                                                                                                                                               |
| ipv6            | -                                 | ipv6                                                                         | ipv6Rule checks if val is a valid IPv6. <br /> It panics if val is not a string. <br />  It returns error if val is not a valid IPv6.                                                                                                                                                                                                                                               |
| ip              | -                                 | ip                                                                           | ipRule checks if val is a valid IP address. <br /> It panics if val is not a string. <br />  It returns error if val is not a valid IP address.                                                                                                                                                                                                                                     |
| mac             | -                                 | mac                                                                          | macRule checks if val is a valid mac address. <br /> It panics if val is not a string. <br />  It returns error if val is not a valid mac address.                                                                                                                                                                                                                                  |
| url             | -                                 | url                                                                          | urlRule checks if val is a valid URL. <br /> It panics if val is not a string. <br />  It returns error if val is not a valid URL.                                                                                                                                                                                                                                                  |
| time            | -                                 | time                                                                         | timeRule checks if val is type of time.Time. <br /> It returns error if val is not type of time.Time.                                                                                                                                                                                                                                                                               |
| timeFormat      | string                            | timeFormat:Monday, 02-Jan-06 15:04:05 MST                                    | timeFormatRule checks if val's format matches ruleVal. <br /> It returns error if val's format doesn't match ruleVal.                                                                                                                                                                                                                                                               |
| timeFormatIn    | string,string,...                 | timeFormatIn:Monday, 02-Jan-06 15:04:05 MST[]Mon, 02 Jan 2006 15:04:05 -0700 | timeFormatInRule checks if val's format matches any of ruleVal[]. <br /> Use [] to split between two formats. <br /> It returns error if val's format doesn't match any of ruleVal[].                                                                                                                                                                                               |
| timeFormatNotIn | string,string,...                 | timeFormatNotIn:02 Jan 06 15:04 MST[]02 Jan 06 15:04 -0700                   | timeFormatNotInRule checks if val's format doesn't match any of ruleVal[]. <br /> Use [] to split between two formats. <br /> It returns error if val's format matches any of ruleVal[].                                                                                                                                                                                            |
| file            | -                                 | file                                                                         | fileRule checks if val is a valid file. <br /> It returns error if val is not a valid file.                                                                                                                                                                                                                                                                                         |
| size            | integer                           | size:12000                                                                   | sizeRule checks if val's size equals ruleVal. <br /> it panics if val is not a valid file. <br /> it panics if ruleVal is not an integer. <br /> It returns error if val's size doesn't equal ruleVal.                                                                                                                                                                              |
| sizeMin         | integer                           | sizeMin:4000                                                                 | sizeMinRule checks if val's size greater than or equal ruleVal or not. <br /> it panics if val is not a valid file. <br /> it panics if ruleVal is not an integer. <br /> It returns error if val's size is lower than ruleVal.                                                                                                                                                     |
| sizeMax         | integer                           | sizeMax:20000                                                                | sizeMaxRule checks if val's size lower than or equal ruleVal or not. <br /> it panics if val is not a valid file. <br /> it panics if ruleVal is not an integer. <br /> It returns error if val's size is greater than ruleVal.                                                                                                                                                     |
| sizeBetween     | integer,integer                   | sizeBetween:2000,8000                                                        | sizeBetweenRule checks if val's size is between ruleVal[0] and ruleVal[1]. <br /> it panics if val is not a valid file. <br /> It panics if min or max is not set. <br /> It panics if min or max is not set. <br /> It panics if min is not an integer. <br /> It panics if max is not an integer. <br /> It returns error if val's size is not between ruleVal[0] and ruleVal[1]. |
| ext             | string                            | ext:zip                                                                      | extRule checks if val's extension equals ruleVal. <br /> it panics if val is not a valid file. <br /> It returns error if val's extension doesn't equal ruleVal.                                                                                                                                                                                                                    |
| notExt          | string                            | notExt:php                                                                   | notExtRule checks if val's extension does not equal ruleVal. <br /> it panics if val is not a valid file. <br /> It returns error if val's extension equals ruleVal.                                                                                                                                                                                                                |
| extIn           | string,string,...                 | extIn:jpeg,png,jpg,gif                                                       | extInRule checks if val's extension equals one of ruleVal[] items. <br /> It panics if val is not a valid file. <br />  It returns error if val's extension doesn't equal any item in ruleVal[].                                                                                                                                                                                    |
| extNotIn        | string,string,...                 | extNotIn:js,ts                                                               | extNotInRule checks if val's extension doesn't equal one of ruleVal[] items. <br /> It panics if val is not a valid file. <br />  It returns error if val's extension equals any item in ruleVal[].                                                                                                                                                                                 |

## Validation functions

| Function           | Takes                           | Returns | Description                                                                     |
|--------------------|---------------------------------|---------|---------------------------------------------------------------------------------|
| IsEmpty            | val interface{}                 | bool    | IsEmpty reports weather value is empty or not.                                  |
| IsKind             | val interface{}, kind string    | bool    | IsKind reports weather value's kind equals kind.                                |
| IsKindIn           | val interface{}, kinds []string | bool    | IsKindIn reports weather value's kind is one of kinds.                          |
| IsType             | val interface{}, typ string     | bool    | IsType reports weather value's type equals typ.                                 |
| IsTypeIn           | val interface{}, types []string | bool    | IsTypeIn reports weather value's type is one of types.                          |
| IsInteger          | val interface{}                 | bool    | IsInteger reports weather value is integer or not.                              |
| IsUnsignedInteger  | val interface{}                 | bool    | IsUnsignedInteger reports weather value is unsigned integer or not.             |
| IsFloat            | val interface{}                 | bool    | IsFloat reports weather value is float or not.                                  |
| IsUnsignedFloat    | val interface{}                 | bool    | IsUnsignedFloat reports weather value is unsigned float or not.                 |
| IsComplex          | val interface{}                 | bool    | IsComplex reports weather value is complex number or not.                       |
| IsNumeric          | val interface{}                 | bool    | IsNumeric reports weather value is numeric or not.                              |
| IsCollection       | val interface{}                 | bool    | IsCollection reports weather value's kins is one of (Array, Slice, Map) or not. |
| IsEmail            | val interface{}                 | bool    | IsEmail reports weather value is a valid email address or not.                  |
| IsJSON             | val interface{}                 | bool    | IsJSON reports weather value is a valid json or not.                            |
| IsIPv4             | val interface{}                 | bool    | IsIPv4 reports weather value is a valid IPv4 or not.                            |
| IsIPv6             | val interface{}                 | bool    | IsIPv6 reports weather value is a valid IPv6 or not.                            |
| IsIP               | val interface{}                 | bool    | IsIP reports weather value is a valid IP or not.                                |
| IsMAC              | val interface{}                 | bool    | IsMAC reports weather value is a valid MAC address or not.                      |
| IsURL              | val interface{}                 | bool    | IsURL reports weather value is a valid URL or not.                              |
| IsFile             | val interface{}                 | bool    | IsFile reports weather value is a valid file or not.                            |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Valdn is open-sourced library licensed under the [MIT license](https://opensource.org/licenses/MIT).

I'm working on the documentation.****

