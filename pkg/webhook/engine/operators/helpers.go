package operators

import (
	"fmt"
	"strconv"

	"strings"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/tidwall/gjson"
)

type Number interface {
	int64 | float64
}

func getNumber[T Number](v string) (T, error) {

	numInt, err := strconv.Atoi(v)
	if err == nil {
		return T(numInt), nil
	}

	numFloat, err := strconv.ParseFloat(v, 64)
	if err == nil {
		return T(numFloat), nil
	}

	return T(0), fmt.Errorf("can't convert string %s into float64 or int64", v)
}

func isInt(v float64) bool {
	return v == float64(int(v))
}

func getValues(address string, jsonData string) []gjson.Result {
	address = strings.TrimPrefix(address, "$_.")
	return gjson.Get(jsonData, address).Array()
}

func getStringValue(v gjson.Result) string {
	if v.Str != "" {
		return v.Str
	}
	return fmt.Sprintf("%v", v)
}

func NewCheckResult() *v1.CheckResult {
	return &v1.CheckResult{
		Result: false,
		Error:  "",
	}
}

func CreateCheckResult(res bool, msg string) *v1.CheckResult {
	return &v1.CheckResult{
		Result: res,
		Error:  msg,
	}
}
