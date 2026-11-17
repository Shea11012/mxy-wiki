---
source_url: "https://www.youtube.com/watch?v=fwHok9ZhQaY"
title: "sync.Pool in Go Explained & Avoid The Heap Allocation Mistake"
author: "Phuong Le"
source_kind: third-party
upload_date: 2025-11-06
ingested: 2026-10-07
format: webvtt
subtitle_language: en
extracted_to:
  - "sync.Pool"
---
WEBVTT
Kind: captions
Language: en

00:00:00.640 --> 00:00:01.520
Alright.

00:00:01.520 --> 00:00:05.200
Let's talk about what sync.Pool&nbsp;
is, why it matters for performance,

00:00:05.200 --> 00:00:08.080
when it clears its objects,&nbsp;
and how to use it properly,

00:00:08.080 --> 00:00:10.880
so you don't fall into the heap allocation trap.

00:00:10.880 --> 00:00:13.680
To understand why sync.Pool exists,

00:00:13.680 --> 00:00:17.280
we need to talk about the problem heap allocation.

00:00:17.280 --> 00:00:18.640
Heap allocation happens

00:00:18.640 --> 00:00:22.480
whenever a value needs to live longer&nbsp;
than the function that created it.

00:00:22.480 --> 00:00:25.680
These heap objects appear in every Go program,

00:00:25.680 --> 00:00:28.800
and working with them comes with two main costs.

00:00:28.800 --> 00:00:30.560
First, allocation.

00:00:30.560 --> 00:00:33.520
That's the process of reserving&nbsp;
memory space for the objects.

00:00:33.520 --> 00:00:34.240
It's pretty fast,

00:00:34.240 --> 00:00:39.120
but still much slower than working with&nbsp;
objects that are already on the stack.

00:00:39.120 --> 00:00:41.360
Second, deallocation.

00:00:41.360 --> 00:00:43.280
That's when the memory gets reclaimed.

00:00:43.280 --> 00:00:46.640
Once the object is no longer in use or reachable,

00:00:46.640 --> 00:00:49.200
Go doesn't require you to manually free heap

00:00:49.200 --> 00:00:52.000
memory like you would with the free function in C.

00:00:52.000 --> 00:00:54.320
Instead, it uses a garbage collector.

00:00:54.320 --> 00:00:57.920
The garbage collector&nbsp;
automatically marks live objects,

00:00:57.920 --> 00:01:01.600
and any objects that are not&nbsp;
marked are considered garbage.

00:01:01.600 --> 00:01:04.720
The garbage collector then&nbsp;
reclaims that memory for you.

00:01:04.720 --> 00:01:07.040
The trade-off is extra overhead.

00:01:07.040 --> 00:01:11.760
It uses CPU time to iterate, mark&nbsp;
objects, and sweep the unused ones.

00:01:11.760 --> 00:01:16.800
It also enables write barriers, which means&nbsp;
updates to global and heap pointers are tracked.

00:01:16.800 --> 00:01:20.640
So those writes become a bit&nbsp;
slower, and it briefly pauses

00:01:20.640 --> 00:01:24.240
the entire program during those&nbsp;
garbage collection phase transitions.

00:01:24.240 --> 00:01:28.000
This is what's commonly&nbsp;
called a stop-the-world pause.

00:01:28.000 --> 00:01:32.000
The solution for this overhead&nbsp;
is obviously sync.Pool.

00:01:32.000 --> 00:01:35.920
It's basically a cache for&nbsp;
objects managed by the Go runtime.

00:01:35.920 --> 00:01:40.640
It stores the objects you give it and&nbsp;
hands them back to whoever needs them.

00:01:40.640 --> 00:01:45.520
Goroutines can put objects into the pool&nbsp;
and let sync.Pool manage them for you.

00:01:45.520 --> 00:01:47.200
You don't really know what's going on inside

00:01:47.200 --> 00:01:51.280
or how many objects are in&nbsp;
there, and honestly, that's fine.

00:01:51.280 --> 00:01:53.520
Later, you can get an object back and reuse it

00:01:53.520 --> 00:01:56.960
instead of making a fresh allocation each time.

00:01:56.960 --> 00:02:00.480
And it's designed to be safe&nbsp;
and efficient with concurrency,

00:02:00.480 --> 00:02:06.000
so multiple goroutines can grab from&nbsp;
or return to the pool at the same time.

00:02:06.000 --> 00:02:08.240
Besides reducing the load on the garbage collector

00:02:08.240 --> 00:02:12.800
by reusing objects, it also manages&nbsp;
the lifecycle of what it holds.

