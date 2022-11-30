package cache

import "fmt"

func (c *CacheController) onAdd(obj interface{}) {
	fmt.Println(obj)
}

func (c *CacheController) onDelete(obj interface{}) {
	fmt.Println(obj)
}
