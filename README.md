# TCP2SSH

Simple way to avoid all Nat/Port Forwarding/Upnp or any other thing that may block you to make any machine expose a bash for you. 

> [!CAUTION]  
> Use at your own risk. Anything may fail at any time. I made this as a tool for me to be able to access some of my machines that were blocked by a NAT.

### How it works

It reverses a tcp connection made from your TARGET server to your actual CONTROLING server and exposes a simple bash terminal 

* Same TCP connection
* Reversed client and server

[You] <===> [Client terminal][Server] <====> [Target Server]

### How use?

* `go build .` for both server dir and target_server dir
* `go run . your_key*32char` for server to run command cli and initiate the listening port
* `go run . you_server_ip same_key*32char` for target server to keep trying to make request to your server, once succeeded we can try to reverse the connection from the server and bash through it

### Limitations

* Only single bash is created when the target server app runs. All TCP connection use that same bash session.
* Only single connection from target server can be handled on the server side at a time.