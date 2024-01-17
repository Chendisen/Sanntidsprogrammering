3: GOMAXPROCS seems to decide how many "threads" that can be operated at the same time. So when we initialize it with 1, it is only able to run one thread at a time, removing the problems with race condition. When sharing the variable poorly, there will appear data race, so that both functions receive the variable with the same initial value at a given iteration of the loop, but the last function to send back their result will overwrite the other functions result. 

4: We should choose the MUTEX over the semaphore since the mutex is a binary flag that only allows the thread that decremented it to increment it again. This will prevent thread 2 from unlocking if thread 1 is the one that locked, and vice versa. 

Some questions:

1: Concurrency is swapping back and fourth really fast between threads so that the operations look like they are being executed in parallell, while parallelism is actually running multiple operations at the same time. 

2: Race condition is when the order of events affects the correctness of the program. Data race is a more specified occurence where several threads access the same variable without synchronisation such that the outcome is wrong. 

3: A scheduler selects among available threads and executes a swap of threads either by cooperative scheduling, where the threads yield to each other in order to ensure that the code finishes in time. It can also use preemptive scheduling where the threads are changed by force after an amount of time, or we can use a mix. 

4: We want to use multiple threads in order to execute programs in an concurrent order rather than having everyting sequentially. Several threads allows us to have a parallellistic program even though it is sequential through concurrent programming. 

5: 

6: Concurrency makes the programmers life both harder and easier. Harder because there is more stuff to implement and account for, but easier since some tasks become more efficient. 

7: We prefer message passing, because it seems more organized, and there seems like less chance of race conditions/data race, which is worth the extra lines of code. 