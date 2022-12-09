package operators

import (
	"strings"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
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
		return float64(v.Num)
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

func Dispatch(rawPayload string, check *v1.Check) bool {
	return false
}

func NewCheckResult() *v1.CheckResult {
	return &v1.CheckResult{
		Result: false,
		Errors: make([]error, 0),
	}
}

func UpdateCheckResult(checkRes *v1.CheckResult, res bool, err error) *v1.CheckResult {
	checkRes.Result = res
	checkRes.Errors = append(checkRes.Errors, err)
	return checkRes
}
