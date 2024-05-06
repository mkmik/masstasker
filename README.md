[![Go Reference](https://pkg.go.dev/badge/mkm.pub/masstasker.svg)](https://pkg.go.dev/mkm.pub/masstasker)
[![Go Report Card](https://goreportcard.com/badge/mkm.pub/masstasker)](https://goreportcard.com/report/mkm.pub/masstasker)

# Mass Tasker

MassTasker is a task queuing system with a small set of primitive operations. It enables primitive workflows thanks to atomic operations.

MassTasker has been inspired by the Workflow TaskMaster tool described in https://sre.google/sre-book/data-processing-pipelines/.

Core philosophy:

  * low-level but core support for atomic operations.
  * write clients in any programming language thanks a thin API based on Protobuf + gRPC.
  * workers can fail at any time; deal with it.
  * this is a building block; users are expected to design their task systems properly and handle errors.

MassTasker is mainly an API and a set of patterns.

You could probably achieve the same goals by using another task queue, such as Redis or even Postgres.
In fact, the MassTasker API can be implemented as a frontend to Redis, Postgres or other databases.

The primitives and the data model defined by MassTasker might help you focusing on the important aspects
of running a distributed job with tens of thousands of workers and encourages you to think about failure modes
that will occur in that regime.

At the same time, the MassTasker API will spare you from havingt to figure out how to properly implement a multi-key transactional write in redis or whatnot.

## API

The API is defined as a gRPC API.

This repo also contains an in-memory server implementing the API and a Go client that is a thin wrapper around
the raw gRPC and it includes a few helpers for common patterns.

The API is designed to fit on a page and be easy to understand, troubleshoot and operate even manually

## Data model

The one and only entity in the TaskMaster data model is the `Task`.

A task is defined as:

```proto
message Task {
  uint64 id = 1;
  string group = 2;
  google.protobuf.Any data = 3;
  google.protobuf.Timestamp not_before = 5;

  // optional error annotation, useful when moving tasks to an error queue.
  // If you need a more structured error please encode it in the payload.
  string error = 6;
}
```

Each task:

  * is immutable
  * has a unique ID
  * lives in exactly one task group
  * has a user defined payload (`data`)

## Operations

The MassTasker state consists of a set of immutable Tasks.

You can update the MassTasker state by creating new tasks and delete existing tasks.


There only two operations:

  * `Update`: updates the MassTasker state (creates and deletes tasks)
  * `Query`: Retrieves one task from a task group, possibly waiting for tasks to become available.

```proto
service MassTasker {
  rpc Update(UpdateRequest) returns (UpdateResponse);
  rpc Query(QueryRequest)   returns (QueryResponse);
}
```

### Update

The `Update` operation can atomically create new tasks and delete existing tasks. It is defined as:

```proto
message UpdateRequest {
  repeated Task created = 1;
  repeated uint64 deleted = 2;
  repeated uint64 predicates = 3;
}
```

If any of the tasks whose IDs are referenced in the `predicates` argument do not exist, the request returns the `FAILED_PRECONDITION` error.

If any of the tasks whose IDs are referenced by the `deleted` argument do not exist, the request returns the `NOT_FOUND` error.

The `predicates` argument is checked first, you're guaranteed to receive a `FAILED_PRECONDITION` if you pass the same task to both the `deleted` and `predicates` arguments.

### Query

The `Query` operation searches for one task from the specified task group whose `not_before` field is less than _now_.

The operation is defined as:

```proto
message QueryRequest {
  string group = 1;
  google.protobuf.Duration own_for = 3;
  bool wait = 5;

  google.protobuf.Timestamp now = 4; // if omitted it defaults to server clock's now.
}
```

If a non-zero value for `own_for` is provided, the task is atomically deleted and its content moved to a new task (with a new ID).
The `not_before` field of this new task is defined as `old_not_before + own_for`.

If no task is found, the `NOT_FOUND` error is returned, unless the `wait` parameter is set to true, in which case the server will find the closest
task in the _future_, compute how much time the caller has to wait and blocks the gRPC connection until the task will become available.

Since `Query` always atomically deletes and re-creates the tasks, at any given time only one task will effectively be the **owner** of a given task.
The worker that obtained the older version of the task now holds a task ID that no longer exists and subsequent `Update` operation will fail when the worker will try to delete the task or use it as a predicate id.

This property provides correctness even in the case that a worker doesn't respect the lease time (clocks are brittle).

### Admin/Debug APIs

There are also other API calls defined in this package that allow the user to inspect the content of a task group, or to perform bulk operations
such as re-setting the `not_before` fields or deleting all tasks in a large task group. These operations are not guaranteed to be atomic.

## Payload

The user can put some user defined payload in the `data` field, which is meant to contain a Protobuf [Any](https://protobuf.dev/programming-guides/proto3/#any) message. `Any` messages can carry a user-defined protobuf message alongside with the package+name of the user defined message type.

User defined workers must unmarshal and remarshal when needed. Workers can mix different message types in the same task group (and dispatch the action depending on the type information present in the `Any` metadata) or they can organize their work items so that each working group contains only tasks
with a uniform message type.

## Common patterns

The MassTasker primitives are very low level and users are required to design their tasking model the way they see fit.
That said, there are some common primitive patterns that you can re-use.

### Task lease

The most common operation in a task queue system is to deliver tasks to workers in such a way that each worker
processes one task at a time.

To achieve that, a worker issues a `Query` request and specifies an `own_for` parameter that defines for how long no query operation will
return this task to any worker.

The `own_for` duration is generally choosen to be greater than the internal timeout the worker would apply to itself while processing the task.

If a worker fails ungracefully, the task(s) owned by it will not be available to other workers for the duration of the lease.

### Committing a task

Once a worker is done processing a task it can simply delete the task or _move_ it into another task queue (move == delete+create).

The "commit" will fail if the _least_ has expired and another worker has grabbed the task.

### Bootstrap task

In order for workers to process tasks, something has to write the tasks into MassTasker in the first place.

You can have an external orchestrator creating all the tasks or you can have the external orchestrator only create a small 
"bootstrap" task and then leave it the first worker that picks it up to create all the other tasks.

For example, when a worker processes the bootstrap task, it may list all files in an objecstore bucket and create
a new task for each file. Upon committing the list of new tasks the bootstrap task itself will be removed atomically.

The provided MassTasker server offers a flag/env config that will define the name of a group where the server itself
will create an initial bootstrap task (with empty payload). This can be useful to bootstrap a system without external orchestrators

### Config task

You can coordinate workers by having them all share a _config task_.

The simplest way is to have a singleton `Config` task in a `config` task group and have each worker `Query` without the `own_for` parameter.

The worker will then grab a task and start working on it. At task commit time, the config task ID will be passed in the `predicates` parameters.

If the config has been changed (e.g. by a human-driven config push, or another task), the commit will fail, and the worker will re-start with the new config. 

### State/Barrier tasks

Another pattern similar to the Config task pattern is to have a singleton task but instead of it being a read-only task, you can have multiple tasks
race to update it and retry on precondition failed error.

This allows you to implement barriers.

For example imagine that your job needs to read N files from object store, chop each file up into M ranges and process each range in parallel and then re-join them once every range has been processed. There are many ways to do this, but one way is to create a join/barrier task that contains some metadata about the progress of the tasks that processed each range. For example when the worker that processes a range is ready to commit the processing task, it will fetch the barrier task, add the URL of the processed chunk to the barrier task issue an update request that deletes the processing task, the old barrier task, creates the new updated barrier task and uses the barrier task as predicate:

```go
myUpdateBarier(barrierTask, myTask) // encode e.g. the output filename from from myTask into barrierTask
mt.Update(ctx, &UpdateRequest(Created: {barrierTask}, Deleted: {myTask.ID, barrierTask.ID}, Predicate {barrierTask.ID}))
```

and retry with exp backoff until the error returned by Update is not `FAILED_PRECONDITION`.

When the state encoded in the barrier task shows that all ranges from the input file have been processed, the worker can then enqueue
a new task that takes care of the next phase in the pipeline.

You can see this as a way to _transfer_ the information from independent task into a a part of the barrier task and vice versa.
A common trick is to reuse the very same protobuf messages.

If a single barrier task cannot fit all the state, you can split it up in multiple barrier tasks living in multiple task groups, each
handling a subset of the tasks (by levering some form of hashing for example).

### Dead letter queue

In general, task procesing errors fall in two broad categories:

  * recoverable errors
  * unrecoverable errors

Recoverable errors can be transient errors such as networking issues, but also errors that don't necessarily resolve by themselves, such as misconfigurations. When recoverable errors are encountered, it's a good idea to just let the worker crash and let the task in the task group
ready to be picked up later by another worker.
This will naturally throttle the job progress and is suitable for handling with downstream rate-limiting / backpressure.
Crashing the workers is an effective way to surface configuration errors.
Once the configuration errors is fixed the job will just continue processing the tasks that have been not processed yet.

Unrecoverable errors OTOH are errors such as _bad data_ for which it doesn't make sense to block the job for.
Instead, you may want to save tasks that have failed into another group and possibly run another job on those tasks later.

These errors may not strictly be un-recoverable in principle. For example, you might encounter data that is not _yet_ supported
by your worker but you may add support for that case (or fix the bug) and re-run the job only on the failed data.

Another common example is handling outliers that would OOM the workers.
Instead of sizing the entire worker pool so it can fit the worst case task size,
you can reject large tasks and put them into another group where a smaller set of workers with more available resources can process them.

The `Task` message contains a optional `error` field that is mean to contain an error message suitable for annotating tasks that have
been moved to such a dead-letter-queue.
If you need more than an error string (such as an error code) you probably should define that in your custom payload message types.

### Monitoring progress

The provided MassTasker server implementation exposes a metric containing a gauge that exposes the task count metric for each task group.

A common pattern to monitor progress of your job is to monitor the task group task count and compute the processing speed by taking the derivative of
the task count (and then computing the ETA by computing how long it will take to reach to zero).

This works well when there are many and relatively short tasks. For example if you have 1000 tasks and 1000 workers each processing one task, then you won't see much of a progress report until the whole job is done.
