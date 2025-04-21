package main

import (
	"log"
	"sync"
)

const (
	MAX_TRANSFR = 1024
)

var (
	cnxLock sync.Mutex
	tcpCnxExist = false
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
