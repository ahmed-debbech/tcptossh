package main

import (
	"log"
	"net"
	"errors"
    "io"
    "fmt"
    _"strings"
    "os"
    "golang.org/x/term"
)

func startServertoServer() error {


	ln, err := net.Listen("tcp", ":42001")
    if err != nil {
        log.Println(err)
        return errors.New("not able to start server socket")
    }

    log.Println("Started Server <---> Target Server TCP connection! on 42001")

    for {  

        conn, err := ln.Accept()
        if err != nil {
            log.Println(err)
            panic("PANIC")
        }

        cnxLock.Lock()
        if !tcpCnxExist {
            tcpCnxExist = true

            log.Println("New connection is accepted")
            go handleServer(conn)
        }else{
            conn.Close()
        }
        cnxLock.Unlock()
    }
    close(inchannel)

    return nil
}

func handleServer(conn net.Conn) {
	defer conn.Close()
    defer func(){
        cnxLock.Lock()
        tcpCnxExist = false
        cnxLock.Unlock()
    }()
    
    go func(){

        oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
        if err != nil {
            panic(err)
        }
        defer term.Restore(int(os.Stdin.Fd()), oldState)
    
        fmt.Println("Press keys (ESC to exit):")
    
        buf := make([]byte, 1)
        for {
            _, err := os.Stdin.Read(buf)
            if err != nil {
                break
            }

            if buf[0] == 27 { // ESC key
                break
            }
            //fmt.Print(string(buf))
            inchannel <- []byte(buf)

        }
    }()

    go func(){
        for{
            buf := make([]byte, MAX_TRANSFR)
            n, err := conn.Read(buf)
            if err != nil {
                if err == io.EOF {
                    log.Println("TCP socket from target server closed!", err)
                    cnxLock.Lock()
                    tcpCnxExist = false
                    cnxLock.Unlock()
                    return;
                }
            }


            if n > 1 {
                ll := buf[:n]
                l := PipeMessages(string(ll))
                for k := 0; k<=len(l)-1; k++{
                    text, err := Decrypt([]byte(key), string(l[k][1:len(l[k])-1]))
                    if err != nil {
                        log.Println("error decrypting", err)
                        panic("PANIC")
                    }
                    fmt.Print(string(text))
                }
            }
        }
    }()

    for{

        cnxLock.Lock()
        if !tcpCnxExist {
            cnxLock.Unlock()
            panic("PANIC")
        }
        cnxLock.Unlock()

        data := <- inchannel
        
        cipher, err := Encrypt([]byte(key), string(data))
        if err != nil {
            log.Println("error encrypting", err)
            panic("PANIC")
        }

        _, err = conn.Write([]byte(cipher))
        if err != nil {
            log.Println("could not write to target server", err)
            panic("PANIC")
        }
    }

}

func PipeMessages(line string) []string{
    ss := make([]string, 0)
    
    opening := false

    s := ""
    for i := 0; i<=len(line)-1; i++{
        if line[i] == '[' && !opening {
            opening = true
            s = ""
        }
        if line[i] == ']' && opening {
            s += string(line[i])
            ss = append(ss, s)
            opening = false
        }
        if line[i] != ']' && opening {
            s += string(line[i])
        }
    }
    return ss
}