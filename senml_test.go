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

	// the same message but with already resolved fields, see https://tools.ietf.org/html/rfc8428#section-5.1.4
	var resolvedTestData string = `[
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

	testDataMessage, err := senml.Decode([]byte(testData), senml.JSON)
	if err != nil {
		t.Error("parsing initial JSON failed: ", err)
		return
	}

	resolvedTestDataMessage, err := senml.Decode([]byte(resolvedTestData), senml.JSON)
	if err != nil {
		t.Error("parsing initial JSON failed: ", err)
		return
	}

	// Resolve

	testDataMessageResolved, err := testDataMessage.Resolve()
	if err != nil {
		t.Error("resolving message failed: ", err)
		return
	}

	// Check Resolve

	compareMessages(t, testDataMessageResolved, resolvedTestDataMessage)

	// Encode

	_, err = testDataMessage.Encode(senml.JSON)
	if err != nil {
		t.Error("encoding message to JSON failed: ", err)
		return
	}
}

func compareMessages(t *testing.T, firstMessage senml.SenMLMessage, secondMessage senml.SenMLMessage) {
	if len(firstMessage.Records) != len(secondMessage.Records) {
		t.Error("Unequal amount of records")
		return
	}
	for i := 0; i < len(firstMessage.Records); i++ {
		firstMessageRecord := firstMessage.Records[i]
		secondMessageRecord := secondMessage.Records[i]

		if firstMessageRecord.Name == nil && secondMessageRecord.Name != nil || firstMessageRecord.Name != nil && secondMessageRecord.Name == nil {
			t.Error("Name of resolved message is not set")
			return
		} else if firstMessageRecord.Name != nil && secondMessageRecord.Name != nil && *firstMessageRecord.Name != *secondMessageRecord.Name {
			t.Error("Name of resolved message doesnt match")
			return
		}

		if firstMessageRecord.Unit == nil && secondMessageRecord.Unit != nil || firstMessageRecord.Unit != nil && secondMessageRecord.Unit == nil {
			t.Error("Unit of resolved message is not set")
			return
		} else if firstMessageRecord.Unit != nil && secondMessageRecord.Unit != nil && *firstMessageRecord.Unit != *secondMessageRecord.Unit {
			t.Error("Unit of resolved message doesnt match")
			return
		}

		if firstMessageRecord.Value == nil && secondMessageRecord.Value != nil || firstMessageRecord.Value != nil && secondMessageRecord.Value == nil {
			t.Error("Value of resolved message is not set")
			return
		} else if firstMessageRecord.Value != nil && secondMessageRecord.Value != nil && *firstMessageRecord.Value != *secondMessageRecord.Value {
			t.Error("Value of resolved message doesnt match")
			return
		}

		if firstMessageRecord.BoolValue == nil && secondMessageRecord.BoolValue != nil || firstMessageRecord.BoolValue != nil && secondMessageRecord.BoolValue == nil {
			t.Error("BoolValue of resolved message is not set")
			return
		} else if firstMessageRecord.BoolValue != nil && secondMessageRecord.BoolValue != nil && *firstMessageRecord.BoolValue != *secondMessageRecord.BoolValue {
			t.Error("BoolValue of resolved message doesnt match")
			return
		}

		if firstMessageRecord.StringValue == nil && secondMessageRecord.StringValue != nil || firstMessageRecord.StringValue != nil && secondMessageRecord.StringValue == nil {
			t.Error("StringValue of resolved message is not set")
			return
		} else if firstMessageRecord.StringValue != nil && secondMessageRecord.StringValue != nil && *firstMessageRecord.StringValue != *secondMessageRecord.StringValue {
			t.Error("StringValue of resolved message doesnt match")
			return
		}

		if firstMessageRecord.DataValue == nil && secondMessageRecord.DataValue != nil || firstMessageRecord.DataValue != nil && secondMessageRecord.DataValue == nil {
			t.Error("DataValue of resolved message is not set")
			return
		} else if firstMessageRecord.DataValue != nil && secondMessageRecord.DataValue != nil && *firstMessageRecord.DataValue != *secondMessageRecord.DataValue {
			t.Error("DataValue of resolved message doesnt match")
			return
		}

		if firstMessageRecord.Sum == nil && secondMessageRecord.Sum != nil || firstMessageRecord.Sum != nil && secondMessageRecord.Sum == nil {
			t.Error("Sum of resolved message is not set, first: ", firstMessageRecord.Sum, ", second: ", secondMessageRecord.Sum)
			return
		} else if firstMessageRecord.Sum != nil && secondMessageRecord.Sum != nil && *firstMessageRecord.Sum != *secondMessageRecord.Sum {
			t.Error("Sum of resolved message doesnt match")
			return
		}

		if firstMessageRecord.Time == nil && secondMessageRecord.Time != nil || firstMessageRecord.Time != nil && secondMessageRecord.Time == nil {
			t.Error("Time of resolved message is not set")
			return
		} else if firstMessageRecord.Time != nil && secondMessageRecord.Time != nil && *firstMessageRecord.Time != *secondMessageRecord.Time {
			t.Error("Time of resolved message doesnt match")
			return
		}

		if firstMessageRecord.UpdateTime == nil && secondMessageRecord.UpdateTime != nil || firstMessageRecord.UpdateTime != nil && secondMessageRecord.UpdateTime == nil {
			t.Error("UpdateTime of resolved message is not set")
			return
		} else if firstMessageRecord.UpdateTime != nil && secondMessageRecord.UpdateTime != nil && *firstMessageRecord.UpdateTime != *secondMessageRecord.UpdateTime {
			t.Error("UpdateTime of resolved message doesnt match")
			return
		}
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
