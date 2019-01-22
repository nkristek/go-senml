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
		{"bn":"urn:dev:ow:10e2073a01080063","bt":1.320067464e+09,
		 "bu":"%RH","v":20},
		{"u":"lon","v":24.30621},
		{"u":"lat","v":60.07965},
		{"t":60,"v":20.3},
		{"u":"lon","t":60,"v":24.30622},
		{"u":"lat","t":60,"v":60.07965},
		{"t":120,"v":20.7},
		{"u":"lon","t":120,"v":24.30623},
		{"u":"lat","t":120,"v":60.07966},
		{"u":"%EL","t":150,"v":98},
		{"t":180,"v":21.2},
		{"u":"lon","t":180,"v":24.30628},
		{"u":"lat","t":180,"v":60.07967}
	  ]`

	// the same message but with already resolved fields
	var resolvedData string = `[
		{"n":"urn:dev:ow:10e2073a01080063","u":"%RH","t":1.320067464e+09,
		 "v":20},
		{"n":"urn:dev:ow:10e2073a01080063","u":"lon","t":1.320067464e+09,
		 "v":24.30621},
		{"n":"urn:dev:ow:10e2073a01080063","u":"lat","t":1.320067464e+09,
		 "v":60.07965},
		{"n":"urn:dev:ow:10e2073a01080063","u":"%RH","t":1.320067524e+09,
		 "v":20.3},
		{"n":"urn:dev:ow:10e2073a01080063","u":"lon","t":1.320067524e+09,
		 "v":24.30622},
		{"n":"urn:dev:ow:10e2073a01080063","u":"lat","t":1.320067524e+09,
		 "v":60.07965},
		{"n":"urn:dev:ow:10e2073a01080063","u":"%RH","t":1.320067584e+09,
		 "v":20.7},
		{"n":"urn:dev:ow:10e2073a01080063","u":"lon","t":1.320067584e+09,
		 "v":24.30623},
		{"n":"urn:dev:ow:10e2073a01080063","u":"lat","t":1.320067584e+09,
		 "v":60.07966},
		{"n":"urn:dev:ow:10e2073a01080063","u":"%EL","t":1.320067614e+09,
		 "v":98},
		{"n":"urn:dev:ow:10e2073a01080063","u":"%RH","t":1.320067644e+09,
		 "v":21.2},
		{"n":"urn:dev:ow:10e2073a01080063","u":"lon","t":1.320067644e+09,
		 "v":24.30628},
		{"n":"urn:dev:ow:10e2073a01080063","u":"lat","t":1.320067644e+09,
		 "v":60.07967}
	  ]`

	// Decode

	message, err := senml.Decode([]byte(testData), senml.JSON)
	if err != nil {
		t.Error("parsing initial JSON failed: ", err)
		return
	}

	resolvedDataMessage, err := senml.Decode([]byte(resolvedData), senml.JSON)
	if err != nil {
		t.Error("parsing initial JSON failed: ", err)
		return
	}

	// Resolve

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("resolving message failed: ", err)
		return
	}

	// Check Resolve

	if len(resolvedMessage.Records) != len(resolvedDataMessage.Records) {
		t.Error("Unequal amount of records")
		return
	}
	for i := 0; i < len(resolvedMessage.Records); i++ {
		record := resolvedMessage.Records[i]
		recordResolved := resolvedDataMessage.Records[i]

		if record.Name == nil && recordResolved.Name != nil || record.Name != nil && recordResolved.Name == nil {
			t.Error("Name of resolved message is not set")
			return
		} else if record.Name != nil && recordResolved.Name != nil && *record.Name != *recordResolved.Name {
			t.Error("Name of resolved message doesnt match")
			return
		}

		if record.Unit == nil && recordResolved.Unit != nil || record.Unit != nil && recordResolved.Unit == nil {
			t.Error("Unit of resolved message is not set")
			return
		} else if record.Unit != nil && recordResolved.Unit != nil && *record.Unit != *recordResolved.Unit {
			t.Error("Unit of resolved message doesnt match")
			return
		}

		if record.Value == nil && recordResolved.Value != nil || record.Value != nil && recordResolved.Value == nil {
			t.Error("Value of resolved message is not set")
			return
		} else if record.Value != nil && recordResolved.Value != nil && *record.Value != *recordResolved.Value {
			t.Error("Value of resolved message doesnt match")
			return
		}

		if record.BoolValue == nil && recordResolved.BoolValue != nil || record.BoolValue != nil && recordResolved.BoolValue == nil {
			t.Error("BoolValue of resolved message is not set")
			return
		} else if record.BoolValue != nil && recordResolved.BoolValue != nil && *record.BoolValue != *recordResolved.BoolValue {
			t.Error("BoolValue of resolved message doesnt match")
			return
		}

		if record.StringValue == nil && recordResolved.StringValue != nil || record.StringValue != nil && recordResolved.StringValue == nil {
			t.Error("StringValue of resolved message is not set")
			return
		} else if record.StringValue != nil && recordResolved.StringValue != nil && *record.StringValue != *recordResolved.StringValue {
			t.Error("StringValue of resolved message doesnt match")
			return
		}

		if record.DataValue == nil && recordResolved.DataValue != nil || record.DataValue != nil && recordResolved.DataValue == nil {
			t.Error("DataValue of resolved message is not set")
			return
		} else if record.DataValue != nil && recordResolved.DataValue != nil && *record.DataValue != *recordResolved.DataValue {
			t.Error("DataValue of resolved message doesnt match")
			return
		}

		if record.Sum == nil && recordResolved.Sum != nil || record.Sum != nil && recordResolved.Sum == nil {
			t.Error("Sum of resolved message is not set")
			return
		} else if record.Sum != nil && recordResolved.Sum != nil && *record.Sum != *recordResolved.Sum {
			t.Error("Sum of resolved message doesnt match")
			return
		}

		if record.Time == nil && recordResolved.Time != nil || record.Time != nil && recordResolved.Time == nil {
			t.Error("Time of resolved message is not set")
			return
		} else if record.Time != nil && recordResolved.Time != nil && *record.Time != *recordResolved.Time {
			t.Error("Time of resolved message doesnt match")
			return
		}

		if record.UpdateTime == nil && recordResolved.UpdateTime != nil || record.UpdateTime != nil && recordResolved.UpdateTime == nil {
			t.Error("UpdateTime of resolved message is not set")
			return
		} else if record.UpdateTime != nil && recordResolved.UpdateTime != nil && *record.UpdateTime != *recordResolved.UpdateTime {
			t.Error("UpdateTime of resolved message doesnt match")
			return
		}
	}

	// Encode

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
