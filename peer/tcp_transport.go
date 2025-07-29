package peer

import (
	"fmt"
	"net"
	"sync"
)

type TCPPeer struct {
	conn     net.Conn
	outbound bool
}

type TCPTransport struct {
	listenAddr string
	handShaker Handshaker
	decoder    Decoder

	msgChan  chan Message
	listener net.Listener
	onPeer   func(Peer) error
	ru       sync.RWMutex
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddr: listenAddr,
		handShaker: NOHandshake,
		decoder:    DefaultDecoder{},
		msgChan:    make(chan Message),
		onPeer: func(p Peer) error {
			fmt.Println("testing...")
			return nil
		},
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}

	go t.acceptLoop()
	return nil
}

func (t *TCPTransport) ConsumeMessage() <-chan Message {
	return t.msgChan
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("tcp accept error")
		}

		fmt.Printf("incoming conn %+v\n", conn)
		go t.manageConn(conn)
	}
}

func (t *TCPTransport) manageConn(conn net.Conn) {
	var err error
	defer func() {
		conn.Close()
		fmt.Printf("dropping connection %s\n", err)
	}()
	peer := NewTCPPeer(conn, true)

	if err = t.handShaker(peer); err != nil {
		return
	}
	if t.onPeer != nil {
		if err = t.onPeer(peer); err != nil {
			return
		}
	}

	msg := Message{}
	for {
		err = t.decoder.Decode(conn, &msg)
		if err != nil {
			return
		}

		msg.From = conn.LocalAddr()
		t.msgChan <- msg
	}
}
