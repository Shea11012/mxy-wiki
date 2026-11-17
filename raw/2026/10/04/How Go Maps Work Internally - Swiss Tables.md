---
source_url: "https://www.youtube.com/watch?v=TdtNViWk45I"
ingested: 2026-10-04
extracted_to:
  - "Go map Swiss Tables 视频学习笔记"
---
# How Go Maps Work Internally: Swiss Tables

- Video: https://www.youtube.com/watch?v=TdtNViWk45I
- Channel: Phuong Le
- Uploaded: 2026-09-28
- Duration: 30:04
- Evidence: manually provided English captions; original VTT text and timestamps preserved.
- Scope: transcript only; visual-only details have not been independently inspected.

## Video description and chapters

This video is meant for learning in a simple way, so it doesn't cover every edge case, and some parts are simplified to keep them easy to follow.

Also, these internals can change from one Go version to the next. If something big changes, let me know and I'll make a follow-up.

Chapters
0:00 Go 1.24 maps are Swiss Tables
2:19 Every map has its own seed
3:04 Storage layout: directory, table, group
3:53 Inside a group: slots and control bytes
4:38 Splitting the hash into H1 and H2
7:09 Checking 8 control bytes at once (SIMD)
8:22 When a group is full: tables
10:21 Triangular probing
11:56 Load factor and the 1024-slot limit
13:26 How the directory grows
17:43 make(map, hint): what the size hint does
19:32 The old map: overflow buckets
22:06 Benchmark: Swiss Table vs the old map
27:35 Go 1.27: mapsplitgroup

## Original English subtitles (VTT)

