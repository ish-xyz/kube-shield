package cache

/*
This file defines the Cache Index logic and methods
*/

import (
	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
)

func NewCacheIndex() *CacheIndex {

	return &CacheIndex{
		Policies: make(map[Verb]map[Namespace]map[Group]map[Resource][]PolicyName),
	}
}

// Delete policy for all entries
func (c *CacheIndex) Delete(entries []*v1.Definition, namespace, name string) {
	for _, def := range entries {
		c.DeleteSingleEntry(
			Verb(def.Verb),
			Namespace(namespace),
			GetGroup(def.ApiGroup),
			Resource(def.Resource),
			PolicyName(name),
		)
	}
}

// Add policy for all entries
func (c *CacheIndex) Add(entries []*v1.Definition, namespace, name string) {
	for _, def := range entries {
		c.AddSingleEntry(
			Verb(def.Verb),
			Namespace(namespace),
			GetGroup(def.ApiGroup),
			Resource(def.Resource),
			PolicyName(name),
		)
	}
}

func (c *CacheIndex) AddSingleEntry(verb Verb, ns Namespace, grp Group, res Resource, name PolicyName) {

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

func (c *CacheIndex) DeleteSingleEntry(verb Verb, ns Namespace, grp Group, res Resource, name PolicyName) {
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

func (c *CacheIndex) Get(verb Verb, ns Namespace, grp Group, res Resource) []PolicyName {

	if _, ok := c.Policies[verb][ns][grp][res]; !ok {
		return []PolicyName{}
	} else {
		return c.Policies[verb][ns][grp][res]
	}
}
