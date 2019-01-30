package senml_test

import (
	"testing"

	senml "github.com/nkristek/go-senml"
)

// https://tools.ietf.org/html/rfc8428#section-5.1.3
const jsonDataUnresolved string = `[
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

// https://tools.ietf.org/html/rfc8428#section-5.1.4
const jsonDataResolved string = `[
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

// https://tools.ietf.org/html/rfc8428#section-7
const xmlDataUnresolved string = `<sensml xmlns="urn:ietf:params:xml:ns:senml">
	<senml bn="urn:dev:ow:10e2073a0108006:" bt="1.276020076001e+09"
	bu="A" bver="5" n="voltage" u="V" v="120.1"></senml>
	<senml n="current" t="-5" v="1.2"></senml>
	<senml n="current" t="-4" v="1.3"></senml>
	<senml n="current" t="-3" v="1.4"></senml>
	<senml n="current" t="-2" v="1.5"></senml>
	<senml n="current" t="-1" v="1.6"></senml>
	<senml n="current" v="1.7"></senml>
  </sensml>`

func TestDecodeJSON(t *testing.T) {
	_, err := senml.Decode([]byte(jsonDataUnresolved), senml.JSON)
	if err != nil {
		t.Error("Decoding JSON failed: ", err)
		return
	}
}

func TestDecodeXML(t *testing.T) {
	_, err := senml.Decode([]byte(xmlDataUnresolved), senml.XML)
	if err != nil {
		t.Error("Decoding XML failed: ", err)
		return
	}
}

func TestDecodeInvalidFormat(t *testing.T) {
	_, err := senml.Decode(nil, -1)
	if err == nil {
		t.Error("Decoding an invalid format should result in an error")
	}
}

func TestEncodeJSON(t *testing.T) {
	message, err := senml.Decode([]byte(jsonDataUnresolved), senml.JSON)
	if err != nil {
		t.Error("Decoding JSON failed: ", err)
		return
	}

	encodedMessage, err := message.Encode(senml.JSON)
	if err != nil {
		t.Error("Encoding message to JSON failed: ", err)
		return
	}
	if encodedMessage == nil || len(encodedMessage) == 0 {
		t.Error("Encoding to JSON resulted in an empty message")
	}
}

func TestEncodeXML(t *testing.T) {
	message, err := senml.Decode([]byte(xmlDataUnresolved), senml.XML)
	if err != nil {
		t.Error("Decoding XML failed: ", err)
		return
	}

	encodedMessage, err := message.Encode(senml.XML)
	if err != nil {
		t.Error("Encoding message to XML failed: ", err)
		return
	}
	if encodedMessage == nil || len(encodedMessage) == 0 {
		t.Error("Encoding to XML resulted in an empty message")
	}
}

func TestEncodeInvalidFormat(t *testing.T) {
	message, err := senml.Decode([]byte(xmlDataUnresolved), senml.XML)
	if err != nil {
		t.Error("Decoding XML failed: ", err)
		return
	}

	_, err = message.Encode(-1)
	if err == nil {
		t.Error("Encoding message to an invalid format should result in an error")
		return
	}
}

func TestResolveExampleData(t *testing.T) {
	message, err := senml.Decode([]byte(jsonDataUnresolved), senml.JSON)
	if err != nil {
		t.Error("Decoding JSON failed: ", err)
		return
	}

	resolvedDataMessage, err := senml.Decode([]byte(jsonDataResolved), senml.JSON)
	if err != nil {
		t.Error("Decoding JSON failed: ", err)
		return
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving message failed: ", err)
		return
	}

	compareMessages(t, resolvedMessage, resolvedDataMessage)
}

func compareMessages(t *testing.T, firstMessage senml.Message, secondMessage senml.Message) {
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

func TestResolveUnsupportedSenMLVersion(t *testing.T) {
	var unsupportedVersion = 11
	var name = "test"
	var value float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				BaseVersion: &unsupportedVersion,
				Name:        &name,
				Value:       &value,
			},
		},
	}

	_, err := message.Resolve()
	if err == nil {
		t.Error("Resolving an unsupported SenML version should result in an error")
		return
	}
}

