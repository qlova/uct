# URI compatibillty table
If any specific URI is not yet supported by the target langauge you want to use, report it as a bug.

## TCP  
    tcp://address:port
    
    eg.
    LOAD "tcp://:8080" //Start listening on port 8080 as a TCP server.
    OPEN "tcp://:8080" //Return a pipe to the next client that connects.
    OPEN "tcp://localhost:8080" //Return a pipe that connects to localhost:8080.
**Supported by:** Go, Python and Java.

## UDP
    udp://address:port
    
    LOAD "udp://:7777" //Start listening on port 7777 as a UDP server.
    OPEN "udp://:7777" //Return a pipe to the next client that connects each call to IN/OUT will recieve/send a single packet..
    OPEN "udp://localhost:8080" //Return a pipe that connects to localhost:7777 each call to OUT will be sent as a single packet.
**Planned Feature**

## Multicast
    multicast://group:port
      
    OPEN "multicast://224.3.29.71:10000" //Return a pipe connected to the multicast group, each call to IN/OUT will recieve/send a single packet..
**Supported by:** Java

## HTTP
    http://domain.com/path
      
    OPEN "http://github.com" //Return a pipe connected to a web address.
**Supported by:** Go

## DNS
    dns://domain.com
    
    LOAD "dns://github.com" //Return an array with hosts seperated by spaces.
**Supported by:** Go, Java, Python

## DIR
    dir://path/to/dir
      
    OPEN "dir:///" //Return a pipe connected to the specified directory, each call to IN will return a file or directory within the specified path.
**Supported by:** Go

