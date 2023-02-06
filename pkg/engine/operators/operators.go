package operators

import (
	"fmt"
	"strings"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/tidwall/gjson"
)

const (
	// Operators
	GREATER_THAN  = "GreaterThan"
	LOWER_THAN    = "LowerThan"
	EQUAL         = "Equal"
	NOT_EQUAL     = "NotEqual"
	EQUAL_ITERATE = "EqualIterate"

	// check statuses
	CHECK_INIT_ERROR = 1
	CHECK_OK         = 2

	// behaviours
	RUN_EQUAL        = true
	RUN_NOT_EQUAL    = false
	RUN_GREATER_THAN = true
	RUN_LOWER_THAN   = false
)

func Run(payload string, check *v1.Check) *v1.CheckResult {
	switch check.Operator {
	case EQUAL:
		return equal(payload, check)
	case NOT_EQUAL:
		return notEqual(payload, check)
	// case EQUAL_ITERATE:
	// 	return equalIterate(payload, check)
	case GREATER_THAN:
		return greaterThan(payload, check)
	case LOWER_THAN:
		return lowerThan(payload, check)
	}
	return &v1.CheckResult{Match: false, Error: fmt.Errorf("unknown operator '%s'", check.Operator)}
}

func getValue(payload, key string) (gjson.Result, error) {
	if strings.HasPrefix(key, `$_.`) {
		return gjson.Get(payload, strings.TrimPrefix(key, "$_.")), nil
	}
	return gjson.Result{}, fmt.Errorf("you are trying to retrieve a non-dynamic value")
}
