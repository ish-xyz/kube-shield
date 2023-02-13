package cache

/*
This file defines the Cache Index logic and methods
*/

import (
	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
)

func NewCacheIndex() *CacheIndex {

	return &CacheIndex{
		Policies: make(map[Operation]map[Namespace]map[Group]map[Resource][]PolicyName),
	}
}

// Delete policy for all entries
func (c *CacheIndex) Delete(entries []*v1.Definition, namespace, name string) {
	for _, def := range entries {
		c.DeleteSingleEntry(
			Operation(def.Operation),
			Namespace(namespace),
			GetGroup(def.Group),
			Resource(def.Resource),
			PolicyName(name),
		)
	}
}

// Add policy for all entries
func (c *CacheIndex) Add(entries []*v1.Definition, namespace, name string) {
	for _, def := range entries {
		c.AddSingleEntry(
			Operation(def.Operation),
			Namespace(namespace),
			GetGroup(def.Group),
			Resource(def.Resource),
			PolicyName(name),
		)
	}
}

func (c *CacheIndex) AddSingleEntry(ops Operation, ns Namespace, grp Group, res Resource, name PolicyName) {

	if _, ok := c.Policies[ops]; !ok {
		c.Policies[ops] = make(map[Namespace]map[Group]map[Resource][]PolicyName)
	}

	if _, ok := c.Policies[ops][ns]; !ok {
		c.Policies[ops][ns] = make(map[Group]map[Resource][]PolicyName)
	}

	if _, ok := c.Policies[ops][ns][grp]; !ok {
		c.Policies[ops][ns][grp] = make(map[Resource][]PolicyName)
	}

	if _, ok := c.Policies[ops][ns][grp][res]; !ok {
		c.Policies[ops][ns][grp][res] = []PolicyName{name}
	} else {
		c.Policies[ops][ns][grp][res] = append(c.Policies[ops][ns][grp][res], name)
	}
}

func (c *CacheIndex) DeleteSingleEntry(ops Operation, ns Namespace, grp Group, res Resource, name PolicyName) {
	if _, exists := c.Policies[ops][ns][grp][res]; exists {
		var newPoliciesArr []PolicyName
		for _, cachedPolicyName := range c.Policies[ops][ns][grp][res] {
			if cachedPolicyName != name {
				newPoliciesArr = append(newPoliciesArr, cachedPolicyName)
			}
		}
		c.Policies[ops][ns][grp][res] = newPoliciesArr
	}
}

func (c *CacheIndex) Get(ops Operation, ns Namespace, grp Group, res Resource) []PolicyName {

	if _, ok := c.Policies[ops][ns][grp][res]; !ok {
		return []PolicyName{}
	} else {
		return c.Policies[ops][ns][grp][res]
	}
}
