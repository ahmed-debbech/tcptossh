package main

import (
	"log"
	"time"
	_"fmt"
	"os"
)

const (
	MAX_TRANSFR = 4096
)

var (
	ip = ""
	key = ""
)

func main(){
	log.Println("Started tunneling")

	ip = os.Args[1]
	key = os.Args[2]
	
	i:=1
	for {
		log.Println("re-running software for the", i, "time")

		out := make(chan []byte, MAX_TRANSFR)
		in := make(chan []byte, MAX_TRANSFR)
	
		StartConnection(out, in)
		i++

		close(out)
		close(in)
		time.Sleep(time.Second * 10)
	}
}