func TestResolveBaseVersionIsSetIfLowerThanMaximumSupported(t *testing.T) {
	var lowerVersion = 5
	var name = "test"
	var value float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				BaseVersion: &lowerVersion,
				BaseName:    &name,
				Value:       &value,
			},
			{
				Value: &value,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving an lower SenML version than the maximum supported version failed:", err)
		return
	}

	for _, record := range resolvedMessage.Records {
		if record.BaseVersion == nil {
			t.Error("The BaseVersion attribute is not set if the version is lower than the maximum supported version")
			return
		}
		if *record.BaseVersion != lowerVersion {
			t.Error("The BaseVersion attribute is not set to the BaseVersion in the unresolved message")
			return
		}
	}
}

func TestResolveRecordsHaveDifferentVersion(t *testing.T) {
	var version = 5
	var differentVersion = 6
	var name = "test"
	var value float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				BaseVersion: &version,
				BaseName:    &name,
				Value:       &value,
			},
			{
				BaseVersion: &differentVersion,
				Value:       &value,
			},
		},
	}

	_, err := message.Resolve()
	if err == nil {
		t.Error("Resolving a SenML message which contains records with different version should result in an error")
		return
	}
}

func TestResolveNameContainsInvalidSymbols(t *testing.T) {
	var name = "test("
	var value float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:  &name,
				Value: &value,
			},
		},
	}

	_, err := message.Resolve()
	if err == nil {
		t.Error("Resolving a record with a name which contains invalid symbols should result in an error")
		return
	}
}

func TestResolveNameStartsWithInvalidSymbols(t *testing.T) {
	var name = "-test"
	var value float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:  &name,
				Value: &value,
			},
		},
	}

	_, err := message.Resolve()
	if err == nil {
		t.Error("Resolving a record with a name which starts with an invalid symbol should result in an error")
		return
	}
}

func TestResolveNoName(t *testing.T) {
	var value float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				Value: &value,
			},
		},
	}

	_, err := message.Resolve()
	if err == nil {
		t.Error("Resolving a record with no name should result in an error")
		return
	}
}

func TestResolveNoValue(t *testing.T) {
	var name = "test"
	message := senml.Message{
		Records: []senml.Record{
			{
				Name: &name,
			},
		},
	}

	_, err := message.Resolve()
	if err == nil {
		t.Error("Resolving a record with no value or sum should result in an error")
		return
	}
}

func TestResolveValue(t *testing.T) {
	var name = "test"
	var value float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:  &name,
				Value: &value,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving a record with a value should not result in an error", err)
		return
	}

	if resolvedMessage.Records[0].Value == nil {
		t.Error("The record in the resolved message has no value")
		return
	}

	if *resolvedMessage.Records[0].Value != value {
		t.Error("The value field has a different value than expected")
		return
	}
}

func TestResolveBoolValue(t *testing.T) {
	var name = "test"
	var boolValue bool = true
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:      &name,
				BoolValue: &boolValue,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving a record with a bool value should not result in an error", err)
		return
	}

	if resolvedMessage.Records[0].BoolValue == nil {
		t.Error("The record in the resolved message has no bool value")
		return
	}

	if *resolvedMessage.Records[0].BoolValue != boolValue {
		t.Error("The bool value field has a different value than expected")
		return
	}
}

func TestResolveStringValue(t *testing.T) {
	var name = "test"
	var stringValue string = "value"
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:        &name,
				StringValue: &stringValue,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving a record with a string value should not result in an error", err)
		return
	}

	if resolvedMessage.Records[0].StringValue == nil {
		t.Error("The record in the resolved message has no string value")
		return
	}

	if *resolvedMessage.Records[0].StringValue != stringValue {
		t.Error("The string value field has a different value than expected")
		return
	}
}

