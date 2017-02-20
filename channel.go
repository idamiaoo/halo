package halo

import (
	log "github.com/bohler/lib/dlog"

	"sync"
)

type SessionFilter func(*Session) bool

type Channel struct {
	sync.RWMutex
	name           string              // channel name
	uidMap         map[string]*Session // uid map to session pointer
	members        []string            // all user ids
	channelService *channelService     // channel service which contain current channel
}

func newChannel(n string, cs *channelService) *Channel {
	return &Channel{
		name:           n,
		channelService: cs,
		uidMap:         make(map[string]*Session)}
}

func (c *Channel) Member(uid string) *Session {
	c.RLock()
	defer c.RUnlock()

	return c.uidMap[uid]
}

func (c *Channel) Members() []string {
	c.RLock()
	defer c.RUnlock()

	return c.members
}

// Push message to partial client, which filter return true
func (c *Channel) Multicast(route string, v interface{}, filter SessionFilter) error {
	data, err := serializeOrRaw(v)
	if err != nil {
		return err
	}

	c.RLock()
	defer c.RUnlock()

	for _, s := range c.uidMap {
		if !filter(s) {
			continue
		}
		err = DefaultNetService.Push(s, route, data)
		if err != nil {
			log.Log.Error(err.Error())
		}
	}

	return nil
}

// Push message to all client
func (c *Channel) Broadcast(route string, v interface{}) error {
	data, err := serializeOrRaw(v)
	if err != nil {
		return err
	}

	c.RLock()
	defer c.RUnlock()

	for _, s := range c.uidMap {
		err = DefaultNetService.Push(s, route, data)
		if err != nil {
			log.Log.Error(err.Error())
		}
	}

	return err
}

func (c *Channel) IsContain(uid string) bool {
	c.RLock()
	defer c.RUnlock()

	if _, ok := c.uidMap[uid]; ok {
		return true
	}

	return false
}

func (c *Channel) Add(session *Session) {
	c.Lock()
	defer c.Unlock()

	c.uidMap[session.UID] = session
	c.members = append(c.members, session.UID)
}

func (c *Channel) Leave(uid string) {
	if !c.IsContain(uid) {
		return
	}

	c.Lock()
	defer c.Unlock()

	var temp []string
	for i, u := range c.members {
		if u == uid {
			temp = append(temp, c.members[:i]...)
			c.members = append(temp, c.members[(i+1):]...)
			break
		}
	}
	delete(c.uidMap, uid)
}

func (c *Channel) LeaveAll() {
	c.Lock()
	defer c.Unlock()

	c.uidMap = make(map[string]*Session)
	c.members = make([]string, 0)
}

func (c *Channel) Count() int {
	c.RLock()
	defer c.RUnlock()

	return len(c.uidMap)
}

func (c *Channel) Destroy() {
	c.channelService.Lock()
	defer c.channelService.Unlock()

	delete(c.channelService.channels, c.name)
}
