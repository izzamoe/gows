package mqtt

import "sync"

type counter struct {
	sync.Mutex
	val int
}

func (c *counter) Add(int) {
	c.Lock()
	c.val++
	c.Unlock()
}

func (c *counter) Value() int {
	return c.val
}
