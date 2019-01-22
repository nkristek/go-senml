/*
	Unit tests, run:
	`go test`
	to execute them.
*/
package senml_test

import (
	"github.com/nkristek/go-senml"
	"testing"
)

func TestJSONParsing(t *testing.T) {
	var testData string = `[
     {"bn":"urn:dev:ow:10e2073a0108006:","bt":1.276020076001e+09,
      "bu":"A","bver":5,
      "n":"voltage","u":"V","v":120.1},
     {"n":"current","t":-5,"v":1.2},
     {"n":"current","t":-4,"v":1.3},
     {"n":"current","t":-3,"v":1.4},
     {"n":"current","t":-2,"v":1.5},
     {"n":"current","t":-1,"v":1.6},
     {"n":"current","v":1.7}
   ]`

	message, err := senml.Decode([]byte(testData), senml.JSON)
	if err != nil {
		t.Error("parsing initial JSON failed: ", err)
		return
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("resolving message failed: ", err)
		return
	}

	_, err = message.Encode(senml.JSON)
	if err != nil {
		t.Error("encoding message to JSON failed: ", err)
		return
	}

	_, err = resolvedMessage.Encode(senml.JSON)
	if err != nil {
		t.Error("encoding resolved message to JSON failed: ", err)
		return
	}
}

func TestXMLParsing(t *testing.T) {
	var testData string = `<sensml xmlns="urn:ietf:params:xml:ns:senml">
     <senml bn="urn:dev:ow:10e2073a0108006:" bt="1.276020076001e+09"
     bu="A" bver="5" n="voltage" u="V" v="120.1"></senml>
     <senml n="current" t="-5" v="1.2"></senml>
     <senml n="current" t="-4" v="1.3"></senml>
     <senml n="current" t="-3" v="1.4"></senml>
     <senml n="current" t="-2" v="1.5"></senml>
     <senml n="current" t="-1" v="1.6"></senml>
     <senml n="current" v="1.7"></senml>
   </sensml>`

	message, err := senml.Decode([]byte(testData), senml.XML)
	if err != nil {
		t.Error("parsing intial XML failed: ", err)
		return
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("resolving message failed: ", err)
		return
	}

	_, err = message.Encode(senml.XML)
	if err != nil {
		t.Error("encoding message to XML failed: ", err)
		return
	}

	_, err = resolvedMessage.Encode(senml.XML)
	if err != nil {
		t.Error("encoding resolved message to XML failed: ", err)
		return
	}
}
