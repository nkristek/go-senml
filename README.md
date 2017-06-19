# go-senml

A go library to parse SenML records. It currently supports only JSON and XML, but more formats like CBOR are planned.

This currently implements 'draft-ietf-core-senml-08' (https://tools.ietf.org/html/draft-ietf-core-senml-08)

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

## Use
```go
// parse using the encoding format
message, err := senml.ParseBytes(payload, senml.JSON)
if err != nil {
	// process error
}

// resolve the message (resolve base attributes, convert relative to absolute time etc.)
resolvedMessage, err := senml.Resolve(message)
if err != nil {
	// process error
}

// encode a new message
encodedMessage, err := message.EncodeToBytes(senml.JSON)
if err != nil {
	// process error
}
```

For further documentation, visit http://godoc.org/github.com/nkristek/go-senml

