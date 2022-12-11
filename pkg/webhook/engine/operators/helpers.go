package operators

import (
	"fmt"
	"strconv"

	"strings"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/tidwall/gjson"
	"golang.org/x/exp/constraints"
)

type GenericNumber interface {
	constraints.Float | constraints.Integer
}

func getNumber(v string) (float64, error) {

	val, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert to number: %v", err)
	}

	return val, err

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

func CreateCheckResult(res bool, msg string) *v1.CheckResult {
	return &v1.CheckResult{
		Result:  res,
		Message: msg,
	}
}
