// Package senml provides an implementation of RFC 8428
package senml

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"time"
)

// SupportedVersion declares the maximum version of the SenML format supported by this library
const SupportedVersion int = 10

// EncodingFormat declares the supported encoding formats of the SenML message
type EncodingFormat int

const (
	// JSON will use encoding/json to serialize/deserialize the message
	JSON EncodingFormat = iota

	// XML will use encoding/xml to serialize/deserialize the message
	XML
)

// Message is used to serialize and deserialize a SenML message
type Message struct {
	/*
		Used for XML parsing
	*/
	XMLName xml.Name `json:"-" xml:"urn:ietf:params:xml:ns:senml sensml"`

	/*
		Records of the message
	*/
	Records []Record `xml:"senml"`
}

// Record is a single record inside a SenML message
type Record struct {
	/*
		Used for XML parsing
	*/
	XMLName xml.Name `json:"-" xml:"senml"`

	/*
		This is a string that is prepended to the names found in the entries.
	*/
	BaseName *string `json:"bn,omitempty" xml:"bn,attr,omitempty"`

	/*
		A base time that is added to the time found in an entry.
	*/
	BaseTime *float64 `json:"bt,omitempty" xml:"bt,attr,omitempty"`

	/*
		A base unit that is assumed for all entries, unless
		otherwise indicated. If a record does not contain a Unit value,
		then the Base Unit is used. Otherwise, the value found in the
		Unit (if any) is used.
	*/
	BaseUnit *string `json:"bu,omitempty" xml:"bu,attr,omitempty"`

	/*
		A base value is added to the value found in an entry, similar to Base Time.
	*/
	BaseValue *float64 `json:"bv:omitempty" xml:"bv,attr,omitempty"`

	/*
		A base sum is added to the sum found in an entry, similar to Base Time.
	*/
	BaseSum *float64 `json:"bs:omitempty" xml:"bs,attr,omitempty"`

	/*
		Version number of the media type format. This field is an optional positive integer and defaults to 10 if not present.
	*/
	BaseVersion *int `json:"bver,omitempty" xml:"bver,attr,omitempty"`

	/*
		Name of the sensor or parameter. When appended to the Base
		Name field, this must result in a globally unique identifier for
		the resource. The name is optional, if the Base Name is present.
		If the name is missing, the Base Name must uniquely identify the
		resource. This can be used to represent a large array of
		measurements from the same sensor without having to repeat its
		identifier on every measurement.
	*/
	Name *string `json:"n,omitempty" xml:"n,attr,omitempty"`

	/*
		Unit for a measurement value.  Optional.
	*/
	Unit *string `json:"u,omitempty" xml:"u,attr,omitempty"`

	/*
		Value of the entry. Optional if a Sum value is present;
		otherwise, it's required. Values are represented using basic data
		types. This specification defines floating-point numbers ("v"
		field for "Value"), booleans ("vb" for "Boolean Value"), strings
		("vs" for "String Value"), and binary data ("vd" for "Data
		Value"). Exactly one Value field MUST appear unless there is a
		Sum field, in which case it is allowed to have no Value field.
	*/
	Value       *float64 `json:"v,omitempty" xml:"v,attr,omitempty"`
	BoolValue   *bool    `json:"vb,omitempty" xml:"vb,attr,omitempty"`
	StringValue *string  `json:"vs,omitempty" xml:"vs,attr,omitempty"`
	DataValue   *string  `json:"vd,omitempty" xml:"vd,attr,omitempty"`

	/*
		Integrated sum of the values over time. Optional. This field
		is in the unit specified in the Unit value multiplied by seconds.
		For historical reasons, it is named "sum" instead of "integral".
	*/
	Sum *float64 `json:"s,omitempty" xml:"s,attr,omitempty"`

	/*
		Time when the value was recorded. Optional.
	*/
	Time *float64 `json:"t,omitempty" xml:"t,attr,omitempty"`

	/*
		Period of time in seconds that represents the maximum
		time before this sensor will provide an updated reading for a
		measurement.  Optional.  This can be used to detect the failure of
		sensors or the communications path from the sensor.
	*/
	UpdateTime *float64 `json:"ut,omitempty" xml:"ut,attr,omitempty"`
}

// Decode parses the message with the given encoding format.
// Returns a non-resolved message, you need to resolve it using Resolve() to get
// base attributes resolution, absolute time, etc.
func Decode(encodedMessage []byte, format EncodingFormat) (message Message, err error) {
	switch {
	case format == JSON:
		err = json.Unmarshal(encodedMessage, &message.Records)
	case format == XML:
		err = xml.Unmarshal(encodedMessage, &message)
	default:
		err = errors.New("Unsupported encoding format")
		return
	}
	return
}

// Encode encodes the message with the given encoding format.
func (message Message) Encode(format EncodingFormat) ([]byte, error) {
	switch {
	case format == JSON:
		return json.Marshal(message.Records)
	case format == XML:
		return xml.Marshal(message)
	default:
		return nil, errors.New("Unsupported encoding format")
	}
}

