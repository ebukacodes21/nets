package main

import (
	"eleniyan/peer"
	"fmt"
	"log"
)

func main() {
	t := peer.NewTCPTransport(":8000")
	go func() error {
		for {
			msg := <-t.ConsumeMessage()
			fmt.Printf("%+v\n", msg)
		}
	}()
	if err := t.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
