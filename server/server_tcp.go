package main

import (
	"log"
	"net"
	"errors"
    "io"
    "fmt"
    "bufio"
    _"strings"
    "os"
    //"time"
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
            continue
        }

        cnxLock.Lock()
        if !tcpCnxExist {
            tcpCnxExist = true
            //currentServerTcp = &conn
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
        for {
            log.Println("shell>")
            var cmd string
            reader := bufio.NewReader(os.Stdin)
            cmd, _ = reader.ReadString('\n')
            //cmd = strings.TrimSpace(cmd)
            inchannel <- []byte(cmd)
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
                fmt.Println(string(buf))
            }
        }
    }()

    for{

        cnxLock.Lock()
        if !tcpCnxExist {
            cnxLock.Unlock()
            return
        }
        cnxLock.Unlock()

        data := <- inchannel
        _, err := conn.Write(data)
        if err != nil {
            log.Println("could not write to target server", err)
            return
        }
        //log.Println("wrote to target server", n, "bytes")

        //time.Sleep(time.Second * 30)
    }

}