/*
	This implements 'RFC8428' (https://tools.ietf.org/rfc/rfc8428.txt)
*/
package senml

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"time"
)

// Supported version
const SenMLVersion int = 10

// Supported encoding formats
type EncodingFormat int

const (
	JSON EncodingFormat = iota
	XML
)

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

// Parse the message with the given encoding format.
// Returns a non-resolved message, you need to resolve it using Resolve() to get
// base attributes resolution, absolute time, etc.
func Decode(encodedMessage []byte, format EncodingFormat) (message Message, err error) {
	switch {
	case format == JSON:
		err = json.Unmarshal(encodedMessage, &message.Records)
	case format == XML:
		err = xml.Unmarshal(encodedMessage, &message)
	}
	return
}

// Encodes the message with the given encoding format.
func (message Message) Encode(format EncodingFormat) (encodedMessage []byte, err error) {
	switch {
	case format == JSON:
		encodedMessage, err = json.Marshal(message.Records)
	case format == XML:
		encodedMessage, err = xml.Marshal(message)
	}
	return
}

// Resolves the base attributes, calculates absolute time from relative time etc.
func (message Message) Resolve() (resolvedMessage Message, err error) {
	var timeNow float64 = float64(time.Now().Unix())

	var baseName *string = nil
	var baseTime *float64 = nil
	var baseUnit *string = nil
	var baseValue *float64 = nil
	var baseSum *float64 = nil
	var baseVersion *int = nil

	for _, record := range message.Records {
		var resolvedRecord Record = Record{}

		// Base attributes

		if record.BaseVersion != nil {
			if *record.BaseVersion > SenMLVersion {
				err = errors.New(fmt.Sprintf("The version of the record is higher than supported. (expected: %v, got: %v)", SenMLVersion, *record.BaseVersion))
				return
			} else if baseVersion == nil {
				baseVersion = record.BaseVersion
			} else if *record.BaseVersion != *baseVersion {
				err = errors.New("The BaseVersion of the records should all be the same.")
				return
			}
		} else if baseVersion == nil {
			var defaultVersion int = SenMLVersion
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

		// Name

		var resolvedName string = ""
		if baseName != nil {
			resolvedName = *baseName
		}
		if record.Name != nil {
			resolvedName += *record.Name
		}
		if len(resolvedName) == 0 {
			err = errors.New("The concatenated name MUST not be empty to uniquely identify and differentiate the sensor from all others.")
			return
		}
		validNameCharsExp := regexp.MustCompile(`^[a-zA-Z0-9\-\:\.\/\_]*$`)
		if !validNameCharsExp.MatchString(resolvedName) {
			err = errors.New("The concatenated name MUST consist only of characters out of the set \"A\" to \"Z\", \"a\" to \"z\", and \"0\" to \"9\", as well as \"-\", \":\", \".\", \"/\", and \"_\".")
			return
		}
		validFirstCharacterExp := regexp.MustCompile(`^[a-zA-Z0-9]*$`)
		if !validFirstCharacterExp.MatchString(resolvedName[:1]) {
			err = errors.New("The concatenated name MUST start with a character out of the set \"A\" to \"Z\", \"a\" to \"z\", or \"0\" to \"9\".")
			return
		}
		resolvedRecord.Name = &resolvedName

		// Unit

		if record.Unit != nil {
			var resolvedUnit string = *record.Unit
			resolvedRecord.Unit = &resolvedUnit
		} else if baseUnit != nil {
			var resolvedUnit string = *baseUnit
			resolvedRecord.Unit = &resolvedUnit
		}

		// Value

		var resolvedValue float64 = 0
		if baseValue != nil {
			resolvedValue = *baseValue
		}
		if record.Value != nil {
			resolvedValue += *record.Value
		}
		if baseValue != nil || record.Value != nil {
			resolvedRecord.Value = &resolvedValue
		}

		// BoolValue

		if record.BoolValue != nil {
			var resolvedBoolValue bool = *record.BoolValue
			resolvedRecord.BoolValue = &resolvedBoolValue
		}

		// StringValue

		if record.StringValue != nil {
			var resolvedStringValue string = *record.StringValue
			resolvedRecord.StringValue = &resolvedStringValue
		}

		// DataValue

		if record.DataValue != nil {
			var resolvedDataValue string = *record.DataValue
			resolvedRecord.DataValue = &resolvedDataValue
		}

		// Sum

		var resolvedSum float64 = 0
		if baseSum != nil {
			resolvedSum = *baseSum
		}
		if record.Sum != nil {
			resolvedSum += *record.Sum
		}
		if baseSum != nil || record.Sum != nil {
			resolvedRecord.Sum = &resolvedSum
		}

		// Time

		var resolvedTime float64 = 0
		if baseTime != nil {
			resolvedTime = *baseTime
		}
		if record.Time != nil {
			resolvedTime += *record.Time
		}
		if baseTime != nil || record.Time != nil {
			if resolvedTime < 2^28 {
				var absoluteTime float64 = resolvedTime + timeNow
				resolvedRecord.Time = &absoluteTime
			} else {
				resolvedRecord.Time = &resolvedTime
			}
		}

		// UpdateTime

		if record.UpdateTime != nil {
			var resolvedUpdateTime float64 = *record.UpdateTime
			resolvedRecord.UpdateTime = &resolvedUpdateTime
		}

		// Check if a value or sum is set

		if resolvedRecord.Value == nil && resolvedRecord.StringValue == nil && resolvedRecord.BoolValue == nil && resolvedRecord.DataValue == nil && resolvedRecord.Sum == nil {
			err = errors.New("The record has no Value, StringValue, BoolValue, DataValue or Sum.")
			return
		}

		resolvedMessage.Records = append(resolvedMessage.Records, resolvedRecord)
	}

	// Set BaseVersion if necessary

	if baseVersion != nil && *baseVersion != SenMLVersion {
		for _, record := range resolvedMessage.Records {
			var resolvedVersion int = *baseVersion
			record.BaseVersion = &resolvedVersion
		}
	}

	// TODO: sort the records to be in chronological order

	return
}
