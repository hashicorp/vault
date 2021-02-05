/*
Package workerpool queues work to a limited number of goroutines.

The purpose of the worker pool is to limit the concurrency of tasks
executed by the workers.  This is useful when performing tasks that require
sufficient resources (CPU, memory, etc.), and running too many tasks at the
same time would exhaust resources.

Non-blocking task submission

A task is a function submitted to the worker pool for execution.  Submitting
tasks to this worker pool will not block, regardless of the number of tasks.
Incoming tasks are immediately dispatched to an available
worker.  If no worker is immediately available, or there are already tasks
waiting for an available worker, then the task is put on a waiting queue to
wait for an available worker.

The intent of the worker pool is to limit the concurrency of task execution,
not limit the number of tasks queued to be executed.  Therefore, this unbounded
input of tasks is acceptable as the tasks cannot be discarded.  If the number
of inbound tasks is too many to even queue for pending processing, then the
solution is outside the scope of workerpool, and should be solved by
distributing load over multiple systems, and/or storing input for pending
processing in intermediate storage such as a database, file system, distributed
message queue, etc.

Dispatcher

This worker pool uses a single dispatcher goroutine to read tasks from the
input task queue and dispatch them to worker goroutines.  This allows for a
small input channel, and lets the dispatcher queue as many tasks as are
submitted when there are no available workers.  Additionally, the dispatcher
can adjust the number of workers as appropriate for the work load, without
having to utilize locked counters and checks incurred on task submission.

When no tasks have been submitted for a period of time, a worker is removed by
the dispatcher.  This is done until there are no more workers to remove.  The
minimum number of workers is always zero, because the time to start new workers
is insignificant.

Usage note

It is advisable to use different worker pools for tasks that are bound by
different resources, or that have different resource use patterns.  For
example, tasks that use X Mb of memory may need different concurrency limits
than tasks that use Y Mb of memory.

Waiting queue vs goroutines

When there are no available workers to handle incoming tasks, the tasks are put
on a waiting queue, in this implementation.  In implementations mentioned in
the credits below, these tasks were passed to goroutines.  Using a queue is
faster and has less memory overhead than creating a separate goroutine for each
waiting task, allowing a much higher number of waiting tasks.  Also, using a
waiting queue ensures that tasks are given to workers in the order the tasks
were received.

Credits

This implementation builds on ideas from the following:

http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang
http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html

*/
package workerpool
