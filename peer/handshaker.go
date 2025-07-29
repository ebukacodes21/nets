package peer

type Handshaker func(any) error

func NOHandshake(any) error { return nil }
