package cache

import (
	"fmt"
)

/*
This a general purpose file, all helper functions for this package should end up here
*/

// Returning the API group name typed as Group
func GetGroup(grp string) Group {
	if grp == "" {
		return CORE_GROUP
	}
	return Group(grp)
}

func GetResource(res, subres string) Resource {
	if subres != "" {
		return Resource(fmt.Sprintf("%s/%s", res, subres))
	}
	return Resource(res)
}
