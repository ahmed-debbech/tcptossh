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
		return
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}
	defer stderr.Close()

	go func() {
        scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanLines)
        for scanner.Scan() {
            //out <- append(scanner.Bytes(), '\n')
			out <- []byte(scanner.Text())
		}	
    }()

	go func() {
        scanner := bufio.NewScanner(stderr)
		scanner.Split(bufio.ScanLines)
        for scanner.Scan() {
            //out <- append(scanner.Bytes(), '\n')
			out <- []byte(scanner.Text())
		}	
    }()

	if err := cmd.Start(); err != nil {
		log.Println("could not enable run command", err)
		return
	}

	cmd.Wait()

}