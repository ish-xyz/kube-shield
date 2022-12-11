package operators

import (
	"strings"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/tidwall/gjson"
)

func getTypedValue(v gjson.Result) interface{} {
	switch v.Type.String() {
	case "String":
		return v.Str
	case "Number":
		// always returns float
		return v.Num
	case "True":
		return true
	case "False":
		return false
	case "Null":
		return nil
	default:
		return v.Raw
	}
}

func getPayloadValues(address string, jsonData string) []gjson.Result {
	address = strings.TrimPrefix(address, "$_.")
	return gjson.Get(jsonData, address).Array()

}

func CreateCheckResult(res bool, msg string) *v1.CheckResult {
	return &v1.CheckResult{
		Result:  res,
		Message: msg,
	}
}
