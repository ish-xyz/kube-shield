package operators

import (
	"fmt"
	"strings"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/tidwall/gjson"
)

const (
	// Operators
	EQUAL         = "Equal"
	NOT_EQUAL     = "NotEqual"
	GREATER       = "Greater"
	LOWER         = "Lower"
	GREATER_EQUAL = "GreaterEqual"
	LOWER_EQUAL   = "LowerEqual"
	// REGEX         = "Regex"
	// EQUAL_ITERATE = "Iterate"

	// check statuses
	CHECK_INIT_ERROR = 1
	CHECK_DONE       = 2
)

func Run(payload string, check *v1.Check) *v1.CheckResult {
	switch check.Operator {
	case EQUAL:
		return equal(payload, check)
	case NOT_EQUAL:
		return notEqual(payload, check)
		// case GREATER:
		// 	return greater(payload, check)
		// case GREATER_EQUAL:
		// 	return greaterOrEqual(payload, check)
		// case LOWER:
		// 	return lower(payload, check)
		// case LOWER_EQUAL:
		// 	return lowerOrEqual(payload, check)
		// case EQUAL_ITERATE:
		// 	return equalIterate(payload, check)
	}
	return &v1.CheckResult{Match: false, Error: fmt.Errorf("unknown operator '%s'", check.Operator)}
}

func getValue(payload, key string) (gjson.Result, error) {
	if strings.HasPrefix(key, `$_.`) {
		return gjson.Get(payload, strings.TrimPrefix(key, "$_.")), nil
	}
	return gjson.Result{}, fmt.Errorf("you are trying to retrieve a non-dynamic value")
}
