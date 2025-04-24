package main

import (
	"log"
)

const (
	MAX_TRANSFR = 4096
)

func main(){
	log.Println("Started tunneling")

	out := make(chan []byte)
	in := make(chan []byte)
	defer close(out)
	defer close(in)

	go StartShell(out, in)
	
	go func(){
		for {
			ff := <- out
			log.Println(string(ff))
		}
	}()

	StartConnection(out, in)

}