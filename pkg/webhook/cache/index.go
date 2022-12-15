package cache

func NewCacheIndex() *CacheIndex {

	return &CacheIndex{
		Policies: make(map[Verb]map[Namespace]map[Group]map[Resource][]PolicyName),
	}
}

// Add resource address into Cache Index
// Resource address example: NS/GROUP/VERSION/KIND/RULENAME
func (c *CacheIndex) Add(verb Verb, ns Namespace, grp Group, res Resource, name PolicyName) {

	if _, ok := c.Policies[verb]; !ok {
		c.Policies[verb] = make(map[Namespace]map[Group]map[Resource][]PolicyName)
	}

	if _, ok := c.Policies[verb][ns]; !ok {
		c.Policies[verb][ns] = make(map[Group]map[Resource][]PolicyName)
	}

	if _, ok := c.Policies[verb][ns][grp]; !ok {
		c.Policies[verb][ns][grp] = make(map[Resource][]PolicyName)
	}

	if _, ok := c.Policies[verb][ns][grp][res]; !ok {
		c.Policies[verb][ns][grp][res] = []PolicyName{name}
	} else {
		c.Policies[verb][ns][grp][res] = append(c.Policies[verb][ns][grp][res], name)
	}
}

func (c *CacheIndex) Delete(verb Verb, ns Namespace, grp Group, res Resource, name PolicyName) {
	if _, exists := c.Policies[verb][ns][grp][res]; exists {
		var newPoliciesArr []PolicyName
		for _, cachedPolicyName := range c.Policies[verb][ns][grp][res] {
			if cachedPolicyName != name {
				newPoliciesArr = append(newPoliciesArr, cachedPolicyName)
			}
		}
		c.Policies[verb][ns][grp][res] = newPoliciesArr
	}
}

func (c *CacheIndex) Get(verb Verb, ns Namespace, grp Group, res Resource, name PolicyName) []PolicyName {

	if _, ok := c.Policies[verb][ns][grp][res]; !ok {
		return []PolicyName{}
	} else {
		return c.Policies[verb][ns][grp][res]
	}
}
