package main 

import (
	"log"
	"os/exec"
	"bufio"
	"github.com/creack/pty"
	"os"
)

func ExecCmd(data []byte, out chan []byte){

	err := os.WriteFile("s.sh", data, 0777)
	if err != nil {
		panic(err)
	}
	defer os.Remove("s.sh")

	cmd := exec.Command("sh", "s.sh")
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

func StartShell(out chan []byte, in chan []byte){

	cmd := exec.Command("bash")

	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Println(err)
		return
	}
	defer ptmx.Close()

	// Simulate typing "echo Hello" into bash
	go func() {
		for {
			newcmd := <- in
			ptmx.Write(newcmd)
			ptmx.Write([]byte("\n"))
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