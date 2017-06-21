/*
	Unit tests, run:
	`go test`
	to execute them.
*/
package senml

import (
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

	message, err := Decode([]byte(testData), JSON)
	if err != nil {
		t.Fatalf("parsing initial JSON failed: ", err)
		return
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Fatalf("resolving message failed: ", err)
		return
	}

	_, err = message.Encode(JSON)
	if err != nil {
		t.Fatalf("encoding message to JSON failed: ", err)
		return
	}

	_, err = resolvedMessage.Encode(JSON)
	if err != nil {
		t.Fatalf("encoding resolved message to JSON failed: ", err)
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

	message, err := Decode([]byte(testData), XML)
	if err != nil {
		t.Fatalf("parsing intial XML failed: ", err)
		return
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Fatalf("resolving message failed: ", err)
		return
	}

	_, err = message.Encode(XML)
	if err != nil {
		t.Fatalf("encoding message to XML failed: ", err)
		return
	}

	_, err = resolvedMessage.Encode(XML)
	if err != nil {
		t.Fatalf("encoding resolved message to XML failed: ", err)
		return
	}
}
