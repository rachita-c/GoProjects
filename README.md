# GoProjects

Playing around with Go and Distributed Systems!

Testing out:
- Boostrap server and client connections
- Implementing heartbeat messages, panic, defer functions; If primary server is down, secondary bootstrap server
 doesn't receive callback and calls the panic function that does dynamic DNS mapping and swaps the primary and secondary
 servers.
- Mutex locking for resources.
- Consensus algorithm integration (Paxos) with Boostrap Logic.

