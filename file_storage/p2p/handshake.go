package p2p

// HandshakeFunc... ?
type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(any) error { return nil }
