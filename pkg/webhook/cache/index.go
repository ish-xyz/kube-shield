package cache

// Add resource address into Cache Index
// Resource address example: NS/GROUP/VERSION/KIND/RULENAME
// TODO: This code is ugly, I'll clean it
func (c *CacheIndex) CachePolicyResourcesMapping(ns Namespace, grp Group, ver Version, kind Kind, rule RuleName) error {

	if _, ok := c.Policies[ns]; !ok {
		c.Policies[ns] = make(map[Group]map[Version]map[Kind][]RuleName)
	}

	if _, ok := c.Policies[ns][grp]; !ok {
		c.Policies[ns][grp] = make(map[Version]map[Kind][]RuleName)
	}

	if _, ok := c.Policies[ns][grp][ver]; !ok {
		c.Policies[ns][grp][ver] = make(map[Kind][]RuleName)
	}

	if _, ok := c.Policies[ns][grp][ver][kind]; !ok {
		c.Policies[ns][grp][ver][kind] = []RuleName{rule}
	} else {
		c.Policies[ns][grp][ver][kind] = append(c.Policies[ns][grp][ver][kind], rule)
	}

	return nil
}