func TestResolveDataValue(t *testing.T) {
	var name = "test"
	var dataValue string = "data"
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:      &name,
				DataValue: &dataValue,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving a record with a data value should not result in an error", err)
		return
	}

	if resolvedMessage.Records[0].DataValue == nil {
		t.Error("The record in the resolved message has no data value")
		return
	}

	if *resolvedMessage.Records[0].DataValue != dataValue {
		t.Error("The data value field has a different value than expected")
		return
	}
}

func TestResolveSum(t *testing.T) {
	var name = "test"
	var sum float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				Name: &name,
				Sum:  &sum,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving a record with a sum should not result in an error", err)
		return
	}

	if resolvedMessage.Records[0].Sum == nil {
		t.Error("The record in the resolved message has no sum")
		return
	}

	if *resolvedMessage.Records[0].Sum != sum {
		t.Error("The sum field has a different value than expected")
		return
	}
}

func TestResolveUnit(t *testing.T) {
	var name = "test"
	var value float64 = 1
	var unit = "unit"
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:  &name,
				Value: &value,
				Unit:  &unit,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].Unit == nil {
		t.Error("The record in the resolved message has no unit")
		return
	}

	if *resolvedMessage.Records[0].Unit != unit {
		t.Error("The unit field has a different value than expected")
		return
	}
}

func TestResolveUpdateTime(t *testing.T) {
	var name = "test"
	var value float64 = 1
	var updateTime float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:       &name,
				Value:      &value,
				UpdateTime: &updateTime,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].UpdateTime == nil {
		t.Error("The record in the resolved message has no update time")
		return
	}

	if *resolvedMessage.Records[0].UpdateTime != updateTime {
		t.Error("The update time field has a different value than expected")
		return
	}
}

func TestResolveRelativeToAbsoluteTime(t *testing.T) {
	var name = "test"
	var value float64 = 1
	var time float64 = 2
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:  &name,
				Value: &value,
				Time:  &time,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].Time == nil {
		t.Error("The record in the resolved message has no time")
		return
	}

	if *resolvedMessage.Records[0].Time == time {
		t.Error("The time field was not resolved")
	}
}

func TestResolveAbsoluteTime(t *testing.T) {
	var name = "test"
	var value float64 = 1
	var time float64 = 2 ^ 28
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:  &name,
				Value: &value,
				Time:  &time,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].Time == nil {
		t.Error("The record in the resolved message has no time")
		return
	}

	if *resolvedMessage.Records[0].Time != time {
		t.Error("The value of the time field was changed, but the RFC specifies that values of 2^28 or over should not be changed")
	}
}