00:02:12.800 --> 00:02:18.400
Your objects don't live there forever,&nbsp;
and we'll find out why in just a moment.

00:02:18.400 --> 00:02:21.120
Using sync.Pool is pretty simple.

00:02:21.120 --> 00:02:23.280
There are just three things.

00:02:23.280 --> 00:02:25.440
First, you don't need a constructor.

00:02:25.440 --> 00:02:28.720
Just declare a variable of type sync.Pool,

00:02:28.720 --> 00:02:31.760
and don't forget to import&nbsp;
the sync standard package.

00:02:31.760 --> 00:02:33.120
If you want it to automatically

00:02:33.120 --> 00:02:36.960
create objects when the pool is empty,&nbsp;
which is what you'll usually do,

00:02:36.960 --> 00:02:41.440
you set its New field to a function that&nbsp;
returns the kind of object you want.

00:02:41.440 --> 00:02:45.600
Next, getting objects, you call&nbsp;
the Get function on the pool.

00:02:45.600 --> 00:02:48.480
It gives back an empty interface, so you'll almost

00:02:48.480 --> 00:02:52.080
always do a type assertion to the type you expect.

00:02:52.080 --> 00:02:55.600
It's also a good idea to reset&nbsp;
the object before you use it,

00:02:55.600 --> 00:02:58.720
just to clear out any leftover&nbsp;
state from its previous use.

00:02:58.720 --> 00:03:02.400
And when you're done, you call the put&nbsp;
function to return the object to the pool,

00:03:02.400 --> 00:03:06.080
so it can be reused. Behind the&nbsp;
scenes, when you call the get function,

00:03:06.080 --> 00:03:08.960
the pool first tries to&nbsp;
hand you an existing object.

00:03:08.960 --> 00:03:12.480
If there isn't one, it calls&nbsp;
the New function you defined,

00:03:12.480 --> 00:03:16.080
and if you didn't provide a new&nbsp;
function, get will just return nil.

00:03:16.080 --> 00:03:19.600
The typical reason you'd skip providing a&nbsp;
New function is when you want to control

00:03:19.600 --> 00:03:21.280
how the object gets created.

00:03:21.280 --> 00:03:24.160
If you look at the sync.Pool struct definition,

00:03:24.160 --> 00:03:27.360
the New field is a function&nbsp;
that takes no parameters.

00:03:27.360 --> 00:03:30.320
So if your object needs&nbsp;
parameters to be initialized,

00:03:30.320 --> 00:03:32.000
you'll need to handle that manually.

00:03:32.000 --> 00:03:35.440
There's actually a good example&nbsp;
of not providing a New field

00:03:35.440 --> 00:03:38.400
in the standard library's HTTP package.

00:03:38.400 --> 00:03:41.520
It defines its own new&nbsp;
function that takes arguments.

00:03:41.520 --> 00:03:46.080
This function wraps the get-from-pool&nbsp;
logic to handle custom object creation,

00:03:46.080 --> 00:03:51.520
based on the argument that was passed in, and&nbsp;
checks when the pool's Get method returns nil.

00:03:51.520 --> 00:03:56.800
If there's no object available, it just&nbsp;
allocates a new one using that argument.

00:03:56.800 --> 00:03:58.320
All right, that's the basics.

00:03:58.320 --> 00:04:01.680
If you take a closer look at&nbsp;
how sync.Pool is usually used,

00:04:01.680 --> 00:04:05.920
you'll notice that the New field&nbsp;
for the most part returns a pointer.

00:04:05.920 --> 00:04:09.120
That's intentional because it&nbsp;
helps avoid extra allocations.

00:04:09.120 --> 00:04:11.120
It's a bit out of scope and a little advanced,

00:04:11.120 --> 00:04:15.200
so I'll save that explanation&nbsp;
for the end of the video.

00:04:15.200 --> 00:04:18.960
Okay, then when should we actually use sync.Pool?

00:04:18.960 --> 00:04:21.520
sync.Pool is designed for short-lived objects

00:04:21.520 --> 00:04:23.920
that are created and discarded frequently.

00:04:23.920 --> 00:04:29.280
A common example is per-operation data,&nbsp;
like when you handle a single HTTP

00:04:29.280 --> 00:04:34.960
request. The goroutine gets an object from&nbsp;
the pool, uses it, and modifies its state.

00:04:34.960 --> 00:04:38.480
Once it's done, it puts the dirty&nbsp;
object back into the pool and responds

00:04:38.480 --> 00:04:40.320
to the client.

00:04:40.320 --> 00:04:44.320
After that, another request comes&nbsp;
in, and another goroutine pulls

