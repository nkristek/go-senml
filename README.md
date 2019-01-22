# go-senml

[![GitHub license](https://img.shields.io/github/license/nkristek/go-senml.svg)](https://github.com/nkristek/go-senml/blob/master/LICENSE)
[![Build Status](https://travis-ci.com/nkristek/go-senml.svg?branch=master)](https://travis-ci.com/nkristek/go-senml)
[![Coverage Status](https://coveralls.io/repos/github/nkristek/go-senml/badge.svg?branch=master)](https://coveralls.io/github/nkristek/go-senml?branch=master)

A go library to parse SenML records. It currently supports only JSON and XML, but other formats like CBOR and EXI are planned.

This library implements [RFC8428](https://tools.ietf.org/rfc/rfc8428.txt).

## Install
```sh
go get github.com/nkristek/go-seml
```

## Import
```go
import(
	"github.com/nkristek/go-seml"
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

For further documentation, visit http://godoc.org/github.com/nkristek/go-senml