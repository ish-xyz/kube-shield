package operators

import (
	"strings"

	"github.com/tidwall/gjson"
)

func isInt(val float64) bool {
	return val == float64(int(val))
}

func getValues(address string, jsonData string) []gjson.Result {
	address = strings.TrimPrefix(address, "$_.")
	return gjson.Get(jsonData, address).Array()
}

func getTypedValue(v gjson.Result) interface{} {
	switch v.Type.String() {
	case "String":
		return v.Str
	case "Number":
		if isInt(v.Num) {
			return int(v.Num)
		}
		return v.Num
	case "True":
		return true
	case "False":
		return false
	case "Null":
		return ""
	default:
		return v.Raw
	}
}
