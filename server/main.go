package main

import (
	"log"
	"sync"
)

const (
	MAX_TRANSFR = 4096
)

var (
	cnxLock sync.Mutex
	tcpCnxExist = false
	inchannel = make(chan []byte, 1)
)

func main(){

	log.Println("Hello world")

	go func(){if err := startServertoServer(); err != nil {
		log.Println("[FATAL]", err)
	}}()

	/*go func(){ if err := startClienttoServer(linkingCh); err != nil {
		log.Println("[FATAL]", err)
	}}()*/

	select{}
}
