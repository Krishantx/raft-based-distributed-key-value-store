# TODO

## KNOWN ERRORS AND BUGS (FIX LATER)
+ Error handling (Low Priority Do Log it though)
- [ ] When a key value pair does not exist for GET endpoint
- [ ] When trying to insert a key value pair in which a key already exists
- [ ] When a key does not exist for PUT
- [ ] When a key does not exits for a DELETE endpoint
- [ ] Fix race condition if 2 candidate requests hit @ the same time and can cause the term to break
- [ ] Use a database to store logs and indexes

## Proxy Node
- [x] Receives HTTP Requests from the client and forwards them to the raft nodes
- [x] Broadcasts the requests to all the nodes in the cluster (Hardcoded for now)

## Raft Node

### Leader Role
- [x] Send Heartbeat using gRPC protocol
- [x] (Using a .log file) Figure out a way to implement logs (Need these logs in persistent store, can write to a .log log file or use a SQLite database)
- [ ] implement a global log to log all operations in a file
### Candidate Role (Work on in the end - This is only required when things go wrong - Design for a perfect system first build fault tolerance later)


### Follower Role



### On Startup
- [ ] Load previous logs