```vtt
WEBVTT
Kind: captions
Language: en

00:00:00.650 --> 00:00:07.580
Go 1.24 replaced the old map implementation 
with a design based on Swiss Tables.

00:00:07.580 --> 00:00:11.810
We can understand that design 
by starting with a small map.

00:00:11.810 --> 00:00:15.750
Here we create a map with 
string keys and integer values.

00:00:15.750 --> 00:00:18.160
That's how we see the map as developers.

00:00:18.160 --> 00:00:22.800
But to see how it really works, we have 
to look at it from the runtime's side.

00:00:22.800 --> 00:00:26.060
The map variable we see 
here is actually a pointer,

00:00:26.060 --> 00:00:29.300
and it points to an internal 
structure called maps.Map.

00:00:29.300 --> 00:00:32.270
This structure holds a 
few details about the map

00:00:32.270 --> 00:00:36.300
and tells the runtime where to 
find the storage for its entries.

00:00:36.300 --> 00:00:38.640
We'll follow that pointer first,

00:00:38.640 --> 00:00:43.220
then look at how Go stores and 
finds the key-value pairs.

00:00:43.220 --> 00:00:46.240
Now we assign m to another variable, m2.

00:00:46.240 --> 00:00:51.200
Go doesn't create a second map here, and 
it doesn't copy any key-value pairs.

00:00:51.200 --> 00:00:53.940
It only copies the pointer to the map.

00:00:53.940 --> 00:00:57.350
So both variables point 
to the same maps.Map.

00:00:57.350 --> 00:00:59.730
When we set dog to 1 through m2

00:00:59.730 --> 00:01:03.950
and then print m with fmt.Println, 
m has the new entry too.

00:01:03.950 --> 00:01:05.960
Even though we never touched m.

00:01:05.960 --> 00:01:09.020
We can actually check 
that they share one map

00:01:09.020 --> 00:01:12.220
with a simple tool, the built-in println.

00:01:12.220 --> 00:01:15.910
Unlike fmt.Println, it 
doesn't print the entries.

00:01:15.910 --> 00:01:20.010
It prints the pointer itself, and both 
variables print the same address.

00:01:20.010 --> 00:01:23.580
So they point to the same place 
in memory and share the same map.

00:01:23.580 --> 00:01:27.630
And the same thing happens in 
something most of us do all the time,

00:01:27.630 --> 00:01:30.090
passing a map to a function.

00:01:30.090 --> 00:01:33.590
The parameter m2 gets 
a copy of the pointer,

00:01:33.590 --> 00:01:37.960
so println inside useMap 
prints the same address again.

00:01:37.960 --> 00:01:40.010
Now let's look inside maps.Map.

00:01:40.010 --> 00:01:42.780
Its first two fields 
describe the map itself,

00:01:42.780 --> 00:01:47.210
and the two below them tell Go how to 
reach the storage for its entries.

00:01:47.210 --> 00:01:48.590
We can start with used.

00:01:48.590 --> 00:01:52.500
This field counts how many entries 
the map currently contains.

00:01:52.500 --> 00:01:56.050
And Go keeps it up to date as 
entries are added or removed.

00:01:56.050 --> 00:02:00.220
So when we ask for the length of 
the map, Go already has the answer.

00:02:00.220 --> 00:02:03.520
It reads this field and 
converts the count to an int.

00:02:03.520 --> 00:02:06.640
Go does not need to visit 
the entries to count them.

00:02:06.640 --> 00:02:10.000
The field is at a known 
position in the Map structure.

00:02:10.000 --> 00:02:13.990
That's why len is still really 
fast, even for a huge map.

00:02:13.990 --> 00:02:18.840
The next field is seed. Go uses it 
to decide where to store each entry.

00:02:18.840 --> 00:02:21.000
Every map gets its own random seed,

00:02:21.000 --> 00:02:23.280
and Go picks it when the map is created,

00:02:23.280 --> 00:02:27.280
no matter whether we create it 
with make or with a map literal.

00:02:27.280 --> 00:02:30.550
A map keeps its entries in 
its storage, the box here.

00:02:30.550 --> 00:02:33.390
To decide where in that 
storage to put a key,

00:02:33.390 --> 00:02:35.930
Go turns the key into a 
number called a hash,

00:02:35.930 --> 00:02:38.990
and it doesn't use just the 
key to make that number,

00:02:38.990 --> 00:02:41.440
it uses the key together 
with the map's seed.

00:02:41.440 --> 00:02:45.230
Now let's create the same map again 
and put in the same three keys.

00:02:45.230 --> 00:02:47.060
The new map gets a new seed,

00:02:47.060 --> 00:02:50.960
so every key gets a new hash and 
is stored in a different place.

00:02:50.960 --> 00:02:55.530
And every time we create it, we get 
another seed and another layout,

00:02:55.530 --> 00:02:57.500
even though the keys never change.

00:02:57.500 --> 00:03:01.280
This box is just a simple way to 
see where each key is stored.

00:03:01.280 --> 00:03:04.760
Next, let's look at how Go 
actually lays out that storage.

00:03:04.760 --> 00:03:09.060
The storage layout changes as the 
map holds more key-value pairs,

00:03:09.060 --> 00:03:11.080
so let's look at three examples.

00:03:11.080 --> 00:03:13.160
With 4 keys, one group is enough,

00:03:13.160 --> 00:03:15.810
because a group has 
room for up to 8 pairs,

00:03:15.810 --> 00:03:19.670
and the dirPtr field in the map 
points straight to this group.

00:03:19.670 --> 00:03:23.280
With 100 keys, the map gets 
a directory with one entry,

00:03:23.280 --> 00:03:26.590
and that entry points to 
a table with 16 groups.

00:03:26.590 --> 00:03:30.210
With 1,000 keys, the 
directory has two entries,

00:03:30.210 --> 00:03:34.420
and each one points to its 
own table with 128 groups.

00:03:34.420 --> 00:03:38.660
A map doesn't switch layouts 
at exactly 100 or 1,000 keys,

00:03:38.660 --> 00:03:43.450
and we're just using these numbers 
to show what the layouts look like.

00:03:43.450 --> 00:03:47.170
So the directory points to 
tables, each table holds groups,

00:03:47.170 --> 00:03:50.980
and the groups are where Go 
keeps the key-value pairs,

00:03:50.980 --> 00:03:53.540
so let's look inside a single group.

00:03:53.540 --> 00:03:57.790
Each group has 8 slots for key-value 
pairs and 8 control bytes,

00:03:57.790 --> 00:04:01.150
and each control byte describes 
the slot right below it.

00:04:01.150 --> 00:04:05.160
The numbers in the control bytes 
come from the keys in our example,

00:04:05.160 --> 00:04:08.090
and they're not the values 
stored with those keys.

00:04:08.090 --> 00:04:13.200
We can write the layout of a group as a 
struct to see how these parts fit together.

00:04:13.200 --> 00:04:17.150
Go stores all of those control 
bytes in a single 64-bit word.

00:04:17.150 --> 00:04:20.140
The slots field holds 
the 8 key-value pairs,

00:04:20.140 --> 00:04:22.930
with one key and one element in each slot.

00:04:22.930 --> 00:04:28.320
This struct is just pseudocode, because 
Key and Elem change with each map type.

00:04:28.320 --> 00:04:31.990
For a map[string]int, the 
compiler creates a group type

00:04:31.990 --> 00:04:34.670
with string keys and int elements,

00:04:34.670 --> 00:04:38.550
and a map[int][]byte 
gets its own group type.

00:04:38.550 --> 00:04:42.560
For each key in a group, the control 
byte comes from the key's hash.

00:04:42.560 --> 00:04:45.710
Before storing cow, Go hashes 
it with the map's seed

00:04:45.710 --> 00:04:48.140
and splits the hash into 2 parts.

00:04:48.140 --> 00:04:52.090
On 64-bit machines, the 
top 57 bits become H1

00:04:52.090 --> 00:04:54.520
and the bottom 7 bits become H2.

00:04:54.520 --> 00:04:57.750
Go uses H1 to pick a 
starting group for the key,

00:04:57.750 --> 00:05:01.810
but our map has only 1 group, 
so there's nothing to pick yet.

00:05:01.810 --> 00:05:04.830
H2 goes into the control byte of the slot,

00:05:04.830 --> 00:05:07.440
and for cow we'll use 42 as an example.

00:05:07.440 --> 00:05:11.480
Dog, cat and fox got their 
control bytes the same way.

00:05:11.480 --> 00:05:14.630
A control byte has 8 bits 
and H2 needs only 7,

00:05:14.630 --> 00:05:19.370
so Go uses the top bit to mark 
whether the slot holds a live entry.

00:05:19.370 --> 00:05:22.860
When that bit is 0, the 
other 7 bits hold H2,

00:05:22.860 --> 00:05:26.160
so a live control byte is 
never higher than 127.

00:05:26.160 --> 00:05:30.120
But if the top bit is 1, the 
slot doesn't hold a live entry.

00:05:30.120 --> 00:05:35.100
It's either empty or deleted, and 
the other 7 bits tell us which one:

00:05:35.100 --> 00:05:39.640
an empty slot is 128, and 
a deleted slot is 254.

00:05:39.640 --> 00:05:42.620
A deleted slot is also called a tombstone,

00:05:42.620 --> 00:05:45.260
and neither one holds a key-value pair.

00:05:45.260 --> 00:05:50.330
The difference between empty and deleted 
only matters if our map has many groups.

00:05:50.330 --> 00:05:53.780
Because, for example, when 
a map needs to find a key,

00:05:53.780 --> 00:05:56.430
an empty slot means the 
group was never full,

00:05:56.430 --> 00:05:59.580
so no key had to move on to the 
next group, and you can stop.

00:05:59.580 --> 00:06:03.760
A tombstone doesn't tell us that, 
because the group could have been full,

00:06:03.760 --> 00:06:05.810
so the key might be in the next group.

00:06:05.810 --> 00:06:08.800
When the map grows, these 
tombstones get cleaned up,

00:06:08.800 --> 00:06:13.230
but a map never shrinks, so deleting 
keys doesn't give its memory back,

00:06:13.230 --> 00:06:16.430
and the freed slots just 
get reused by new keys.

00:06:16.430 --> 00:06:19.000
Now before storing cow, Go first checks

00:06:19.000 --> 00:06:23.335
whether this adds a new key or 
just updates an existing one.

00:06:23.335 --> 00:06:27.390
Cow's H2 is 42, and dog's 
control byte is 42 also,

00:06:27.390 --> 00:06:29.190
so dog becomes a candidate.

00:06:29.190 --> 00:06:31.750
But H2 is just a small part of the hash,

00:06:31.750 --> 00:06:35.680
then two different keys can 
easily have the same H2.

00:06:35.680 --> 00:06:38.760
Go needs to compare the 
full keys after H2,

00:06:38.760 --> 00:06:42.290
and dog obviously isn't cow, 
so cow must be a new key.

00:06:42.290 --> 00:06:46.560
Go picks the first empty slot and 
puts cow and its value there,

00:06:46.560 --> 00:06:49.510
then writes H2 into that 
slot's control byte.

00:06:49.510 --> 00:06:53.530
So H2 lets Go skip every slot 
with a different control byte,

00:06:53.530 --> 00:06:56.530
and a new key goes into 
the first empty slot.

00:06:56.530 --> 00:07:00.120
Go also adds 1 to the used 
field in the map struct.

00:07:00.120 --> 00:07:04.040
From 3 to 4 entries, everything 
still fits in 1 group,

00:07:04.040 --> 00:07:06.890
so dirPtr points straight to this group,

00:07:06.890 --> 00:07:09.210
with no table or directory in between.

00:07:09.210 --> 00:07:11.960
Dog and cow now have 
the same control byte.

00:07:11.960 --> 00:07:16.200
When we set cow to 5, Go must find 
which of these keys to update.

00:07:16.200 --> 00:07:21.880
Go checks H2 against all 8 control bytes, 
but it doesn't compare them one by one.

00:07:21.880 --> 00:07:25.320
All 8 bytes sit together 
in one 64-bit word,

00:07:25.320 --> 00:07:29.800
so Go can check the whole word 
at once with a few CPU ops.

00:07:29.800 --> 00:07:33.990
On AMD64, Go uses SIMD instructions 
to speed up this check.

00:07:33.990 --> 00:07:37.440
They compare H2 with the 8 
control bytes in parallel,

00:07:37.440 --> 00:07:41.330
and Go turns the result into a 
bitmap with one bit per slot.

00:07:41.330 --> 00:07:45.100
Other architectures use 
arithmetic and bitwise operations

00:07:45.100 --> 00:07:46.500
on the same 64-bit word.

00:07:46.500 --> 00:07:50.010
Either way, Go finds all the 
candidate slots at once,

00:07:50.010 --> 00:07:53.130
and the old map had to 
check its slots one by one.

00:07:53.130 --> 00:07:55.800
We'll come back to that 
at the end of the video.

00:07:55.800 --> 00:07:58.350
Here, dog and cow are the only candidates

00:07:58.350 --> 00:08:01.180
because their control 
bytes both store H2 42.

00:08:01.180 --> 00:08:05.440
Go can skip the other 6 slots 
without reading their complete keys.

00:08:05.440 --> 00:08:10.130
It then compares the full keys, and 
dog doesn't match but cow does,

00:08:10.130 --> 00:08:11.870
so Go puts 5 into cow's slot.

00:08:11.870 --> 00:08:14.080
From here on, we'll show only the keys,

00:08:14.080 --> 00:08:17.330
as the values don't matter 
for the rest of the video.

00:08:17.330 --> 00:08:22.070
And it's time to see how the map grows, 
probably the most interesting section.

00:08:22.070 --> 00:08:26.550
Now our small map fills all 8 slots 
and a group cannot grow beyond that.

00:08:26.550 --> 00:08:30.140
So when we add a ninth key, we 
need more room for the map.

00:08:30.140 --> 00:08:32.620
Go doubles the number 
of groups from 1 to 2

00:08:32.620 --> 00:08:35.070
and introduces a table to manage them.

00:08:35.070 --> 00:08:38.560
Then Go moves the 8 old 
entries into the 2 new groups,

00:08:38.560 --> 00:08:41.190
and it stores the new key only after that.

00:08:41.190 --> 00:08:44.270
The table struct keeps 
track of the 2 groups.

00:08:44.270 --> 00:08:46.540
Its used field counts the live entries,

00:08:46.540 --> 00:08:50.040
its capacity field counts all 
the slots, full and empty,

00:08:50.040 --> 00:08:53.250
and its groups field 
points to the group array.

00:08:53.250 --> 00:08:56.370
Before this, dirPtr pointed 
straight to the group.

00:08:56.370 --> 00:08:59.460
Now it points to an array 
called the directory,

00:08:59.460 --> 00:09:03.130
and each entry in the 
directory points to a table.

00:09:03.130 --> 00:09:07.740
Right now the directory has just 1 
entry, so dirLen goes from 0 to 1.

00:09:07.740 --> 00:09:11.520
But how does Go choose which 
group each key goes into?

00:09:11.520 --> 00:09:16.160
To answer that, we need the 
other part of the hash, H1.

00:09:16.160 --> 00:09:18.020
Go uses H1 to pick a group.

00:09:18.020 --> 00:09:22.720
It divides H1 by the number of groups, 
and the remainder is the group index.

00:09:22.720 --> 00:09:26.390
We call it the starting group, 
and it will be explained soon.

00:09:26.390 --> 00:09:29.420
With 2 groups, the 
remainder is always 0 or 1,

00:09:29.420 --> 00:09:31.810
and that's just the last bit of H1.

00:09:31.810 --> 00:09:35.110
Cow's last bit is 0, so 
cow starts in group 0.

00:09:35.110 --> 00:09:38.100
Dog's last bit is 1, so 
dog starts in group 1.

00:09:38.100 --> 00:09:40.990
Now let's take 5 keys and see what happens

00:09:40.990 --> 00:09:43.970
when the table grows from 2 groups to 4.

00:09:43.970 --> 00:09:47.550
With 4 groups, Go uses 1 
more bit from the end of H1.

00:09:47.550 --> 00:09:50.930
If that new bit is 0, the key 
keeps the same starting group.

00:09:50.930 --> 00:09:53.590
If it's 1, the group index goes up by 2,

00:09:53.590 --> 00:09:56.980
so group 0 becomes group 2 
and group 1 becomes group 3.

00:09:56.980 --> 00:10:01.430
Keep in mind, the table gets a whole 
new set of groups when it grows.

00:10:01.430 --> 00:10:04.920
So every key moves, even the 
ones with the same group index.

00:10:04.920 --> 00:10:08.330
But the starting group doesn't 
always have an empty slot.

00:10:08.330 --> 00:10:11.910
Here's another map with 10 
keys in a table of 16 slots.

00:10:11.910 --> 00:10:14.690
Group 0 is full, and 
group 1 still has room.

00:10:14.690 --> 00:10:17.180
A new key can still start at group 0,

00:10:17.180 --> 00:10:21.580
so Go has to look for an empty slot 
in another group of this table.

00:10:21.580 --> 00:10:24.880
Let's add hen to see how 
Go finds an empty slot.

00:10:24.880 --> 00:10:27.350
H1 gives us group 0, but it's full,

00:10:27.350 --> 00:10:30.410
so Go checks group 1 next 
and stores hen there.

00:10:30.410 --> 00:10:32.460
It checks groups in a fixed order,

00:10:32.460 --> 00:10:35.090
and we call it the 
triangular probe sequence.

00:10:35.090 --> 00:10:37.580
Let's look at a table with 8 groups,

00:10:37.580 --> 00:10:40.580
and say H1 gives us 
group 3, but it's full.

00:10:40.580 --> 00:10:42.880
Go moves forward 1 group to group 4.

00:10:42.880 --> 00:10:45.940
If that's full too, it 
moves 2 groups to group 6,

00:10:45.940 --> 00:10:48.710
then 3 groups, and that goes past group 7,

00:10:48.710 --> 00:10:50.540
so it wraps around to group 1.

00:10:50.540 --> 00:10:52.810
If we keep going, each step grows by 1,

00:10:52.810 --> 00:10:57.310
so the steps stack up like a triangle, 
and that's where the name comes from.

00:10:57.310 --> 00:11:01.060
Go reaches all 8 groups without 
hitting the same one twice.

00:11:01.060 --> 00:11:04.940
This works because the number of 
groups is always a power of 2.

00:11:04.940 --> 00:11:08.210
So why this order? Let's 
compare it with linear probing.

00:11:08.210 --> 00:11:11.390
Linear probing just checks 
the next group each time.

00:11:11.390 --> 00:11:14.000
Here are two tables with the same keys,

00:11:14.000 --> 00:11:16.840
and group 0 is already full in both.

00:11:16.840 --> 00:11:19.320
Now 8 more keys start at group 0.

00:11:19.320 --> 00:11:23.640
Linear probing puts them in groups 
1 to 4, one right after another,

00:11:23.640 --> 00:11:27.410
so groups 0 to 3 are all full 
and sit next to each other.

00:11:27.410 --> 00:11:29.000
We call this a cluster.

00:11:29.000 --> 00:11:34.180
Triangular probing puts the same 8 
keys in groups 1, 3 and 6 instead,

00:11:34.180 --> 00:11:35.810
so they are spread out.

00:11:35.810 --> 00:11:40.020
Let's add the same new key to 
both tables, starting at group 1.

00:11:40.020 --> 00:11:43.880
Linear probing checks groups 1, 
2 and 3, and they're all full,

00:11:43.880 --> 00:11:46.620
so it checks group 4 and 
stores the key there.

00:11:46.620 --> 00:11:49.430
Triangular probing checks 
group 1, then group 2,

00:11:49.430 --> 00:11:53.790
and group 2 still has room, because 
none of those 8 keys went there.

00:11:53.790 --> 00:11:56.140
So it checks only 2 groups instead of 4.

00:11:56.140 --> 00:11:59.330
The table grows by doubling 
the number of groups

00:11:59.330 --> 00:12:04.800
until it reaches 128 groups, 
which means it has 1,024 slots.

00:12:04.800 --> 00:12:08.740
But Go doesn't wait until every 
slot is full before it grows.

00:12:08.740 --> 00:12:11.870
Instead, it uses something 
called the load factor.

00:12:11.870 --> 00:12:15.090
Go adds the number of live 
entries and deleted slots,

00:12:15.090 --> 00:12:18.190
then divides that by the 
total number of slots.

00:12:18.190 --> 00:12:23.570
For example, our table has 16 slots, 
10 live entries, and no deleted slots.

00:12:23.570 --> 00:12:27.120
That gives us a load factor of 62.5%.

00:12:27.120 --> 00:12:32.190
Go lets the table go up to 
7 out of 8 slots, or 87.5%.

00:12:32.190 --> 00:12:36.570
Without this limit, a nearly full 
table would make inserts and lookups

00:12:36.570 --> 00:12:39.870
check many groups before 
they find an empty slot.

00:12:39.870 --> 00:12:44.650
Our table has 6 empty slots, but 
new keys can use only 4 of them,

00:12:44.650 --> 00:12:49.760
because 7 out of 8 of 16 slots 
is 14, and 10 are already used.

00:12:49.760 --> 00:12:53.330
After 4 more insertions, the 
table reaches the limit,

00:12:53.330 --> 00:12:55.690
even though 2 slots are still empty.

00:12:55.690 --> 00:12:59.590
If one more new key needs an 
empty slot, Go grows the table.

00:12:59.590 --> 00:13:03.510
It doubles the group array, 
hashes every live key again,

00:13:03.510 --> 00:13:05.780
moves the entries into the new groups.

00:13:05.780 --> 00:13:09.200
When that work is done, then 
Go tries the insertion again.

00:13:09.200 --> 00:13:10.310
At 1,024 slots,

00:13:10.310 --> 00:13:13.710
the table can't double its 
number of groups anymore,

00:13:13.710 --> 00:13:17.490
but we still need more space 
to store more key-value pairs.

00:13:17.490 --> 00:13:20.810
When that happens, Go 
splits the table in two,

00:13:20.810 --> 00:13:23.760
and both new tables have 1,024 slots.

00:13:23.760 --> 00:13:25.890
This is a split, not a doubling.

00:13:25.890 --> 00:13:29.850
After the split, Go has to know 
which table a key goes into,

00:13:29.850 --> 00:13:32.540
and then which group inside that table.

00:13:32.540 --> 00:13:35.720
We already used the low 
bits of H1 to pick a group.

00:13:35.720 --> 00:13:39.880
To pick a table, Go uses the 
other end, the highest bit of H1,

00:13:39.880 --> 00:13:42.520
so the two choices use different bits.

00:13:42.520 --> 00:13:46.510
For cow, that bit is 0, so 
it picks directory entry 0,

00:13:46.510 --> 00:13:48.900
and that entry points to table 0.

00:13:48.900 --> 00:13:53.290
Then the low bits pick a 
starting group inside table 0,

00:13:53.290 --> 00:13:58.750
and with 128 groups, that takes 7 
bits, because 2 to the 7 is 128.

00:13:58.750 --> 00:14:02.960
Dog goes through the same steps, 
but its highest bit is 1,

00:14:02.960 --> 00:14:06.870
so it goes to table 1, and its 
low bits pick group 1 there.

00:14:06.870 --> 00:14:11.740
But what if table 1, the table on the 
right here, needs more space again?

00:14:11.740 --> 00:14:15.220
Let's say cow and cat have 
0 as the highest bit of H1,

00:14:15.220 --> 00:14:19.610
so they're in table 0, and dog and 
bat have 1, so they're in table 1.

00:14:19.610 --> 00:14:23.780
When table 1 splits into 2 new 
tables, Go has to move its entries,

00:14:23.780 --> 00:14:27.480
so it reads one more high bit 
to decide where each one goes.

00:14:27.480 --> 00:14:31.090
Those top bits are the key's 
index into the directory.

00:14:31.090 --> 00:14:35.670
Dog's next bit is 0, so its top 2 
bits are 1 0, which is index 2,

00:14:35.670 --> 00:14:37.450
and dog moves to table 2.

00:14:37.450 --> 00:14:41.830
Bat's next bit is 1, so it gets 
index 3 and moves to table 3.

00:14:41.830 --> 00:14:44.650
Now Go reads 2 bits for every key.

00:14:44.650 --> 00:14:48.450
Cat's top 2 bits are 0 
0, so it gets index 0,

00:14:48.450 --> 00:14:52.680
and cow's top 2 bits are 
0 1, so it gets index 1.

00:14:52.680 --> 00:14:58.440
Now we have 4 indexes but only 3 
tables. Does table 0 have to split too?

00:14:58.440 --> 00:15:01.690
It doesn't, and this is 
what the directory is for.

00:15:01.690 --> 00:15:04.060
Each index picks a directory entry,

00:15:04.060 --> 00:15:07.050
and that entry tells 
Go which table to use.

00:15:07.050 --> 00:15:09.660
The directory now has 4 entries,

00:15:09.660 --> 00:15:12.730
and entries 0 and 1 both point to table 0.

00:15:12.730 --> 00:15:15.100
So a split only touches one table,

00:15:15.100 --> 00:15:18.600
and Go moves just the 
entries in that table.

00:15:18.600 --> 00:15:21.350
A table has at most 1,024 slots,

00:15:21.350 --> 00:15:26.540
so however big the map gets, one 
split never moves more than that.

00:15:26.540 --> 00:15:29.780
Table 0 has 2 directory 
entries pointing to it,

00:15:29.780 --> 00:15:33.130
so Go can split it without 
growing the directory.

00:15:33.130 --> 00:15:35.120
But table 2 has only 1 entry,

00:15:35.120 --> 00:15:38.420
so before it can split, 
the directory has to grow.

00:15:38.420 --> 00:15:41.830
To know when that happens, 
Go keeps two numbers:

00:15:41.830 --> 00:15:43.670
global depth and local depth.

00:15:43.670 --> 00:15:46.730
Global depth tells Go how 
many high bits to use

00:15:46.730 --> 00:15:48.540
when it picks a directory entry.

00:15:48.540 --> 00:15:50.860
Local depth, on each table, tells Go

00:15:50.860 --> 00:15:54.250
how many of those bits it 
needs to identify that table.

00:15:54.250 --> 00:15:57.970
Here, global depth is 2 because 
the directory uses 2 bits.

00:15:57.970 --> 00:16:00.520
Two directory entries share table 0,

00:16:00.520 --> 00:16:04.720
so 1 bit is enough to identify 
it, and its local depth is 1.

00:16:04.720 --> 00:16:08.810
The other tables need both bits, 
so their local depth is 2.

00:16:08.810 --> 00:16:12.400
When table 0 splits, Go 
sees that its local depth

00:16:12.400 --> 00:16:14.250
is smaller than global depth.

00:16:14.250 --> 00:16:18.850
That means the directory already has 
enough entries for both new tables.

00:16:18.850 --> 00:16:21.970
Go uses the next high bit 
to separate the keys,

00:16:21.970 --> 00:16:24.550
then points one entry to each new table,

00:16:24.550 --> 00:16:26.670
and the directory doesn't grow.

00:16:26.670 --> 00:16:31.040
Both new tables now have a local 
depth of 2, the same as global depth.

00:16:31.040 --> 00:16:33.590
Let's go back to the original 3-table map,

00:16:33.590 --> 00:16:35.850
and this time we split table 2 instead.

00:16:35.850 --> 00:16:38.660
Its local depth already 
equals global depth,

00:16:38.660 --> 00:16:42.910
so Go first doubles the directory 
and raises global depth to 3.

00:16:42.910 --> 00:16:45.320
Now table 2 has 2 directory entries,

00:16:45.320 --> 00:16:48.600
so Go splits it and points 
one entry to each new table.

00:16:48.600 --> 00:16:51.480
Let's follow one lookup 
through this new map.

00:16:51.480 --> 00:16:55.280
The high bits of H1 pick the 
directory entry and its table,

00:16:55.280 --> 00:16:58.230
and the low 7 bits pick 
the starting group.

00:16:58.230 --> 00:17:00.790
The H1 bits in between aren't used here.

00:17:00.790 --> 00:17:03.320
Then H2 checks the 8 slots of that group

00:17:03.320 --> 00:17:07.180
and narrows the search down to the 
keys that Go needs to compare.

00:17:07.180 --> 00:17:10.610
So why keep 2 depths instead 
of just a simple counter

00:17:10.610 --> 00:17:13.730
for how many directory 
entries point to each table?

00:17:13.730 --> 00:17:17.040
These 2 depths actually answer 
more than one question.

00:17:17.040 --> 00:17:20.440
How many hash bits does the 
directory use? Global depth.

00:17:20.440 --> 00:17:23.590
How many directory entries 
point to this table?

00:17:23.590 --> 00:17:27.170
You can see the formula here, 
and for table 0, that's 4.

00:17:27.170 --> 00:17:29.930
Which bit should Go use 
when a table splits?

00:17:29.930 --> 00:17:34.390
The next bit after its current local 
depth, so for table 0, the second bit.

00:17:34.390 --> 00:17:37.800
Does the directory need to 
grow before that split?

00:17:37.800 --> 00:17:43.180
Only when local depth equals global 
depth, like table 5, where both are 3.

00:17:43.180 --> 00:17:47.170
One more practical thing and I think 
it should be clear for you now.

00:17:47.170 --> 00:17:50.490
When we put a size on a 
make call to create a map,

00:17:50.490 --> 00:17:54.750
like a map[string]int with 100, 
Go calls that number a hint.

00:17:54.750 --> 00:17:57.340
That's because unlike an array or a slice,

00:17:57.340 --> 00:18:00.020
Go doesn't reserve 
exactly that many slots.

00:18:00.020 --> 00:18:02.760
It builds the same pieces we just saw:

00:18:02.760 --> 00:18:07.770
groups of 8 slots, tables of up 
to 1,024 slots and a directory.

00:18:07.770 --> 00:18:09.960
With a hint of 8 or less such as 4,

00:18:09.960 --> 00:18:12.910
the map starts as a 
single group of 8 slots

00:18:12.910 --> 00:18:16.140
and this is the same as 
a map without a hint.

00:18:16.140 --> 00:18:19.230
With 100 entries, the map needs a table.

00:18:19.230 --> 00:18:22.220
100 keys need 13 groups, 
which is 104 slots.

00:18:22.220 --> 00:18:24.940
But the load factor applies to a table

00:18:24.940 --> 00:18:28.890
and Go only lets a table fill up to 87.5%,

00:18:28.890 --> 00:18:32.190
which is 91 slots here for 13 groups.

00:18:32.190 --> 00:18:38.380
So to hold 100 keys without growing, 
Go needs 15 groups which is 120 slots.

00:18:38.380 --> 00:18:41.930
Then the number of groups 
has to be a power of 2,

00:18:41.930 --> 00:18:46.430
so 15 rounds up to 16 
groups, which is 128 slots.

00:18:46.430 --> 00:18:50.390
Finally with 1,000, the 
keys need 125 groups,

00:18:50.390 --> 00:18:53.800
which is less than one 
table of 128 groups.

00:18:53.800 --> 00:18:56.300
But because of the load factor,

00:18:56.300 --> 00:19:00.340
one table can only hold 896 
keys so it's not enough.

00:19:00.340 --> 00:19:04.770
So Go makes a directory with 
2 entries and 2 tables,

00:19:04.770 --> 00:19:07.700
and each table has 128 groups.

00:19:07.700 --> 00:19:10.590
That's 2,048 slots for 1,000 keys.

00:19:10.590 --> 00:19:12.930
So the hint isn't an exact size.

00:19:12.930 --> 00:19:15.860
It just means the map 
usually doesn't need to grow

00:19:15.860 --> 00:19:17.870
before it reaches that many keys,

00:19:17.870 --> 00:19:20.780
and the real number of 
slots can be much bigger.

00:19:20.780 --> 00:19:23.260
And that's all for the data structure.

00:19:23.260 --> 00:19:26.930
We've gone through the hash, 
H1, H2, the directory,

00:19:26.930 --> 00:19:31.880
the table, the group, the 
control byte, the slot, etc.

00:19:31.880 --> 00:19:35.110
It's time to go back to the 
previous map implementation

00:19:35.110 --> 00:19:38.580
and see where the extra 
work and memory came from.

00:19:38.580 --> 00:19:43.240
Before Go 1.24, a map kept its 
regular buckets in one main array.

00:19:43.240 --> 00:19:45.470
Here are 4 buckets side by side,

00:19:45.470 --> 00:19:48.460
with room for 8 key-value 
pairs in each one.

00:19:48.460 --> 00:19:50.900
The hash picks a bucket for each key,

00:19:50.900 --> 00:19:54.250
but sometimes too many keys 
go into the same bucket.

00:19:54.250 --> 00:19:58.300
When that bucket has no room 
left, Go adds an overflow bucket

00:19:58.300 --> 00:20:01.380
and puts its address in 
the overflow pointer.

00:20:01.380 --> 00:20:07.310
If the overflow bucket fills up too, Go adds 
another one, and the chain gets longer.

00:20:07.310 --> 00:20:10.840
To find a key, Go checks the 
8 slots of bucket 1 first.

00:20:10.840 --> 00:20:14.360
If the key isn't there, Go 
loads the overflow pointer

00:20:14.360 --> 00:20:16.250
and checks the next 8 slots.

00:20:16.250 --> 00:20:19.710
So every overflow bucket 
adds one more pointer load,

00:20:19.710 --> 00:20:23.570
and Go can't check the next 
bucket until that load is done.

00:20:23.570 --> 00:20:28.050
The overflow buckets can also sit 
far away from bucket 1 in memory,

00:20:28.050 --> 00:20:31.870
so the CPU often has to wait 
for main memory to read them.

00:20:31.870 --> 00:20:34.040
This way, bucket 1 gets 8 more slots

00:20:34.040 --> 00:20:36.160
without growing the main bucket array.

00:20:36.160 --> 00:20:40.060
But each overflow bucket takes extra 
memory outside the main array,

00:20:40.060 --> 00:20:44.410
and every lookup that reaches it 
has to load one more pointer first.

00:20:44.410 --> 00:20:48.630
Go 1.24 doesn't use overflow buckets 
or overflow pointers at all.

00:20:48.630 --> 00:20:51.530
Each table keeps all of 
its groups in one array.

00:20:51.530 --> 00:20:56.310
If a group is full, Go uses the 
triangular probe sequence we saw earlier

00:20:56.310 --> 00:20:59.570
to pick the next one, here 
group 1 and then group 3.

00:20:59.570 --> 00:21:03.950
It finds each group by its index in 
that array, without loading a pointer.

00:21:03.950 --> 00:21:06.640
The old map had another 
problem when it grew.

00:21:06.640 --> 00:21:10.130
Go made a new main bucket array 
with twice as many buckets,

00:21:10.130 --> 00:21:12.960
but it didn't move all 
the keys right away.

00:21:12.960 --> 00:21:17.170
When your code added a key, 
updated a value, or deleted a key,

00:21:17.170 --> 00:21:19.720
Go moved a few old buckets at a time.

00:21:19.720 --> 00:21:24.520
This avoided doing all the work of 
growing the entire map in one operation.

00:21:24.520 --> 00:21:29.260
Until every old bucket had moved, some 
keys were still in the old array,

00:21:29.260 --> 00:21:32.010
so Go had to keep both arrays in memory.

00:21:32.010 --> 00:21:35.670
Here the old array has 4 
buckets and the new one has 8,

00:21:35.670 --> 00:21:39.590
so together they take 12 
buckets, 3 times the old array,

00:21:39.590 --> 00:21:41.870
even before counting overflow buckets.

00:21:41.870 --> 00:21:45.640
A lookup also had to check 
whether its old bucket had moved,

00:21:45.640 --> 00:21:47.910
to know which array to search.

00:21:47.910 --> 00:21:54.340
Go 1.24 grows only one table at a time, 
and each table has at most 1,024 slots.

00:21:54.340 --> 00:21:58.470
So when table 1 runs out of 
room, Go splits it into 2 tables

00:21:58.470 --> 00:22:02.190
and moves its keys all at 
once, during that same insert.

00:22:02.190 --> 00:22:06.430
The other tables, and the keys in 
them, still stay where they are.

00:22:06.430 --> 00:22:10.960
To see how much faster the new map 
really is, let's run a benchmark.

00:22:10.960 --> 00:22:14.460
The easy way would be to 
compare Go 1.23 with Go 1.24,

00:22:14.460 --> 00:22:18.340
but a lot more than the map 
changed between those 2 versions,

00:22:18.340 --> 00:22:21.280
and all of it would be 
mixed into the result.

00:22:21.280 --> 00:22:25.380
Go 1.25 still has the old map 
behind an experiment flag,

00:22:25.380 --> 00:22:27.070
and Go 1.26 removed it.

00:22:27.070 --> 00:22:31.850
So we'll use Go 1.25.7 and build 
the same benchmark two times:

00:22:31.850 --> 00:22:34.630
once normally which uses the Swiss Table,

00:22:34.630 --> 00:22:39.730
and once with that experiment flag 
on, which brings back the old map.

00:22:39.730 --> 00:22:44.900
We run it on an Apple M4 Pro, which 
is ARM64, and that matters here.

00:22:44.900 --> 00:22:49.530
Because on AMD64, Go compares 
the 8 control bytes of a group

00:22:49.530 --> 00:22:51.730
with one SIMD instruction.

00:22:51.730 --> 00:22:54.120
On ARM64 it doesn't use SIMD.

00:22:54.120 --> 00:23:00.410
It reads the 8 bytes as one 64-bit integer 
and finds the matches with bit tricks.

00:23:00.410 --> 00:23:03.190
So none of these numbers include SIMD,

00:23:03.190 --> 00:23:08.170
and the Swiss Table could do better on 
an AMD64 machine than it does here.

00:23:08.170 --> 00:23:13.430
And before we trust any result, let's 
run the Swiss Table against itself.

00:23:13.430 --> 00:23:16.190
Those 2 runs never differ by more than 3%,

00:23:16.190 --> 00:23:20.210
so a difference only counts 
when it's bigger than that.

00:23:20.210 --> 00:23:26.060
In a map with about 65,000 keys, 
most lookups are 58 to 68% faster.

00:23:26.060 --> 00:23:30.640
The Go team reported something 
similar: in their microbenchmarks,

00:23:30.640 --> 00:23:34.800
some map operations ran up to 
60% faster than in Go 1.23.

00:23:34.800 --> 00:23:38.180
In their full application 
benchmarks, though,

00:23:38.180 --> 00:23:40.790
they used only about 1.5% less CPU time.

00:23:40.790 --> 00:23:43.600
To see why, let's look 
at where the time goes

00:23:43.600 --> 00:23:46.540
when Go looks up a key 
that isn't in the map.

00:23:46.540 --> 00:23:50.450
The old map checks the 8 
slots of a bucket one by one,

00:23:50.450 --> 00:23:53.360
and that loop alone takes 58% of the time.

00:23:53.360 --> 00:23:56.880
The Swiss Table checks all 
8 control bytes at once,

00:23:56.880 --> 00:24:00.550
and it only looks at the 
slots with a matching H2.

00:24:00.550 --> 00:24:05.150
But with 1 or 4 million keys, the 
map is too big for the CPU cache.

00:24:05.150 --> 00:24:09.270
Both maps spend most of their 
time waiting for memory,

00:24:09.270 --> 00:24:11.870
so lookups are only 15 to 45% faster.

00:24:11.870 --> 00:24:16.270
The new design also helps in a 
different way, while the map grows.

00:24:16.270 --> 00:24:18.430
As the map grows to 4 million keys,

00:24:18.430 --> 00:24:22.120
the old map keeps both bucket 
arrays during each grow,

00:24:22.120 --> 00:24:25.340
so its peak memory is about 
1.6 times its final size.

00:24:25.340 --> 00:24:28.050
The Swiss Table splits one table at a time

00:24:28.050 --> 00:24:30.840
so its peak is the same as its final size.

00:24:30.840 --> 00:24:34.550
The slowest single insert 
goes from 8 milliseconds

00:24:34.550 --> 00:24:36.690
down to about 0.1 milliseconds.

00:24:36.690 --> 00:24:41.440
The old map allocates its whole 
new bucket array in one insert,

00:24:41.440 --> 00:24:45.810
but the Swiss Table only splits 
one table of up to 1,024 slots.

00:24:45.810 --> 00:24:49.180
If we sort all the inserts 
from slowest to fastest,

00:24:49.180 --> 00:24:53.600
the old map has a few very long 
pauses, like the 8 millisecond one

00:24:53.600 --> 00:24:55.540
and then its times drop fast.

00:24:55.540 --> 00:25:00.400
The Swiss Table never has a pause that 
long, but it has many medium ones.

00:25:00.400 --> 00:25:02.850
Those medium pauses are table splits.

00:25:02.850 --> 00:25:06.820
They happen often and each 
one takes a few microseconds,

00:25:06.820 --> 00:25:11.080
so the slowest 1 in 10,000 
inserts get about 5 times slower.

00:25:11.080 --> 00:25:13.780
So the Swiss Table isn't always faster.

00:25:13.780 --> 00:25:16.990
The Go team also said some 
edge cases got slower,

00:25:16.990 --> 00:25:19.370
and our benchmark finds a few of them.

00:25:19.370 --> 00:25:22.990
When we keep deleting one 
key and inserting another

00:25:22.990 --> 00:25:27.140
in a map with 1 million 
keys, it's 17 to 38% slower.

00:25:27.140 --> 00:25:33.210
That's because of tombstones. A delete leaves 
a tombstone when its group has no empty slot.

00:25:33.210 --> 00:25:36.350
After that run there are 
about 219,000 of them.

00:25:36.350 --> 00:25:38.240
They're spread over many tables,

00:25:38.240 --> 00:25:40.270
but they also fill up many groups,

00:25:40.270 --> 00:25:43.820
so an insert often has to 
probe one or two more groups

00:25:43.820 --> 00:25:45.490
before it finds an empty slot.

00:25:45.490 --> 00:25:48.770
That's because an insert first 
has to look for the key.

00:25:48.770 --> 00:25:52.270
If it's already in the map, this 
is an update, not a new key.

00:25:52.270 --> 00:25:54.600
So it can't stop at the first tombstone.

00:25:54.600 --> 00:25:58.210
It only knows the key is new 
when it reaches an empty slot,

00:25:58.210 --> 00:26:01.010
and then it puts the key 
into the first tombstone.

00:26:01.010 --> 00:26:05.400
And when we create a map with room 
for 1 million keys and then fill it,

00:26:05.400 --> 00:26:08.080
the Swiss Table is 30 to 52% slower.

00:26:08.080 --> 00:26:10.960
This one comes from waiting for memory.

00:26:10.960 --> 00:26:16.230
When the CPU needs data from main 
memory, it waits about 100 nanoseconds.

00:26:16.230 --> 00:26:20.820
While it waits, a modern CPU can 
keep running the code after that,

00:26:20.820 --> 00:26:23.870
but only if it knows 
which way the code goes.

00:26:23.870 --> 00:26:26.780
When there's an if, the 
CPU guesses the answer,

00:26:26.780 --> 00:26:30.140
usually from what that if did 
last time and keeps going.

00:26:30.140 --> 00:26:33.810
When the data comes back it 
checks if the guess was right.

00:26:33.810 --> 00:26:35.810
This is called branch prediction.

00:26:35.810 --> 00:26:39.100
The old map checks the 8 
slots of a bucket one by one

00:26:39.100 --> 00:26:41.660
with if statements to find a free one.

00:26:41.660 --> 00:26:44.150
So the CPU can easily guess the slot,

00:26:44.150 --> 00:26:47.640
move on to the next insert 
before the data comes back,

00:26:47.640 --> 00:26:50.860
and several inserts wait for 
memory at the same time.

00:26:50.860 --> 00:26:55.120
The Swiss Table computes the slot 
number from the 8 control bytes

00:26:55.120 --> 00:26:58.330
with bit tricks, so there's 
no branch to predict.

00:26:58.330 --> 00:27:01.120
The CPU has to wait for those bytes first,

00:27:01.120 --> 00:27:05.190
and on this M4 Pro, the inserts 
wait for memory one after another.

00:27:05.190 --> 00:27:08.050
If we change only that 
step in the Go runtime

00:27:08.050 --> 00:27:11.510
to check the 8 control bytes 
one by one like the old map,

00:27:11.510 --> 00:27:14.500
filling a map with room for 1 million keys

00:27:14.500 --> 00:27:16.810
becomes as fast as in the old map.

00:27:16.810 --> 00:27:21.460
Of course we optimize a case, it may 
be not good for other common cases.

00:27:21.460 --> 00:27:26.130
But our benchmark only measures the map 
by itself, I mean a microbenchmark,

00:27:26.130 --> 00:27:30.040
so it can not tell us how much 
a real application saves.

00:27:30.040 --> 00:27:35.540
For a real app, the Go team's 1.5% 
is much closer to what we'd see.

00:27:35.540 --> 00:27:40.710
OK finally, Go 1.27 has an experimental 
layout called mapsplitgroup.

00:27:40.710 --> 00:27:44.780
It changes where a group keeps 
its keys and values in memory.

00:27:44.780 --> 00:27:48.330
The original layout stores 
each key beside its value.

00:27:48.330 --> 00:27:53.380
The split layout puts all 8 keys in 
one array and all 8 values in another.

00:27:53.380 --> 00:27:56.450
A key and its value still 
share the same index.

00:27:56.450 --> 00:28:00.540
During a lookup for bat, the 
control bytes identify a candidate.

00:28:00.540 --> 00:28:03.320
Go reads that key and 
compares it with bat.

00:28:03.320 --> 00:28:06.310
The value is not needed 
until the key matches.

00:28:06.310 --> 00:28:09.310
In the original layout, 
values sit between the keys.

00:28:09.310 --> 00:28:12.310
The split layout puts the 
keys next to each other

00:28:12.310 --> 00:28:15.770
so the keys Go checks during 
a lookup sit side by side.

00:28:15.770 --> 00:28:19.090
But Go added this layout 
mainly to save memory.

00:28:19.090 --> 00:28:21.800
Let's use a map[int64]struct{} this time.

00:28:21.800 --> 00:28:24.410
We only care whether a key is in the map,

00:28:24.410 --> 00:28:28.290
so each value is an empty 
struct and takes zero bytes.

00:28:28.290 --> 00:28:30.260
Each int64 key takes 8 bytes.

00:28:30.260 --> 00:28:33.670
In the original layout, each 
slot is a small struct:

00:28:33.670 --> 00:28:37.390
the key comes first, and the 
empty value is its last field.

00:28:37.390 --> 00:28:39.270
The empty value takes no space

00:28:39.270 --> 00:28:43.050
so a pointer to it would point 
past the end of the slot.

00:28:43.050 --> 00:28:48.610
To keep that pointer inside the slot, Go 
adds 1 byte of padding after the value.

00:28:48.610 --> 00:28:53.570
Then on a 64-bit machine, Go rounds 
the slot up to a multiple of 8 bytes

00:28:53.570 --> 00:28:56.930
so the next key still 
starts at a multiple of 8.

00:28:56.930 --> 00:29:01.330
One slot now takes 16 bytes: 8 
for the key and 8 for padding.

00:29:01.330 --> 00:29:05.540
We have 8 slots so that same 
padding is added 8 times.

00:29:05.540 --> 00:29:10.210
The keys and padding take 128 bytes 
plus the 8-byte control word.

00:29:10.210 --> 00:29:13.740
Now put the keys together 
with the split layout.

00:29:13.740 --> 00:29:18.450
All 8 keys fit in 64 bytes and the 
empty value array takes zero bytes.

00:29:18.450 --> 00:29:22.120
This time the empty value array is 
the last field of the whole group

00:29:22.120 --> 00:29:25.390
so Go adds padding only once 
at the end of the group.

00:29:25.390 --> 00:29:28.120
That's 8 bytes of padding instead of 64.

00:29:28.120 --> 00:29:30.380
The control word doesn't change.

00:29:30.380 --> 00:29:34.400
For the keys and padding, 
we go from 128 bytes to 72

00:29:34.400 --> 00:29:37.350
and saving 56 bytes in this one group.

00:29:37.350 --> 00:29:40.590
The keys haven't changed; 
we've just removed the padding

00:29:40.590 --> 00:29:42.800
that was repeated in every slot.

00:29:42.800 --> 00:29:46.240
In Go 1.27, mapsplitgroup 
is still an experiment

00:29:46.240 --> 00:29:49.280
and you enable it with 
the GOEXPERIMENT setting.

00:29:49.280 --> 00:29:53.660
On the dev branch now, the Go team 
has already made it the default

00:29:53.660 --> 00:29:57.510
so I expect it would be 
the default from Go 1.28.

00:29:57.510 --> 00:30:00.790
And that's all. See you in the next video.


```
