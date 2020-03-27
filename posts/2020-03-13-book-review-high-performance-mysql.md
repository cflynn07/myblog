[My notes from reading High Performance MySQL][notes]

<img src="/static/images/high-performance-mysql.png" style="width:50%;">

I just finished reading "High Performance MySQL" by Baron Schwartz, Peter
Zaitsev and Vadim Tkachenko. The printed version is 790 pages. I read this book
daily, sequentially... and it took me over a month to finish it. Absorbing the
information was a little bit like drinking from a firehose. In this post I want
to focus the content but also my thoughts on digesting dense technical books.

I'm working my way through a pile of technical books as part of a new years
resolution to "level up" in software engineering. Simultaneously, I'm working
on enhancing on my concentration abilities. While reading HPM this past month,
I progressively got better at catching myself "changing the channel" in my
brain. Whenever some thought or idea would enter my brain while I was working
on downloading the information in this book I would actively fight back and
restore my focus. Over the month I got a lot better at this. I subjectively
feel like my mind tries less often to "change the channel" and when it does I
find it easier to overcome the urge. My learning pace quickened and I'm eager
to see how I perform when I dive into the next book in my queue.

<img src="/static/images/high-performance-mysql-time-blocks.png" style="width:50%;">
###### start and end time for 10 page "blocks"

Quantifying my progress by dividing the book into increments of 10 pages
helped. I noted the start and end time for every block. Initially, my average
time spent per block of 10 pages was approximately 1 hour. After reading 1/2
the book I noticed a 25-50% increase in how quickly I could read and comprehend
10 page blocks. I attribute this to fewer incidents of "changing the channel"
in my mind and thus being able to focus on the text for a greater percentage of
my time. To read the entire book and comprehend it well probably took on the
order of 70-80 hours.

I took a pretty comprehensive database class back in college and I'm no
stranger to RDBMSs but my work as a "full stack" web application software
engineer over the past few years has precluded me from doing a really "deep
dive" into MySQL for awhile. A primary goal of reading HPM was to build my
[tree trunk of understanding][tree_trunk_understanding] so that I could build
better software that makes use of MySQL.

The third edition was published in 2012, around the time of MySQL 5.5, so some
of the information is out of date. I found myself constantly googling claims
made throughout the book thinking "is that still the case?" Things move quickly
in the world of software so often the answer to my question was: no. That being
said, much of the information is still good even 8 years since publication.
High Performance MySQL is still worth a read for anyone that's looking to go
deep in their understanding. The authors, being as knowledgeable as they are,
are foretelling in many instances. For instance: the MySQL query cache (chapter
7), now deprecated and removed as of MySQL8, was always a bit problematic and
difficult to scale. The authors discuss this in detail and as a result I feel
pretty informed why it was formerly a part of MySQL and is no longer.

Being an 800 page book, the information covered is enormous and whatever I can
talk about in this post will be a random tiny sample. My guess is most
engineers have questions like "How can I write more performant queries" and
"How to optimize for indexes" and "How can I diagnose slow queries" and "What's
the most optimal schema design for my data and querying needs." From reading
this book I feel I can provide more confident answers to these questions.
That's in large part to the detailed explanations of B-Tree+ indexes, covering
indexes, hash indexes, clustering indexes, composite indexes, query parsing,
index selectivity, selective denormalization strategies, MySQL statistics and
logging, etc.

The book even recommends other books for guidance on how to go deeper. One
recommended book I plan to read is "Relational Database Index Design
and the Optimizers" by Mike Leach and Tapio Lahdenmaki. I particularly thought
the "3 star" ranking system for quantifying the utility of indexes was
interesting.

###### 3 Star Index Ranking
- 1 star: places relevant rows adjacent to each other
- 2 star: rows are sorted in the order the query needs
- 3 star: contains all columns needed (a "covering" index)

#### My top takeaways

##### Scaling
The book covers scaling in great detail. Both scaling up and scaling out
(vertical/horizontal scaling) through replication, partitioning and sharding.
With intelligent adaptation of these tools, the authors explain how MySQL can
scale to match the needs of giant services like Facebook or Wikipedia.

Most applications will never need to scale beyond a single database server. On
modern hardware MySQL can perform tens of thousands to millions of queries per
second, varying widely based on factors like read/write workload balance,
indexing, etc. ([Percona Sysbench Benchmarks][percona_sysbench]). Scaling "up"
or increasing server resources like I/O capacity, memory, CPU, network is the
easiest option for scaling and thus should be the first option employed. The
typical next step in scaling is setting up read replicas for an increase in
read capacity. After that, "functional partitioning" and then sharding can aid
in "scaling out."