00:04:44.320 --> 00:04:49.360
that dirty object out, resets&nbsp;
its state, and uses it again.

00:04:49.360 --> 00:04:53.360
Since these objects get reused so&nbsp;
quickly, sync.Pool is perfectly suited

00:04:53.360 --> 00:04:55.680
for this kind of use case.

00:04:55.680 --> 00:04:59.360
Now why short-lived and frequently used objects?

00:04:59.360 --> 00:05:01.760
That's because if your objects&nbsp;
sit in the pool for too long

00:05:01.760 --> 00:05:05.680
without being reused, the pool&nbsp;
will eventually discard them,

00:05:05.680 --> 00:05:08.080
and they'll be garbage collected.

00:05:08.080 --> 00:05:10.640
But how long is long enough?

00:05:10.640 --> 00:05:14.240
The sync.Pool documentation says&nbsp;
that any item stored in the pool

00:05:14.240 --> 00:05:18.560
may be removed automatically at&nbsp;
any time without notification.

00:05:18.560 --> 00:05:22.800
Though, to be fair, that&nbsp;
explanation feels a bit uncertain.

00:05:22.800 --> 00:05:26.240
It's actually quite complicated&nbsp;
under the hood, but for our purposes,

00:05:26.240 --> 00:05:28.000
let's try to simplify it.

00:05:28.000 --> 00:05:31.440
The lifecycle of objects in&nbsp;
the pool is basically tied

00:05:31.440 --> 00:05:33.280
to the garbage collection cycle.

00:05:33.280 --> 00:05:38.280
Let's say the current garbage collection&nbsp;
cycle is one. Behind the scenes, sync

00:05:38.280 --> 00:05:42.960
pool keeps two internal pools&nbsp;
called locals and victim.

00:05:42.960 --> 00:05:46.080
When you put an object into the&nbsp;
pool, or when the pool creates

00:05:46.080 --> 00:05:49.760
an object using the new field,&nbsp;
it goes into the locals pool.

00:05:49.760 --> 00:05:52.640
Then some goroutines might take a few objects out,

00:05:52.640 --> 00:05:54.560
and some get left behind in the pool.

00:05:54.560 --> 00:05:57.280
When the second garbage collection cycle starts,

00:05:57.280 --> 00:06:01.440
all the remaining objects in the locals&nbsp;
pool get moved over to the victim pool.

00:06:01.440 --> 00:06:05.840
Any objects that were already checked&nbsp;
out by goroutines aren't affected.

00:06:05.840 --> 00:06:10.640
After that, both new objects and any&nbsp;
objects you put back into the pool

00:06:10.640 --> 00:06:12.960
will go into the locals pool as usual.

00:06:12.960 --> 00:06:16.320
Then the third garbage collection cycle starts.

00:06:16.320 --> 00:06:20.320
Everything in the victim pool gets&nbsp;
cleared out and garbage collected.

00:06:20.320 --> 00:06:23.920
If you don't put anything back into the pool&nbsp;
for a while, it eventually becomes empty.

00:06:24.880 --> 00:06:27.920
In other words, any object stored&nbsp;
in sync.Pool will be cleared

00:06:27.920 --> 00:06:31.920
after at least two garbage collection cycles.

00:06:31.920 --> 00:06:35.200
I say at least because some objects in the victim

00:06:35.200 --> 00:06:39.120
pool can be moved back and&nbsp;
reused when the locals are empty.

00:06:39.120 --> 00:06:40.720
But that's another story.

00:06:40.720 --> 00:06:43.040
And if this video is helpful so far,

00:06:43.040 --> 00:06:46.400
give me a sign by hitting subscribe&nbsp;
and tapping the little bell so

00:06:46.400 --> 00:06:49.520
you'll get notified the&nbsp;
moment I post something new.

00:06:49.520 --> 00:06:50.480
Alright,

00:06:50.480 --> 00:06:54.160
not all objects are, in fact, worth pooling.

00:06:54.160 --> 00:06:56.720
The best candidates are objects&nbsp;
that are used frequently

00:06:56.720 --> 00:07:00.800
and are expensive to initialize or maintain.

00:07:00.800 --> 00:07:04.880
A good example of something you'd&nbsp;
put in a sync.Pool is a bytes.Buffer.

00:07:04.880 --> 00:07:08.480
Let's understand why by looking&nbsp;
at how it works under the hood.

00:07:08.480 --> 00:07:12.000
It's basically a thin wrapper around a byte slice.

00:07:12.000 --> 00:07:15.920
And as you probably know, a byte&nbsp;
slice is just an array underneath

00:07:15.920 --> 00:07:20.160
with a fixed number of elements, like&nbsp;
one, two, three, four, or five, whatever.

