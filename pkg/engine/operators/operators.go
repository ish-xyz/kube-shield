package operators

import (
	"fmt"
	"strings"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/tidwall/gjson"
)

const (
	// Operators
	GREATER       = "GreaterThan"
	LOWER         = "LowerThan"
	EQUAL         = "Equal"
	NOT_EQUAL     = "NotEqual"
	EQUAL_ITERATE = "EqualIterate"

	// check statuses
	CHECK_EXECUTED   = 1
	CHECK_INIT_ERROR = 2
)

func Run(payload string, check *v1.Check) *v1.CheckResult {
	switch check.Operator {
	case EQUAL:
		equal(payload, check)
	}
	return &v1.CheckResult{Match: false, Error: fmt.Errorf("unknown operator '%s'", check.Operator)}
}

func getValue(payload, key string) (gjson.Result, error) {
	if strings.HasPrefix(key, `$_.`) {
		return gjson.Get(payload, strings.TrimPrefix(key, "$_.")), nil
	}
	return gjson.Result{}, fmt.Errorf("you are trying to retrieve a non-dynamic value")
}
