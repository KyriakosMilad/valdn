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

I'm working on the rest of the documentation.