00:07:20.160 --> 00:07:23.520
So imagine we grab a bytes.Buffer from the pool,

00:07:23.520 --> 00:07:26.640
and its underlying array&nbsp;
currently has two elements.

00:07:26.640 --> 00:07:28.480
These elements are stored on the heap.

00:07:28.480 --> 00:07:30.080
So when you start appending data

00:07:30.080 --> 00:07:33.200
but you fill up the first and&nbsp;
second slots, everything's fine.

00:07:33.200 --> 00:07:36.640
But once you go past the&nbsp;
slice's capacity, Go quietly

00:07:36.640 --> 00:07:40.080
allocates a bigger one on the&nbsp;
heap and copies everything over.

00:07:40.080 --> 00:07:42.000
That's how it handles growth behind the scenes,

00:07:42.000 --> 00:07:44.960
and you just keep appending&nbsp;
like nothing ever changed.

00:07:44.960 --> 00:07:48.000
And that original two-element&nbsp;
array becomes garbage,

00:07:48.000 --> 00:07:51.520
which of course adds extra&nbsp;
work for the garbage collector.

00:07:51.520 --> 00:07:55.040
If you create a new byte buffer for every request,

00:07:55.040 --> 00:07:58.640
you end up paying that growth&nbsp;
cost over and over again,

00:07:58.640 --> 00:08:02.000
while also generating a lot of garbage&nbsp;
for the garbage collector to clean up.

00:08:02.000 --> 00:08:06.800
With sync.Pool, you can recycle buffers&nbsp;
along with their already expanded slices.

00:08:06.800 --> 00:08:10.320
If a buffer has grown large in a&nbsp;
previous use, you just call the reset

00:08:10.320 --> 00:08:14.960
function to move its length back&nbsp;
to zero while keeping its capacity.

00:08:14.960 --> 00:08:17.760
And that pretty much covers all the fundamentals.

00:08:17.760 --> 00:08:21.120
Now it's time to go back to our&nbsp;
earlier question and see why

00:08:21.120 --> 00:08:24.800
we usually put pointers into the&nbsp;
pool to avoid extra allocations.

00:08:24.800 --> 00:08:29.280
To understand this, we need to look&nbsp;
deeper into how the empty interface works.

00:08:29.280 --> 00:08:35.440
Because when you put an object into a pool,&nbsp;
it gets boxed into the empty interface.

00:08:35.440 --> 00:08:40.560
In Go, an empty interface is represented&nbsp;
internally by a struct called eface.

00:08:40.560 --> 00:08:44.800
This struct has two parts: one&nbsp;
pointer to the type information

00:08:44.800 --> 00:08:47.520
and another pointer to the actual data.

00:08:47.520 --> 00:08:51.840
The key detail here is that the&nbsp;
data field is itself a pointer.

00:08:51.840 --> 00:08:56.800
When you assign a byte, an int, or&nbsp;
even a struct to an empty interface,

00:08:56.800 --> 00:09:00.800
Go copies the value and then sets&nbsp;
the data pointer to that copy.

00:09:00.800 --> 00:09:02.240
Let's take an example.

00:09:02.240 --> 00:09:07.280
Suppose we have an integer a, and we&nbsp;
assign it to an empty interface b.

00:09:07.280 --> 00:09:11.280
The values of both the integer a&nbsp;
and the empty interface b are stored

00:09:11.280 --> 00:09:12.800
on the stack.

00:09:12.800 --> 00:09:17.440
As we just discussed, assigning&nbsp;
a to b creates a copy of a,

00:09:17.440 --> 00:09:22.400
and in this case that copy also lives on&nbsp;
the stack frame of the sample function.

00:09:22.400 --> 00:09:26.720
The data field of the empty interface&nbsp;
then points to this hidden copy.

00:09:26.720 --> 00:09:30.000
That's why if we change the value of a,

00:09:30.000 --> 00:09:34.160
the value stored inside the&nbsp;
interface doesn't change.

00:09:34.160 --> 00:09:35.600
But here's the thing.

00:09:35.600 --> 00:09:38.080
The empty interface has an optimization

00:09:38.080 --> 00:09:42.400
that completely changes the&nbsp;
situation when we use pointers.

00:09:42.400 --> 00:09:46.480
Let's say we have a pointer x&nbsp;
that points to the variable a.

00:09:46.480 --> 00:09:50.240
Instead of assigning a directly&nbsp;
to the empty interface b,

00:09:50.240 --> 00:09:53.520
we assign the pointer x to b.

00:09:54.320 --> 00:09:58.400
Let's map these variables to the stack&nbsp;
frame to make it easier to visualize.

