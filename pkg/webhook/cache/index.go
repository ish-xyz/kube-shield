package cache

import (
	"fmt"
	"strings"
)

// Add resource address into Cache Index
// Resource address example: NS/GROUP/VERSION/KIND/RULENAME
func (c *CacheIndex) AddPolicyFromAddress(resourceAddress string) error {

	if len(resourceAddress) < 5 {
		return fmt.Errorf("not enough information to process the resource address")
	}
	addr := strings.Split(resourceAddress, "/")

	ns := Namespace(addr[0])
	group := Group(addr[1])
	if group == "" {
		group = Group("null")
	}
	version := Version(addr[2])
	kind := Kind(addr[3])
	rule := RuleName(addr[4])

	if _, ok := c.Policies[ns]; !ok {
		c.Policies[ns] = make(map[Group]map[Version]map[Kind][]RuleName)
	}

	if _, ok := c.Policies[ns][group]; !ok {
		c.Policies[ns][group] = make(map[Version]map[Kind][]RuleName)
	}

	if _, ok := c.Policies[ns][group][version]; !ok {
		c.Policies[ns][group][version] = make(map[Kind][]RuleName)
	}

	if _, ok := c.Policies[ns][group][version][kind]; !ok {
		c.Policies[ns][group][version][kind] = []RuleName{rule}
	} else {
		c.Policies[ns][group][version][kind] = append(c.Policies[ns][group][version][kind], rule)
	}

	return nil
}
