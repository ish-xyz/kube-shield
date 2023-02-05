package operators

import (
	"fmt"
	"strings"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/tidwall/gjson"
)

const (
	GREATER       = "GreaterThan"
	LOWER         = "LowerThan"
	EQUAL         = "Equal"
	NOT_EQUAL     = "NotEqual"
	EQUAL_ITERATE = "EqualIterate"
)

func Run(payload string, check *v1.Check) (*v1.CheckResult, error) {
	switch check.Operator {
	case EQUAL:
		equal(payload, check.Field, check.Value)
	}
	return nil, fmt.Errorf("unknown operator '%s'", check.Operator)
}

func getPayloadValues(payload, key string) (gjson.Result, error) {
	if strings.HasPrefix(key, `$_.`) {
		return gjson.Get(payload, strings.TrimPrefix(key, "$_.")), nil
	}
	return gjson.Result{}, fmt.Errorf("not a dynamic json value")
}