00:09:58.400 --> 00:10:00.960
In this case, there's no copy at all.

00:10:00.960 --> 00:10:04.880
The data field of the empty interface&nbsp;
points directly to the memory slot

00:10:04.880 --> 00:10:08.160
that x refers to, which is a.

00:10:08.160 --> 00:10:12.160
This is the key idea behind&nbsp;
using pointers with a sync.Pool.

00:10:12.160 --> 00:10:15.280
By adding just one level of&nbsp;
indirection through a pointer,

00:10:15.280 --> 00:10:20.000
we avoid making a copy every time we use it.

00:10:20.000 --> 00:10:22.320
Suppose we have a pool of integers.

00:10:22.320 --> 00:10:24.960
After retrieving one and setting it to 500,

00:10:24.960 --> 00:10:28.640
we return it to the pool using the Put function.

00:10:28.640 --> 00:10:31.280
Since the Put function takes an empty interface

00:10:31.280 --> 00:10:34.880
as its argument, behind the&nbsp;
scenes, passing an integer causes

00:10:34.880 --> 00:10:39.440
a copy of that integer to be&nbsp;
created somewhere in memory.

00:10:39.440 --> 00:10:43.680
The variable a and the empty&nbsp;
interface are allocated on the stack,

00:10:43.680 --> 00:10:47.040
but unfortunately, the copy&nbsp;
is allocated on the heap.

00:10:47.040 --> 00:10:49.600
The Go compiler's escape analysis is what decides

00:10:49.600 --> 00:10:52.080
which values get placed on the heap.

00:10:52.080 --> 00:10:56.160
You can confirm this by building&nbsp;
the application with the -gcflags="-m".

00:10:57.760 --> 00:11:02.080
The compiler might print something&nbsp;
like “a escapes to heap.”

00:11:02.080 --> 00:11:06.080
Escapes to heap doesn't mean the&nbsp;
variable itself is moved to the heap.

00:11:06.080 --> 00:11:09.600
It just means a heap allocation&nbsp;
is required to hold the value,

00:11:09.600 --> 00:11:14.160
while interestingly, the original variable&nbsp;
can still remain on the stack frame.

00:11:14.160 --> 00:11:16.640
If you're still confused&nbsp;
about how the compiler decides

00:11:16.640 --> 00:11:20.560
this during static analysis, think&nbsp;
about it this way: a goroutine's

00:11:20.560 --> 00:11:24.880
stack is private to that&nbsp;
goroutine, like private storage,

00:11:24.880 --> 00:11:27.760
but the copy you put into a&nbsp;
pool needs to be shared across

00:11:27.760 --> 00:11:31.440
multiple goroutines, so it has&nbsp;
to move into a shared memory

00:11:31.440 --> 00:11:35.680
region that all goroutines&nbsp;
can access, which is the heap.

00:11:35.680 --> 00:11:39.600
That's the logical answer, but&nbsp;
if you want a precise one, it's

00:11:39.600 --> 00:11:43.520
because Go's escape analysis&nbsp;
simply can't prove that your object

00:11:43.520 --> 00:11:48.480
is safe to be allocated on the stack&nbsp;
after being passed to the Put function.

00:11:48.480 --> 00:11:50.640
Let's change our example to use a pointer.

00:11:50.640 --> 00:11:54.080
When you get a pointer from the pool,&nbsp;
that pointer refers to an object

00:11:54.080 --> 00:11:56.560
that's already allocated on the heap.

00:11:56.560 --> 00:11:58.880
It doesn't matter whether&nbsp;
that object was just created

00:11:58.880 --> 00:12:01.520
or has existed for a long time.

00:12:01.520 --> 00:12:05.440
What matters is that when you call&nbsp;
put, since it's already a pointer,

00:12:05.440 --> 00:12:09.200
the resulting interface points&nbsp;
directly to the same object.

00:12:09.200 --> 00:12:13.840
No extra heap allocation is needed,&nbsp;
which is exactly what we want.

00:12:13.840 --> 00:12:16.560
Alright, that's the story of sync.Pool.

00:12:16.560 --> 00:12:19.840
Short-lived objects, pointer&nbsp;
tricks, and some clever garbage

00:12:19.840 --> 00:12:21.760
collector behaviors under the hood.

00:12:21.760 --> 00:12:24.800
If you've got an idea that&nbsp;
could be a good next topic,

00:12:24.800 --> 00:12:26.160
leave a comment.

00:12:26.160 --> 00:12:28.400
And that's all for now.

00:12:28.400 --> 00:12:30.160
See you in the next discussion.

