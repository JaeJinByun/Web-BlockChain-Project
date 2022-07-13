package rpcserver

import (
	"blockwiki/custom/rpcserver/rpccustom"
	"crypto/elliptic"
	"encoding/gob"
	"log"
	"net"
	"net/rpc"
)

func Start() {
	rpc.Register(new(rpccustom.Bundle))
	gob.Register(elliptic.P256())
	log.Print("Listening port 7775...")
	l, e := net.Listen("tcp", ":7775")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		defer conn.Close()
		go rpc.ServeConn(conn)
	}
}
