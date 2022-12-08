package operators

import (
	"strings"

	"github.com/tidwall/gjson"
)

func getValues(address string, jsonData string) []gjson.Result {
	address = strings.TrimPrefix(address, "$.")
	return gjson.Get(jsonData, address).Array()
}