##### Backups and Recovery
Also covered, backups and recovery. Naturally the database and the data within
is an incredibly critical part of any system. Having backups and restoring from
those backups is important. But creating backups can be time and resource
intensive, and restoring from backups can be complex. HPM explains the
differences and pros and cons of using logical or binary backups. It also
explains the folly of believing your replica or snapshots can be a backup.
There are many tools that can be useful for creating and restoring from
backups, HPM explains these in detail.

HPM also covers features of Linux operating systems like LVM snapshots and how
this is useful for creating "online" backups that don't interrupt OLTP workloads.

HPM explains point-in-time recovery, or the process of restoring the last full
backup and then replaying the binary log from that time forward (sometimes
called "roll-forward recovery").

I found the idea of "delayed replication" or intentionally having a delayed
replica for use in recovery very interesting. If you have a delayed replica and
you notice an accident before the delayed replica executes the offending
statement it can make recovery much faster.

##### Optimizing performance via schema design, indexing and query design
###### Schema
Chapters 4, 5 and 6 probably contain the information that's most sought after
by users of MySQL who want to get the most performance out of their database. 

Sometimes it makes sense to strategically denormalize data to suit querying
needs, or use cache tables to avoid continuously computing expensive statistics.

Using appropriate column datatypes can drastically improve performance as well.
General advice is to use the smallest datatype possible for every column and to
avoid NULL (p116). BLOB/TEXT columns almost always result in using on-disk
temporary tables for sorting operations and therefore can have serious
detrimental impact on performance and should be avoided if possible.

###### Indexing
Indexing is a complex topic and HPM breaks it down nicely. HPM explains the
differences between clustered indexes and secondary indexes. Often, if you
limit the number of columns returned by your query you can make use of
"covering indexes" which avoid the need for performing disk IO (p177).

Also explained is the subtle differences between redundant and duplicate
indexes. Duplicate indexes are wasteful and should be avoided but redundant
indexes can sometimes be useful.

Column order in indexes is also important and can aid queries or affect insert
performance. HPM explains how to quantify index column selectivity and design
indexes accordingly by choosing an optimal column index order.

Unused indexes are dead weight and HPM describes tools and strategies that can
be used to identify and remove them.

###### Querying 
MySQL executes queries using a nontrivial process and understanding this can
help users write queries that MySQL can execute more efficiently.

The query optimizer turns queries into "query execution plans" for the
database. The optimizer is "cost-based" and uses statistics about the data to
try and predict which execution plan will be least expensive.

Sorting operations can potentially be very expensive. Indexes can be very
helpful for sorting but when it's not possible to use them MySQL must sort the
rows itself. It can do this either in memory using the sort buffer, or on disk.
The data to be returned determines which type of sorting MySQL will perform.

Advanced users can provide "hints" to the optimizer to control the query plan.

The appendix contains a very detailed explanation of how to use `EXPLAIN` to
understand what plan the query optimizer will come up with.

From the book:
<pre>
"Optimization always requires a three-pronged approach: stop doing things, do
them fewer times, and do them more quickly."
</pre>

###### Percona
Lastly the authors are closely tied with Percona, a drop in replacement for
MySQL. Throughout the book they highlight the utility of their company's tools
(such as [percona-toolkit][percona_toolkit]) and the advantages/differences of
their "flavor" of MySQL.

High Performance MySQL was a great read and I'd recommend it to anyone looking
to enhance their "tree trunk of understanding" on MySQL/RDBMSs. It's old and
partially dated but still worthwhile to purchase. Through reading this book I
feel my focusing abilities have been enhanced and I'm looking foward to
tackling the next book in my queue.

[notes]: https://gist.github.com/cflynn07/46c564935607c444e2258b23f490168f "My notes from reading HPM"
[tree_trunk_understanding]: https://waitbutwhy.com/2016/05/mailbag-1.html "Tree Trunk of Understanding"
[percona_sysbench]: https://www.percona.com/blog/2017/01/06/millions-queries-per-second-postgresql-and-mysql-peaceful-battle-at-modern-demanding-workloads/ "Percona Sysbench Benchmarks"
[percona_toolkit]: https://www.percona.com/software/database-tools/percona-toolkit "Percona Toolkit"
