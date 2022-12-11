package operators

import (
	"fmt"
	"strconv"

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

func getStringValue(v gjson.Result) string {
	if v.Str != "" {
		return v.Str
	}
	return fmt.Sprintf("%v", v)
}
