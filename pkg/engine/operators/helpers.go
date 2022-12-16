package operators

import (
	"fmt"
	"strings"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/tidwall/gjson"
)

func getTypedPayloadValue(v gjson.Result) interface{} {
	switch v.Type.String() {
	case "String":
		return v.Str
	case "Number":
		if float64(int(v.Num)) == v.Num {
			return int(v.Num)
		}
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

func getPolicyValue(v interface{}, payload string) interface{} {
	if strings.HasPrefix(fmt.Sprintf("%v", v), "$_.") {
		payloadValues := getPayloadValues(v.(string), payload)
		if len(payload) > 0 {
			return getTypedPayloadValue(payloadValues[0])
		}
		return ""
	}
	return v
}

func CreateCheckResult(res bool, msg string) *v1.CheckResult {
	return &v1.CheckResult{
		Result:  res,
		Message: msg,
	}
}
