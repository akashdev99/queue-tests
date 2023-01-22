# Bench-Bulk-Queueimpl8
The original implementation uses interface and does not provide bulk push or bulk pop capabilities . Hence , you will have to run Push/Pop in a loop in order to get “X” number of operations(push or pop) performed .

As an enhancement to queueimpl7 — queueimpl8 , i have added the following :

Added generics support
Added queue bulk push/queue bulk pop functionality . (Enqueue/Dequeue)


## Tests
There is 1 distinct benchmark tests:
1) Adds N values to the queue in batch and then remove all afterwards in bulk

We also have a repeat cycle added to simulat M pusp and pop cycles within an iteration to simulate the behaviour in producer consumer architecture . repeat =1 , negates this assumption 

## Results