func TestResolveOrderIsChronological(t *testing.T) {
	var baseName = "test"
	var value float64 = 1
	var value2 float64 = 2
	var value3 float64 = 3
	var value4 float64 = 4
	var time3 float64 = 3
	var time4 float64 = 4
	message := senml.Message{
		Records: []senml.Record{
			{
				BaseName: &baseName,
				Value:    &value4,
				Time:     &time4,
			},
			{
				Value: &value,
			},
			{
				Value: &value2,
			},
			{
				Value: &value3,
				Time:  &time3,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].Value == nil {
		t.Error("The record in the resolved message has no value")
		return
	}

	if *resolvedMessage.Records[0].Value != value {
		t.Error("The records are not in chronological order")
		return
	}

	if resolvedMessage.Records[1].Value == nil {
		t.Error("The record in the resolved message has no value")
		return
	}

	if *resolvedMessage.Records[1].Value != value2 {
		t.Error("The records are not in chronological order")
		return
	}

	if resolvedMessage.Records[2].Value == nil {
		t.Error("The record in the resolved message has no value")
		return
	}

	if *resolvedMessage.Records[2].Value != value3 {
		t.Error("The records are not in chronological order")
		return
	}

	if resolvedMessage.Records[3].Value == nil {
		t.Error("The record in the resolved message has no value")
		return
	}

	if *resolvedMessage.Records[3].Value != value4 {
		t.Error("The records are not in chronological order")
		return
	}
}

func TestResolveBaseName(t *testing.T) {
	var baseName = "base/"
	var name = "test"
	var value float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				BaseName: &baseName,
				Value:    &value,
			},
			{
				Name:  &name,
				Value: &value,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].BaseName != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[0].Name == nil {
		t.Error("The record in the resolved message has no name")
		return
	}

	if *resolvedMessage.Records[0].Name != baseName {
		t.Error("The base attribute was not properly concatenated with the field")
		return
	}

	if resolvedMessage.Records[1].BaseName != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[1].Name == nil {
		t.Error("The record in the resolved message has no name")
		return
	}

	if *resolvedMessage.Records[1].Name != baseName+name {
		t.Error("The base attribute was not properly concatenated with the field")
		return
	}
}

func TestResolveBaseTime(t *testing.T) {
	var baseTime float64 = 2 ^ 28
	var baseName = "test"
	var baseValue float64 = 1
	var time float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				BaseTime:  &baseTime,
				BaseName:  &baseName,
				BaseValue: &baseValue,
			},
			{
				Time: &time,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].BaseTime != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[0].Time == nil {
		t.Error("The record in the resolved message has no time")
		return
	}

	if *resolvedMessage.Records[0].Time != baseTime {
		t.Error("The base attribute was not properly added to the field")
		return
	}

	if resolvedMessage.Records[1].BaseTime != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[1].Time == nil {
		t.Error("The record in the resolved message has no time")
		return
	}

	if *resolvedMessage.Records[1].Time != baseTime+time {
		t.Error("The base attribute was not properly added to the field")
		return
	}
}

func TestResolveBaseUnit(t *testing.T) {
	var baseUnit = "bu"
	var baseName = "test"
	var baseValue float64 = 1
	var unit string = "u"
	message := senml.Message{
		Records: []senml.Record{
			{
				BaseUnit:  &baseUnit,
				BaseName:  &baseName,
				BaseValue: &baseValue,
			},
			{
				Unit: &unit,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].BaseUnit != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[0].Unit == nil {
		t.Error("The record in the resolved message has no unit")
		return
	}

	if *resolvedMessage.Records[0].Unit != baseUnit {
		t.Error("The base attribute was not properly set")
		return
	}

	if resolvedMessage.Records[1].BaseUnit != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[1].Unit == nil {
		t.Error("The record in the resolved message has no unit")
		return
	}

	if *resolvedMessage.Records[1].Unit != unit {
		t.Error("The field was replaced with the base attribute")
		return
	}
}

func TestResolveBaseValue(t *testing.T) {
	var baseValue float64 = 1
	var baseName = "test"
	var value float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				BaseValue: &baseValue,
				BaseName:  &baseName,
			},
			{
				Value: &value,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].BaseValue != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[0].Value == nil {
		t.Error("The record in the resolved message has no value")
		return
	}

	if *resolvedMessage.Records[0].Value != baseValue {
		t.Error("The base attribute was not properly added to the field")
		return
	}

	if resolvedMessage.Records[1].BaseValue != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[1].Value == nil {
		t.Error("The record in the resolved message has no value")
		return
	}

	if *resolvedMessage.Records[1].Value != baseValue+value {
		t.Error("The base attribute was not properly added to the field")
		return
	}
}

func TestResolveBaseSum(t *testing.T) {
	var baseSum float64 = 1
	var baseName = "test"
	var sum float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				BaseSum:  &baseSum,
				BaseName: &baseName,
			},
			{
				Sum: &sum,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].BaseSum != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[0].Sum == nil {
		t.Error("The record in the resolved message has no sum")
		return
	}

	if *resolvedMessage.Records[0].Sum != baseSum {
		t.Error("The base attribute was not properly added to the field")
		return
	}

	if resolvedMessage.Records[1].BaseSum != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[1].Sum == nil {
		t.Error("The record in the resolved message has no sum")
		return
	}

	if *resolvedMessage.Records[1].Sum != baseSum+sum {
		t.Error("The base attribute was not properly added to the field")
		return
	}
}
