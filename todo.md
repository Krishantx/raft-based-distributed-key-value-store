# TODO

## DO Next
- [ ] Get a kafka docker image.
- [ ] make the proxy a kafka producer, make it produce every request into the broker
- [ ] make the node a kafka consumer, make it consume every event in kafka.




- [ ] Fix race condition if 2 requests hit @ the same time and can cause the term to break
- - Use an event driven tool from proxy to the leader node. The leader node will use that subscribe to the kafka.

- [ ] Use a database to store logs and indexes
- - Use postgresql maybe? The rows(LOG) can be: Term, indexes, Operation, Key (NOT NULL), Value Can be structured schema.
- - Each Node has its own database. Slows down for writes but speeds up reads since reads are in memory.

- [ ] On boot up load the previous logs & come to a specific state and then join the raft
- - 




- [ ] Leader Role
- - [ ] Consumes messages from the event stream
- - [ ] Asks the followers to append logs to their database.
- - [ ] Commit to the database only after the data is in the hashmap memory
- - [ ] Once majority of followers commit. The leader will also commit.


- [ ] Candidate Role
- - [ ] Send a contesting in election. Also increases its term.
- - [ ] Votes for itself and waits for others to vote.
- - [ ] If it recieves a majority of votes it becomes the leader

- [ ] Follower Role
- - [ ] Recieves requests from the leader.
- - [ ] 
