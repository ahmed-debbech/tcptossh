# TCP2SSH

### How it works

* Same TCP connection
* Reversed client and server

[You] <===> [Client terminal][Server] <====> [Target Server]


### Limitations

* Only single bash is created when the target server app runs. All TCP connection use that same bash session.
* Only single connection from target server can be handled on the server side at a time.