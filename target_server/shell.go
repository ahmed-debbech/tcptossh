package main 

import (
	"log"
	"os/exec"
	"bufio"
	"os"
)

func ExecCmd(data []byte, out chan []byte){

	err := os.WriteFile("s.sh", data, 0777)
	if err != nil {
		panic(err)
	}
	defer os.Remove("s.sh")

	cmd := exec.Command("sh", "s.sh")
	log.Println(cmd)
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
			out <- scanner.Bytes()
		}	
    }()

	go func() {
        scanner := bufio.NewScanner(stderr)
		scanner.Split(bufio.ScanLines)
        for scanner.Scan() {
            //out <- append(scanner.Bytes(), '\n')
			out <- scanner.Bytes()
		}	
    }()

	if err := cmd.Start(); err != nil {
		log.Println("could not enable run command", err)
		return
	}

	cmd.Wait()

}