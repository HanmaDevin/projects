package p2p

// Peer is an interface that represents the remote node.
type Peer interface {
}

// Transport is anything that represents the communication
// between the nodes in the network. This network can be
// of the form (TCP, UDP, websockets, ...)
type Transport interface {
}
