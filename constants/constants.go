package constants

type ConfigFileFormat string

const (
	//Properties Properties
	PROPERTIES ConfigFileFormat = "properties"
	//XML XML
	XML ConfigFileFormat = "xml"
	//JSON JSON
	JSON ConfigFileFormat = "json"
	//YML YML
	YML ConfigFileFormat = "yml"
	//YAML YAML
	YAML ConfigFileFormat = "yaml"
	// DEFAULT DEFAULT
	DEFAULT ConfigFileFormat = ""
)

// Time format
type timeFormat string

var TimeFormat timeFormat = "2006-01-02 15:04:05"

func (timeFormat) NORMAL() string {
	return "2006-01-02 15:04:05"
}
func (timeFormat) FLOAT1() string {
	return "2006-01-02 15:04:05.9"
}
func (timeFormat) FLOAT2() string {
	return "2006-01-02 15:04:05.99"
}
func (timeFormat) FLOAT3() string {
	return "2006-01-02 15:04:05.999"
}
func (timeFormat) FLOAT01() string {
	return "2006-01-02 15:04:05.0"
}
func (timeFormat) FLOAT02() string {
	return "2006-01-02 15:04:05.00"
}
func (timeFormat) FLOAT03() string {
	return "2006-01-02 15:04:05.000"
}
func (timeFormat) SYMBOL03() string {
	return "20060102150405.000"
}

// File size
type ByteSize float64

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)
