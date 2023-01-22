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

Bulk(orange) vs Original(blue) implemntation . Value is in nanoseconds . Hence lower the better 

# 1 repeat cycle

![Screenshot 2023-01-22 at 2 25 42 PM](https://user-images.githubusercontent.com/41006458/213907837-6a5b2ab4-c433-40cc-ae1e-6fc9d59c4a40.png)

![Screenshot 2023-01-22 at 2 26 28 PM](https://user-images.githubusercontent.com/41006458/213907863-dc2cc441-edf1-4bbd-9455-598153f83586.png)

# 30 repeat cycle

![Screenshot 2023-01-22 at 2 27 05 PM](https://user-images.githubusercontent.com/41006458/213907881-a6256ede-524d-45c2-8aeb-17d52a9686cf.png)

![Screenshot 2023-01-22 at 2 26 53 PM](https://user-images.githubusercontent.com/41006458/213907876-7460b1a1-5dad-4dfa-9dae-b68d4901e12c.png)

# 70 repeat cycle

![Screenshot 2023-01-22 at 2 27 35 PM](https://user-images.githubusercontent.com/41006458/213907896-53cc5c7f-74f9-49cd-a993-584ee8fd4fde.png)

![Screenshot 2023-01-22 at 2 28 03 PM](https://user-images.githubusercontent.com/41006458/213907916-8f8dc933-199d-4df5-9c4a-3da1ad262b33.png)



