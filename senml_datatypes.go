package senml

type SenMLMessage struct {
	/*
		Used for XML parsing
	*/
	XmlName *bool  `json:"_,omitempty" xml:"sensml"`
	Xmlns   string `json:"_,omitempty" xml:"xmlns,attr"`

	/*
		Records of the message
	*/
	Records []SenMLRecord `xml:"senml"`
}

type SenMLRecord struct {
	/*
		Used for XML parsing
	*/
	XmlName *bool `json:"_,omitempty" xml:"senml"`

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
