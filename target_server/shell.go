package main 

import (
	"log"
	"os/exec"
	"github.com/creack/pty"
)

func StartShell(out chan []byte, in chan []byte){

	cmd := exec.Command("bash")

	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Println(err)
		return
	}
	defer ptmx.Close()

	size := &pty.Winsize{
		Rows : 500,
		Cols : 500,
		X    : 0,
		Y    : 0,
	}
	pty.Setsize(ptmx, size)
	// Simulate typing "echo Hello" into bash
	go func() {
		for {
			newcmd := <- in
			ptmx.Write(newcmd)
			//ptmx.Write([]byte("\n"))
		}
	}()

	// Read all output until bash exits
	output := make([]byte, MAX_TRANSFR)
	for {
		n, err := ptmx.Read(output)
		if n > 0 {
			//fmt.Print(string(output[:n]))
			out <- output[:n]
		}
		if err != nil {
			break
		}
	}
	
}