// Resolve adds the base attributes to the normal attributes, calculates absolute time from relative time etc.
func (message Message) Resolve() (resolvedMessage Message, err error) {
	var timeNow = float64(time.Now().Unix())

	var baseName *string
	var baseTime *float64
	var baseUnit *string
	var baseValue *float64
	var baseSum *float64
	var baseVersion *int

	for _, record := range message.Records {
		var resolvedRecord = Record{}

		if record.BaseVersion != nil {
			if *record.BaseVersion > SupportedVersion {
				err = fmt.Errorf("The version of the record is higher than supported. (expected: %v, got: %v)", SupportedVersion, *record.BaseVersion)
				return
			} else if baseVersion == nil {
				baseVersion = record.BaseVersion
			} else if *record.BaseVersion != *baseVersion {
				err = errors.New("The BaseVersion of the records should all be the same")
				return
			}
		} else if baseVersion == nil {
			var defaultVersion = SupportedVersion
			baseVersion = &defaultVersion
		}
		if record.BaseName != nil {
			baseName = record.BaseName
		}
		if record.BaseTime != nil {
			baseTime = record.BaseTime
		}
		if record.BaseUnit != nil {
			baseUnit = record.BaseUnit
		}
		if record.BaseValue != nil {
			baseValue = record.BaseValue
		}
		if record.BaseSum != nil {
			baseSum = record.BaseSum
		}

		resolvedRecord.Name, err = resolveName(baseName, record.Name)
		if err != nil {
			return
		}
		resolvedRecord.Unit = resolveUnit(baseUnit, record.Unit)
		resolvedRecord.Value = resolveValue(baseValue, record.Value)
		resolvedRecord.BoolValue = resolveBoolValue(record.BoolValue)
		resolvedRecord.StringValue = resolveStringValue(record.StringValue)
		resolvedRecord.DataValue = resolveDataValue(record.DataValue)
		resolvedRecord.Sum = resolveSum(baseSum, record.Sum)
		resolvedRecord.Time = resolveTime(baseTime, record.Time, timeNow)
		resolvedRecord.UpdateTime = resolveUpdateTime(record.UpdateTime)

		err = validateRecordHasValue(resolvedRecord)
		if err != nil {
			return
		}

		resolvedMessage.Records = append(resolvedMessage.Records, resolvedRecord)
	}

	setBaseVersionIfNecessary(resolvedMessage, baseVersion)
	sortRecordsChronologically(resolvedMessage.Records)
	return
}

func resolveName(baseName *string, name *string) (*string, error) {
	var resolvedName string
	if baseName != nil {
		resolvedName = *baseName
	}
	if name != nil {
		resolvedName += *name
	}
	if len(resolvedName) == 0 {
		return nil, errors.New("The concatenated name MUST not be empty to uniquely identify and differentiate the sensor from all others")
	}
	validNameCharsExp := regexp.MustCompile(`^[a-zA-Z0-9\-\:\.\/\_]*$`)
	if !validNameCharsExp.MatchString(resolvedName) {
		return nil, errors.New("The concatenated name MUST consist only of characters out of the set \"A\" to \"Z\", \"a\" to \"z\", and \"0\" to \"9\", as well as \"-\", \":\", \".\", \"/\", and \"_\"")
	}
	validFirstCharacterExp := regexp.MustCompile(`^[a-zA-Z0-9]*$`)
	if !validFirstCharacterExp.MatchString(resolvedName[:1]) {
		return nil, errors.New("The concatenated name MUST start with a character out of the set \"A\" to \"Z\", \"a\" to \"z\", or \"0\" to \"9\"")
	}
	return &resolvedName, nil
}

func resolveUnit(baseUnit *string, unit *string) *string {
	if unit != nil {
		var resolvedUnit = *unit
		return &resolvedUnit
	} else if baseUnit != nil {
		var resolvedUnit = *baseUnit
		return &resolvedUnit
	}
	return nil
}

func resolveValue(baseValue *float64, value *float64) *float64 {
	var resolvedValue float64
	if baseValue != nil {
		resolvedValue = *baseValue
	}
	if value != nil {
		resolvedValue += *value
	}
	if baseValue != nil || value != nil {
		return &resolvedValue
	}
	return nil
}

func resolveBoolValue(value *bool) *bool {
	if value != nil {
		var resolvedBoolValue = *value
		return &resolvedBoolValue
	}
	return nil
}

func resolveStringValue(value *string) *string {
	if value != nil {
		var resolvedStringValue = *value
		return &resolvedStringValue
	}
	return nil
}

func resolveDataValue(value *string) *string {
	if value != nil {
		var resolvedDataValue = *value
		return &resolvedDataValue
	}
	return nil
}

func resolveSum(baseSum *float64, sum *float64) *float64 {
	var resolvedSum float64
	if baseSum != nil {
		resolvedSum = *baseSum
	}
	if sum != nil {
		resolvedSum += *sum
	}
	if baseSum != nil || sum != nil {
		return &resolvedSum
	}
	return nil
}

func resolveTime(baseTime *float64, time *float64, timeNow float64) *float64 {
	var resolvedTime float64
	if baseTime != nil {
		resolvedTime = *baseTime
	}
	if time != nil {
		resolvedTime += *time
	}
	if baseTime != nil || time != nil {
		if resolvedTime < 2^28 {
			resolvedTime += timeNow
		}
		return &resolvedTime
	}
	return nil
}

func resolveUpdateTime(updateTime *float64) *float64 {
	if updateTime != nil {
		var resolvedUpdateTime = *updateTime
		return &resolvedUpdateTime
	}
	return nil
}

func validateRecordHasValue(record Record) error {
	if record.Value == nil && record.StringValue == nil && record.BoolValue == nil && record.DataValue == nil && record.Sum == nil {
		return errors.New("The record has no Value, StringValue, BoolValue, DataValue or Sum")
	}
	return nil
}

func setBaseVersionIfNecessary(message Message, baseVersion *int) {
	if baseVersion != nil && *baseVersion < SupportedVersion {
		for i := range message.Records {
			var resolvedVersion = *baseVersion
			message.Records[i].BaseVersion = &resolvedVersion
		}
	}
}

func sortRecordsChronologically(records []Record) {
	sort.SliceStable(records, func(i, j int) bool {
		var first, second = records[i], records[j]
		if second.Time == nil {
			return false
		}
		if first.Time == nil {
			return true
		}
		return *first.Time < *second.Time
	})
}
