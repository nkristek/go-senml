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

// current supported version
const SenMLVersion int = 10

// Supported encoding formats
type EncodingFormat int

const (
	JSON EncodingFormat = iota
	XML
)

// Parse the []byte with the given encoding format.
// It returns a non-resolved message, when it succeeds,
// you need to resolve it using the Resolve function to get
// base attributes resolution, absolute time, etc.
func Decode(encodedMessage []byte, format EncodingFormat) (message SenMLMessage, err error) {
	message.XmlName = nil
	message.Xmlns = "urn:ietf:params:xml:ns:senml"

	switch {
	case format == JSON:
		err = json.Unmarshal(encodedMessage, &message.Records)
	case format == XML:
		err = xml.Unmarshal(encodedMessage, &message)
	}

	return
}

// Encodes the message with the given encoding format.
// Please try to use base attributes as often as possible
// to make sure that the encoded data is as small as possible.
// (Basically a non-resolved message.)
func (message SenMLMessage) Encode(format EncodingFormat) (encodedMessage []byte, err error) {
	message.Xmlns = "urn:ietf:params:xml:ns:senml"

	switch {
	case format == JSON:
		encodedMessage, err = json.Marshal(message.Records)
	case format == XML:
		encodedMessage, err = xml.Marshal(message)
	}

	return
}

// Resolves the base attributes and deletes the base attributes afterwards
// and calculates absolute time from relative time.
func (message SenMLMessage) Resolve() (resolvedMessage SenMLMessage, err error) {
	var timeNow float64 = float64(time.Now().Unix())

	var basename *string = nil
	var basetime *float64 = nil
	var baseunit *string = nil
	var basevalue *float64 = nil
	var basesum *float64 = nil
	var baseversion *int = nil

	resolvedMessage.XmlName = message.XmlName
	resolvedMessage.Xmlns = message.Xmlns

	for _, record := range message.Records {

		// Base attributes

		if record.BaseVersion != nil {
			if *record.BaseVersion > SenMLVersion {
				err = errors.New(fmt.Sprintf("The version of the record is higher than supported. (expected: %v, got: %v)", SenMLVersion, *record.BaseVersion))
				return
			} else if baseversion == nil {
				baseversion = record.BaseVersion
			} else if *record.BaseVersion != *baseversion {
				err = errors.New("The version of the records should all be the same.")
				return
			}
		}

		if record.BaseName != nil {
			basename = record.BaseName
		}
		if record.BaseTime != nil {
			basetime = record.BaseTime
		}
		if record.BaseUnit != nil {
			baseunit = record.BaseUnit
		}
		if record.BaseValue != nil {
			basevalue = record.BaseValue
		}
		if record.BaseSum != nil {
			basesum = record.BaseSum
		}

		// Name

		var resolvedName string = ""
		if basename != nil {
			resolvedName = *basename
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

		record.Name = &resolvedName

		// Time

		var resolvedTime float64 = 0
		if basetime != nil {
			resolvedTime = *basetime
		}
		if record.Time != nil {
			resolvedTime += *record.Time
		}

		if resolvedTime < 2^28 {
			var absoluteTime float64 = resolvedTime + timeNow
			record.Time = &absoluteTime
		} else {
			record.Time = &resolvedTime
		}

		// Unit

		if record.Unit == nil && baseunit != nil {
			var unit = *baseunit
			record.Unit = &unit
		}

		// Value

		var resolvedValue float64 = 0
		if basevalue != nil {
			resolvedValue = *basevalue
		}
		if record.Value != nil {
			resolvedValue += *record.Value
		}
		if basevalue != nil || record.Value != nil {
			record.Value = &resolvedValue
		}

		// Sum

		var resolvedSum float64 = 0
		if basesum != nil {
			resolvedSum = *basesum
		}
		if record.Sum != nil {
			resolvedSum += *record.Sum
		}
		if basesum != nil || record.Sum != nil {
			record.Sum = &resolvedSum
		}

		// Check if a value or sum is set

		if record.Value == nil && record.StringValue == nil && record.BoolValue == nil && record.DataValue == nil && record.Sum == nil {
			err = errors.New("The record has no Value, StringValue, BoolValue, DataValue or Sum.")
			return
		}

		resolvedMessage.Records = append(resolvedMessage.Records, record)
	}

	// Clear base attributes

	for _, record := range message.Records {
		record.BaseName = nil
		record.BaseTime = nil
		record.BaseUnit = nil
		record.BaseValue = nil
		record.BaseSum = nil

		if baseversion != nil && *baseversion != SenMLVersion {
			version := *baseversion
			record.BaseVersion = &version
		} else {
			record.BaseVersion = nil
		}
	}

	// TODO: sort the records to be in chronological order

	return
}
