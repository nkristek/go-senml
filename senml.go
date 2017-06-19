/*
	author: Niclas Kristek
	github.com/nkristek
*/
package senml

/*
	This implements 'draft-ietf-core-senml-08' (https://tools.ietf.org/html/draft-ietf-core-senml-08)
*/

import (
	"encoding/json" // json formatting
	"encoding/xml"  // xml formatting
	"errors"        // throw error when processing failed
	"time"          // get current time
)

type Encoding int

const (
	JSON Encoding = iota
	XML
)

func ParseBytes(payload []byte, format Encoding) (message SenMLMessage, err error) {
	switch {
	case format == JSON:
		err = json.Unmarshal(payload, &message.Records)
	case format == XML:
		err = xml.Unmarshal(payload, &message)
	}
	return
}

// populates the base attributes in the regular attributes and deletes the base attributes in the process
// also deletes records with no value and sum (following the guidelines of the rfc)
func Resolve(message SenMLMessage) (resolvedMessage SenMLMessage, err error) {
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
	var baseversion int = 5 // current version in the draft

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
