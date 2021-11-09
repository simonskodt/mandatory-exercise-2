
# Mandatory Exercise 2 - Distributed Mutual Exclusion

Implementation of a distributed mutual exclusion between nodes in a distributed system

---

## Description

You have to implement distributed mutual exclusion between nodes in your distributed system.

You can choose to implement any of the algorithms, that were discussed in lecture 7.

---

## System Requirements

R1: Any node can at any time decide it wants access to the Critical Section

R2: Only one node at the same time is allowed to enter the Critical Section

R2: Every node that requests access to the Critical Section, will get access to the Critical Section (at some point in time)

---

## Technical Requirements

1. Use Golang to implement the service's nodes
2. Use gRPC for message passing between nodes
3. Your nodes need to find each other.  For service discovery, you can choose one of the following options
    1. supply a file with  ip addresses/ports of other nodes
    2. enter ip adress/ports trough command line
    3. use the Serf package for service discovery
4. Demonstrate that the system can be started with at least 3 nodes
5. Demonstrate using logs,  that a node gets access to the Critical Section

---

## Hand-in requirements

1. Hand in a single report in a pdf file
2. Provide a link to a Git repo with your source code in the report
3. Include system logs, that document the requirements are met, in the appendix of your report

---

## Grading notes

Partial implementations may be accepted, if the students can reason what they should have done in the report.

---
