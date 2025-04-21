package main

import (
	"log"
	"net"
)


func StartConnection(){
	conn, err := net.Dial("tcp", "localhost:42001")
    if err != nil {
        log.Println("Could not connect to server to build tunnel:", err)
        return
    }
    defer conn.Close()

	for {
		data := make([]byte, MAX_TRANSFR)
		_, err = conn.Read(data)
		if err != nil {
			log.Println("Could not write to server", err)
			return
		}
	}
}