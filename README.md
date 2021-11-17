
# Distributed Mutual Exclusion

Implementation of a distributed mutual exclusion between nodes in a distributed system

---

## How to run

Before starting a node please read "Before running".

### Manually

To start a node, run the following in the node directory: 

> `go run . -name <name> -sport <port> -ips <other address>:<other port>`

### With Script
You can also start the program with a script for easy creation of nodes.
This can be done by running `run.ps1` in PowerShell in the root directory of the project.

You can specify how many nodes to start (supports up to 5) by writing:

> `run.ps1 -n <number of nodes>`

E.g., to start 4 nodes write:

> `run.ps1 -n 4`

(If no argument is given, it will start 3 nodes)

## Before running

Before starting a Node you must know what arguments to give the commandline.

The node must be run with the following required arguments:
- name
- server port of node
- ip addresses/ports of other nodes

#### Name

The name of the node must be unique between all nodes.

It can be specified with the `-name <name>` argument.

E.g., to start a node with the name `node0` do:

>`go run .-name node0`

#### Address
The address of the node is optional.

It can be specified with: `-address <address>`

E.g., to start a node with the address `128.0.0.1` do: 

> `go run . -address 128.0.0.1`

If no address argument is giving, it will pick a default value ("localhost")

#### Server port

The server port of the node must be unique.

This can be done with the `-sport <port>` argument.

E.g., to start a node on port `8080` do: 

> `go run . -sport 8080`

#### Ip addresses/ports of other nodes

The node must know the addresses and server ports of the other nodes.

This can be done with the `-ips <address>:<port>` argument.

Each peer node's ip address must be seperated by a comma i.e.: 
`-ips <address1>:<port1>,<address2>:<port2>,...`

You don't have to enter the whole ip address. 
You can simply enter the port numbers (the address will then be set to "localhost"):
`-ips <port1>,<port2>,...`

E.g., to start a node which have to know about another node having the ip address `128.0.0.1:8081` do:

> `go run . -ips 128.0.0.1:8081`

(You can mix and match the above for each ip address)

#### Start 3 Nodes Example

Let's say we have three nodes:
1. `node0` to run at `localhost:8080`
2. `node1` to run at `128.0.0.1:8081`
3. `node2` to run at `127.0.0.1:8083`

To run the nodes as separate go processes, enter the following in separate terminals in the node directory:
> `go run . -name node0 -sport 8080 -ips 128.0.0.1:8081,127.0.0.1:8083`

> `go run . -name node1 -address 128.0.0.1 -sport 8081 -ips 8080,127.0.0.1:8083`

> `go run . -name node2 -address 127.0.0.1 -sport 8083 -ips 8080,128.0.0.1:8081`

---

## Mandatory Exercise 2 - Distributed Mutual Exclusion

### Description

You have to implement distributed mutual exclusion between nodes in your distributed system.

You can choose to implement any of the algorithms, that were discussed in lecture 7.

---

### System Requirements

R1: Any node can at any time decide it wants access to the Critical Section

R2: Only one node at the same time is allowed to enter the Critical Section

R2: Every node that requests access to the Critical Section, will get access to the Critical Section (at some point in time)

---

### Technical Requirements

1. Use Golang to implement the service's nodes
2. Use gRPC for message passing between nodes
3. Your nodes need to find each other.  For service discovery, you can choose one of the following options
    1. supply a file with  ip addresses/ports of other nodes
    2. enter ip address/ports trough command line
    3. use the [Serf package](https://github.com/hashicorp/serf) for service discovery
4. Demonstrate that the system can be started with at least 3 nodes
5. Demonstrate using logs,  that a node gets access to the Critical Section

---

### Hand-in requirements

1. Hand in a single report in a pdf file
2. Provide a link to a Git repo with your source code in the report
3. Include system logs, that document the requirements are met, in the appendix of your report

---

### Grading notes

Partial implementations may be accepted, if the students can reason what they should have done in the report.

---
