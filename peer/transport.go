package peer

type Peer interface {
	Close() error
}

type Transport interface {
	ListenAndAccept() error
	ConsumeMessage() <-chan Message
}
