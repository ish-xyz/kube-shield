package cache

// Add resource address into Cache Index
// Resource address example: NS/GROUP/VERSION/KIND/RULENAME
func (c *CacheIndex) Add(ns Namespace, grp Group, ver Version, kind Kind, name PolicyName) {

	if _, ok := c.Policies[ns]; !ok {
		c.Policies[ns] = make(map[Group]map[Version]map[Kind][]PolicyName)
	}

	if _, ok := c.Policies[ns][grp]; !ok {
		c.Policies[ns][grp] = make(map[Version]map[Kind][]PolicyName)
	}

	if _, ok := c.Policies[ns][grp][ver]; !ok {
		c.Policies[ns][grp][ver] = make(map[Kind][]PolicyName)
	}

	if _, ok := c.Policies[ns][grp][ver][kind]; !ok {
		c.Policies[ns][grp][ver][kind] = []PolicyName{name}
	} else {
		c.Policies[ns][grp][ver][kind] = append(c.Policies[ns][grp][ver][kind], name)
	}
}

func (c *CacheIndex) Delete(ns Namespace, grp Group, ver Version, kind Kind, name PolicyName) {
	if _, exists := c.Policies[ns][grp][ver][kind]; exists {
		var newPoliciesArr []PolicyName
		for _, cachedPolicyName := range c.Policies[ns][grp][ver][kind] {
			if cachedPolicyName != name {
				newPoliciesArr = append(newPoliciesArr, cachedPolicyName)
			}
		}
		c.Policies[ns][grp][ver][kind] = newPoliciesArr
	}
}
