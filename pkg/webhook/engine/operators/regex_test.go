package operators

import (
	"testing"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/stretchr/testify/assert"
)

func TestRegexMatch(t *testing.T) {
	jsonData := `{"value": "true"}`
	check1 := &v1.Check{
		Field:    "$_.value",
		Operator: "Regex",
		Value:    "^true$",
	}

	assert.Equal(t, true, regex(jsonData, check1).Result)
}
