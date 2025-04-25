package main

import (
	"log"
	"net"
	"time"
)


func StartConnection(out chan []byte,in chan []byte){
	conn, err := net.Dial("tcp", "localhost:42001")
    if err != nil {
        log.Println("Could not connect to server to build tunnel:", err)
        return
    }
    defer conn.Close()

	go func(){ //keep the cnx alive
		for {
			time.Sleep(time.Second * 5)
			_, err := conn.Write([]byte{0x01})
			if err != nil {
				log.Println("connection to server dropped!")
				conn.Close() //TODO this will break fix it
				return
			}
		}
	}()


	go func(){
		for {
			stream := <- out
			_, err := conn.Write(stream)
			if err != nil {
				log.Println("connection to server dropped!")
				conn.Close() //TODO this will break fix it
				return
			}
		}
	}()
	
	for {

		data := make([]byte, MAX_TRANSFR)
		_, err = conn.Read(data)
		if err != nil {
			log.Println("Could not read from server", err)
			conn.Close() //TODO this will break fix it
			return
		}

		//splt := splitByNilByte(data)

		in <- data
		//ExecCmd(splt, outchannel)
	}
}

func splitByNilByte(b []byte) []byte{
	
	k := make([]byte, 0)
	for i:=0; i<=len(b)-1; i++{
		if b[i] != 0x00 {
			k = append(k, b[i])
		}else{
			break
		}
	}
	return k

}