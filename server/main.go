package main

import (
	"log"
	"sync"
	"os"
)

const (
	MAX_TRANSFR = 4096
)

var (
	key = ""
	cnxLock sync.Mutex
	tcpCnxExist = false
	inchannel = make(chan []byte, 1)
)

func main(){

	log.Println("Hello world")
	key = os.Args[1]

	go func(){if err := startServertoServer(); err != nil {
		log.Println("[FATAL]", err)
	}}()

	/*go func(){ if err := startClienttoServer(linkingCh); err != nil {
		log.Println("[FATAL]", err)
	}}()*/

	select{}
}
