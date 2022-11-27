package loader

import (
	"fmt"
	"sync"
)

// Loader is the package that loads the policies and creates the POlicies object
type ApiVersion string

type Kind string

type Policy struct {
	Action   string      // CREATE, DELETE, UPDATE, LIST, READ, CONNECT
	Operator string      // Equal, NotEqual, Regex, GreaterThan, LowerThan
	Expected interface{} // Expected result (not tied to a type)
}

type Store struct {
	sync.Mutex
	Policies map[ApiVersion]map[Kind][]*Policy
}

func NewStore() *Store {

	return &Store{
		Policies: make(map[ApiVersion]map[Kind][]*Policy, 1),
	}
}

// Meant to be executed as goroutine, it watches the filesystem and refreshes the policies
func (st *Store) Watch() {
	for {
		st.Lock()
		// reload
		fmt.Println("reload policies here")
		st.Unlock()
	}
}

func (st *Store) GetMatchingPolicies() []*Policy {
	st.Lock()
	fmt.Println("find matching policies here")
	st.Unlock()

	return nil
}
