# TODO

## KNOWN ERRORS AND BUGS (FIX LATER)
+ Error handling (Low Priority Do Log it though)
- [ ] When a key value pair does not exist for GET endpoint
- [ ] When trying to insert a key value pair in which a key already exists
- [ ] When a key does not exist for PUT
- [ ] When a key does not exits for a DELETE endpoint
- [ ] Fix race condition if 2 candidate requests hit @ the same time and can cause the term to break
- [ ] Use a database to store logs and indexes
- [ ] I am currently checking weather the operation can be performed AFTER Appending to the logs. This is not incorrect but not isnt correct either because on paper when I reiterate through the logs I am going to arrive at the same state but I am storing more logs this way. I can skip saving the logs if the operation would give me an error for example PUT in a key that does not exist.
## Proxy Node
- [DONE] Receives HTTP Requests from the client and forwards them to the raft nodes
- [DONE] Broadcasts the requests to all the nodes in the cluster (Hardcoded for now)
- [MAYBE] Set up a service registry to make the nodes register to the proxy first


## Raft Node

### Leader Role
- [TODO] On every operation -> Append Entry to the RaftLog file as a source of truth. Only after the majority of the nodes append it.
- [DONE] Send Heartbeat using gRPC protocol
- [DONE] (Using a .log file) Figure out a way to implement logs (Need these logs in persistent store, can write to a .log log file or use a SQLite database)
- [DONE] implement a global log to log all operations in a file




### Candidate Role (Work on in the end - This is only required when things go wrong - Design for a perfect system first build fault tolerance later)


### Follower Role
- [DONE] Receives Heartbeats and Resets Timer to ensure leader is healthy
- [DONE] Change role to candidate once the leader does not send any heartbeats assume that it is dead


### On Startup
- [ ] Load previous logs
- [ ] So basically iterate through the logs and recreate whatever happened
- [ ] Step 1: Persistent Volumes for all docker images
- [ ] 


## Appending all logs
- [DONE] Step 1: Receives request from the proxy 
- [DONE] Step 2: Appends log to the RaftLog file
- [DONE] Step 3: Send a AppendEntry request to the follower nodes.
- [DONE] Step 4: Once majority nodes send confirmation of commitment => The leader commits it too.

## Candidate Role
* This will be the death of me
