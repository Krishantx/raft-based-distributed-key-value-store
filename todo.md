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
- [DONE] Load previous logs
- [DONE] So basically iterate through the logs and recreate whatever happened
- [DONE] Step 1: Persistent Volumes for all docker images
- [DONE] Iterate through the entire logs and modify the repository. 


## Appending all logs
- [DONE] Step 1: Receives request from the proxy 
- [DONE] Step 2: Appends log to the RaftLog file
- [DONE] Step 3: Send a AppendEntry request to the follower nodes.
- [DONE] Step 4: Once majority nodes send confirmation of commitment => The leader commits it too.

## Candidate Role
* This will be the death of me


## New Leader Election on dead leader nodes
- [ ] Step 1 - Detection on dead leaders, a follower once not received a timeout would assume that the leader is dead and become a candidate.
- [ ] Step 2 - The candidate would send a gRPC request to all the follower nodes asking them to vote.
- [ ] Step 3 - The followers would then compare the candidates indexes and term to its own index and term.
- [ ] Step 4 - A follower would only  vote for a new leader if the `candidateIndex >= followerIndex` AND  `candidateTerm >= followerTerm`. Regardless a node can only vote once per term, So even if it votes false. It would need to increment its term.
- [ ] Step 5 - If the candidate receives a majority of votes. It becomes the leader and working like a leader

## Technical Debt and Code Cleanup:
- [ ] Too many global vars are currently storing the same redundant information. These multiple widen possible places bugs can arise. Clean them up
- [ ] Heartbeat is just an empty append request. Can reduce the amount of code I have to maintain.
- [ ] Follower conttroller has 2 controllers `Server` and `RaftOperations` Server is just there for heartbeats. Can be merged into the same controller
