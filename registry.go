package gostun

//Registry maintains the mappings required for P2P transmission
type Registry struct {
	mappings   map[string]*Client
	numClients int
}

//RegisterClient registers a client to its uid
func (r *Registry) RegisterClient(c *Client) {
	r.mappings[c.UID] = c
	r.numClients++
}

//RemoveClient removes a client fromt he registry
func (r *Registry) RemoveClient(uid string) {
	delete(r.mappings, uid)
	r.numClients--
}
