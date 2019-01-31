# go-senml

[![GitHub license](https://img.shields.io/github/license/nkristek/go-senml.svg)](https://github.com/nkristek/go-senml/blob/master/LICENSE)
[![Build Status](https://travis-ci.com/nkristek/go-senml.svg?branch=master)](https://travis-ci.com/nkristek/go-senml)
[![Coverage Status](https://coveralls.io/repos/github/nkristek/go-senml/badge.svg?branch=master)](https://coveralls.io/github/nkristek/go-senml?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/nkristek/go-senml)](https://goreportcard.com/report/github.com/nkristek/go-senml)

A go library to parse SenML records. It currently supports JSON and XML.

This library implements [RFC 8428](https://tools.ietf.org/rfc/rfc8428.txt) (SenML version 10).

## Install
```sh
go get github.com/nkristek/go-senml
```

## Import
```go
import(
	"github.com/nkristek/go-senml"
)
```

## Usage
```go
// parse using the encoding format
message, err := senml.Decode(payload, senml.JSON)
if err != nil {
	// process error
}

// resolve the message (resolve base attributes, convert relative to absolute time etc.)
resolvedMessage, err := message.Resolve()
if err != nil {
	// process error
}

// encode a new message
encodedMessage, err := message.Encode(senml.JSON)
if err != nil {
	// process error
}
```

For further documentation, visit [GoDoc](http://godoc.org/github.com/nkristek/go-senml).
