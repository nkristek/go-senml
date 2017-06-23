/*
	This implements 'draft-ietf-core-senml-08' (https://tools.ietf.org/html/draft-ietf-core-senml-08)
*/
package senml

import (
	"encoding/json" // json formatting
	"encoding/xml"  // xml formatting
	"errors"        // throw error when processing failed
	"time"          // get current time
)

// current supported version
const SenMLVersion int = 5

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
// Also deletes records with no value or sum.
func (message SenMLMessage) Resolve() (resolvedMessage SenMLMessage, err error) {
	/*
		Direct quote from the draft:
		Each SenML Pack carries a single array that represents a set of
		measurements and/or parameters.  This array contains a series of
		SenML Records with several attributes described below.  There are two
		kind of attributes: base and regular.  The base attributes can be
		included in the any SenML Record and they apply to the entries in the
		Record.  Each base attribute also applies to all Records after it up
		to, but not including, the next Record that has that same base
		attribute.  All base attributes are optional.  Regular attributes can
		be included in any SenML Record and apply only to that Record.
	*/

	var basename string = ""
	var basetime float64 = 0
	var baseunit string = ""
	var basevalue float64 = 0
	var basesum float64 = 0
	var baseversion int = SenMLVersion

	resolvedMessage.XmlName = message.XmlName
	resolvedMessage.Xmlns = message.Xmlns

	for _, record := range message.Records {
		// get base attributes from current record
		if record.BaseName != nil && len(*record.BaseName) > 0 {
			basename = *record.BaseName
		}
		if record.BaseTime != nil && *record.BaseTime > 0 {
			basetime = *record.BaseTime
		}
		if record.BaseUnit != nil && len(*record.BaseUnit) > 0 {
			baseunit = *record.BaseUnit
		}
		if record.BaseValue != nil && *record.BaseValue > 0 {
			basevalue = *record.BaseValue
		}
		if record.BaseSum != nil && *record.BaseSum > 0 {
			basesum = *record.BaseSum
		}
		if record.Version != nil && *record.Version > baseversion {
			err = errors.New("version number is higher than supported")
			return
		}

		// delete base attributes from record
		record.BaseName = nil
		record.BaseTime = nil
		record.BaseUnit = nil
		record.BaseValue = nil
		record.BaseSum = nil

		// 1. prepend the basename
		combinedName := basename
		if record.Name != nil {
			combinedName += *record.Name
		}
		record.Name = &combinedName

		// 2. add base time to every time field and convert time to absolute
		/*
			Direct quote from the draft:
			A time of zero indicates that the sensor does not know the absolute
			time and the measurement was made roughly "now".  A negative value is
			used to indicate seconds in the past from roughly "now".  A positive
			value is used to indicate the number of seconds, excluding leap
			seconds, since the start of the year 1970 in UTC.
		*/
		combinedTime := basetime
		if record.Time != nil {
			combinedTime += *record.Time
		}
		record.Time = &combinedTime
		if *record.Time <= 0 {
			var now int64 = time.Now().UnixNano() / 1000000000.0
			absoluteTime := float64(now) + *record.Time
			record.Time = &absoluteTime
		}

		// 3. populate base unit on empty unit fields
		if record.Unit == nil || len(*record.Unit) <= 0 {
			currentBaseUnit := baseunit
			record.Unit = &currentBaseUnit
		}

		// 4. add base value to every value field
		combinedValue := basevalue
		if record.Value != nil {
			combinedValue += *record.Value
		}
		record.Value = &combinedValue

		// 5. add base sum to every sum field
		combinedSum := basesum
		if record.Sum != nil {
			combinedSum += *record.Sum
		}
		record.Sum = &combinedSum

		// 6. set version to baseversion
		record.Version = &baseversion

		// add the record to the output message if a value is set
		if record.Value != nil || record.StringValue != nil || record.BoolValue != nil || record.DataValue != nil || record.Sum != nil {
			resolvedMessage.Records = append(resolvedMessage.Records, record)
		}
	}

	return
}
