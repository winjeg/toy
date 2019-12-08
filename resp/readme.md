# RESP 
RESP is *Redis Protocol specification*  

### Networking layer  
A client connects to a Redis server creating a TCP connection to the port 6379.  
While RESP is technically non-TCP specific, in the context of Redis the protocol is only used with TCP connections (or equivalent stream oriented connections like Unix sockets).

### Request-Response model  
- Redis supports pipelining (covered later in this document). So it is possible for clients to send multiple commands at once, and wait for replies later.  
- When a Redis client subscribes to a Pub/Sub channel, the protocol changes semantics and becomes a push protocol, that is, the client no longer requires sending commands, because the server will automatically send to the client new messages (for the channels the client is subscribed to) as soon as they are received.  

### RESP protocol description
- For Simple Strings the first byte of the reply is "+"  
- For Errors the first byte of the reply is "-"  
- For Integers the first byte of the reply is ":"  
- For Bulk Strings the first byte of the reply is "$"  
- For Arrays the first byte of the reply is "*"  

```
*3\r\n
$3\r\n
foo\r\n
$-1\r\n
$3\r\n
bar\r\n
```
The second element is a Null. The client library should return something like this:  
`["foo",nil,"bar"]`

#### example
```
C: *2\r\n
C: $4\r\n
C: LLEN\r\n
C: $6\r\n
C: mylist\r\n

S: :48293\r\n
```

As usual we separate different parts of the protocol with newlines for simplicity, but the actual interaction is the client sending `*2\r\n$4\r\nLLEN\r\n$6\r\nmylist\r\n` as a whole.