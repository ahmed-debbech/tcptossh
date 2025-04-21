package main 

import (
	"log"
	"os/exec"
	"bufio"
)

func ExecCmd(data []byte, out chan []byte){

	cmd := exec.Command("sh", "-c", string(data))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stdout.Close()
	

	go func() {
        scanner := bufio.NewScanner(stdout)
        for scanner.Scan() {
            out <- []byte(scanner.Text())
        }	
    }()

	if err := cmd.Start(); err != nil {
		log.Println("could not enable run command", err)
	}

	cmd.Wait()

}