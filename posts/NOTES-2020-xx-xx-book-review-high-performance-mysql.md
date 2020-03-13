# Chapter 1

###### p1
- Reason for storage engine architecture is to separate query processing and data
  storage & retrieval

###### p4
- two lock types: shared locks, exclusive locks (read locks, write locks) - lock granularity is customizable in MySQL
    - table locks - lowest overhead
    - READ LOCAL - table lock, allows some types of concurrent write ops

###### SIDE THOUGHT
```
the more I read the better at efficient note taking I'm getting. Rather than
just writing down things immediately after reading, I'm reading sections and
revisiting to understand what highlights are best to grab.
```

###### p7
- `LOCK TABLES` mentioned as a way to get transaction like processing on tables
with engines that don't support transactions
  - https://www.oreilly.com/library/view/mysql-reference-manual/0596002653/ch06s07.html

###### p8
- isolation levels
  - side quest (https://blog.pythian.com/understanding-mysql-isolation-levels-repeatable-read/)
  - in MySQL 8, `tx_isolation` removed (https://dev.mysql.com/doc/refman/8.0/en/mysql-nutshell.html)
      - replaced with `transaction_isolation`
  - book states `REPEATABLE READ` *will* allow ***phantom reads***, but side
  quest blog post and testing on mysql8 show no phantom reads (at least in certain situations)
  ![](/static/images/2020-xx-xx-high-performance-mysql/isolation-level-repeatable-read.png)
      - this whole thing is a bit more complex than described in the book

###### p9
- innodb resolves deadlocks by rolling back transactions with fewest exclusive
  row locks

###### p10
- transaction logging / write ahead logging - remember this from college
  good description -> in-memory update first, then sequential IO to append log
  events on disk, later time table updated on disk (so data written to disk
  twice)

###### p11
- Can't mix transactional/non-transactional tables in transaction queries (InnoDB/MyISAM)
- two phase locking, usually locks obtained implicitly, you can explicitly obtain locks

###### p12
- MVCC, transactions consistent view of data, different transactions see
  different data same tables same time <- p8 blog post
- InnoDB MCCC 2 hidden values per row: created, expired ("system version number")
  - https://dev.mysql.com/doc/refman/8.0/en/innodb-multi-versioning.html
    refers to how the hidden system columns are treated re indexes
- system version number = a number that increments each time a transaction begins

###### p13
- .frm files for each table in database. According to side quest, frm files removed in mysql8
  - https://dev.mysql.com/doc/refman/8.0/en/data-dictionary-file-removal.html
  - data now stored in "data dictionary tables" https://dev.mysql.com/doc/refman/8.0/en/data-dictionary-schema.html

###### p16
- innodb created by Oracle before Oracle owned Sun/MySQL - later bundled into MySQL after Oracle bought mysql
- innodb stores data in series of files known as "tablespace"
  - https://mariadb.com/kb/en/innodb-file-per-table-tablespaces/ ".ibd" file extension

**I wanted to see .ibd files on disk for myself (running mysql in docker container)**
<pre class="prettyprint">
$ docker exec -it `cids mysql` ls -l /var/lib/mysql/crashcourse | sort
total 692
-rw-r----- 1 mysql mysql 131072 Jan 22 10:57 products.ibd
-rw-r----- 1 mysql mysql 131072 Jan 22 10:57 orderitems.ibd
-rw-r----- 1 mysql mysql 131072 Feb  9 10:02 orders.ibd
-rw-r----- 1 mysql mysql 114688 Jan 22 10:57 vendors.ibd
-rw-r----- 1 mysql mysql 114688 Feb  9 12:09 ordertotals.ibd
-rw-r----- 1 mysql mysql 114688 Feb  8 15:52 archive_orders.ibd
-rw-r----- 1 mysql mysql 114688 Feb  3 10:36 customers.ibd
-rw-r----- 1 mysql mysql   6144 Feb  9 10:02 productnotes.MYI
-rw-r----- 1 mysql mysql   4397 Jan 21 14:33 productnotes_350.sdi
-rw-r----- 1 mysql mysql   4020 Feb  9 10:02 productnotes.MYD
</pre>

###### p17
- Book recommends reading innodb transaction model/locking: https://dev.mysql.com/doc/refman/8.0/en/innodb-transaction-model.html
- innodb can support "hot" backups, other storage engines require halting all writes to table
- MyISAM used to be default (before 5.1)
- MyISAM no transactions or row level locks, not crash safe

###### p18
- productnotes.MYI and productnotes.MYD are the index and data files (see above ls -l)
- MyISAM can only be single file, therefore limited by disk space and largest allowed file on OS
- DELAY_KEY_WRITE
  - https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html#sysvar_delay_key_write
  - https://www.petefreitag.com/item/441.cfm info on when mysql closes tables (resulting in updating index on disk)
- Can compress MyISAM tables if table is never written to to reduce disk space and IO for fetching

###### p20
- CSV engine, I wanted to see it so:
<pre class="prettyprint">
cat <<- EOF | mysql -u root -ppassword -h 127.0.0.1 -v --table crashcourse &>/dev/null
DROP TABLE IF EXISTS testcsv;
CREATE TABLE testcsv(id INT NOT NULL, c VARCHAR(20) NOT NULL) ENGINE=CSV;
INSERT INTO testcsv(id, c) VALUES(1, "first");
INSERT INTO testcsv(id, c) VALUES(2, "second");
INSERT INTO testcsv(id, c) VALUES(3, "third");
EOF
docker exec -it `cids mysql` cat /var/lib/mysql/crashcourse/testcsv.CSV
1,"first"
2,"second"
3,"third"
</pre>

###### p22
- XtraDB drop in replacement for InnoDB, some enhancements
- PBXT - storage engine w/ MariaDB
- TokuDB - Fractal Trees data structures for indexes. "Big Data" storage engine.

###### p24
- if need full-text search, rec to use innodb w/ sphinx instead of MyISAM
- Recommended not to mix storage engines, will lead to bugs

###### p27
- authors seen InnoDB do fine 3-5 TB, single server not sharded. Beyond 10's of
  TB, data warehouse. Infobright or TokuDB

###### p31
- benchmarks across MySQL versions. Used `Cisco UCS 250` server

###### p32
- buffer pool, https://mariadb.com/kb/en/innodb-buffer-pool/
  - The most important server system variable is innodb_buffer_pool_size, which
  you can set from 70-80% of the total available memory on a dedicated database
  server with only or primarily XtraDB/InnoDB tables.

###### p36
- limitations of benchmarks, many artificial dimensions - differ from real world data

###### p37
- distinction between benchmarks and load testing
- two benchmarking strategies: the full stack or just MySQL(single component)

###### p38
- TPC-C standardized benchmark, widely quoted. http://www.tpc.org/tpcc/default5.asp

###### p39
- concurrency, should only care about benchmarking/testing "working concurrency"
- scalability: ideal system should get twice as much work done (twice as much
  throughput) when you double the number of workers

###### p40
- common benchmarking mistakes, design of benchmark differs from real world
  leading to unusable/inaccurate results

###### p41
- knowing when to push for realism or when to accept differences. Example
  application & database normally on different hosts but just running benchmark
  on same host because that's "good enough"

###### p42
- good source of realistic benchmark queries is from logging prod during a
  representative time period
- gotta run the benchmark for a while to observe system "steady state" perf

###### p45
- good idea sync to evenly divisable timestamps for data collection for easier
  correlation

###### p46
- MySQL default config settings tuned for tiny apps

###### p47
- null hypothesis https://en.wikipedia.org/wiki/Null_hypothesis

###### p49
- introduces `gnuplot`
- interesting how book does a HOW-TO for pulling metrics but hasn't covered
  using the tools to run a benchmark
- "furious flushing", averages don't show this but a graph will

###### p51
- lists popular benchmarking tools, jmeter most sophisticated
- mysqlslap

###### p53
- sysbench, all around favorite tool. Scripting in Lua. https://github.com/akopytov/sysbench

###### p57
Is my Mac much faster than the servers in this book?
<pre class="prettyprint">
$ sysbench --test=cpu --cpu-max-prime=20000 run
sysbench 1.0.19 (using bundled LuaJIT 2.1.0-beta2)

Running the test with following options:
Number of threads: 1
Initializing random number generator from current time


Prime numbers limit: 20000

Initializing worker threads...

Threads started!

CPU speed:
    events per second:   470.71

General statistics:
    total time:                          10.0006s
    total number of events:              4708

Latency (ms):
         min:                                    2.02
         avg:                                    2.12
         max:                                    3.36
         95th percentile:                        2.35
         sum:                                 9997.73

Threads fairness:
    events (avg/stddev):           4708.0000/0.00
    execution time (avg/stddev):   9.9977/0.00
</pre>

<pre class="prettyprint">
$ sysbench --test=fileio --file-total-size=150G prepare
</pre>

###### p60
- I'm wondering how useful running these benchmarks on a laptop are. I think we
  need to watch some live demos on youtube.

###### p69
- `profiling` introduced. This ch defines performance in a context to mean
  response time of an operation. Performance optimization = reducing response time for a given workload
- performance optimization contrasted to throughput optimization (queries per second)
  - increased throughput a side effect of performance optimization
- principle: cannot reliably optimize what you cannot measure

###### p72
- profiling tools all measure start/end time to deduce length. Construct "call graphs"
- pt-query-digest. https://www.percona.com/doc/percona-toolkit/LATEST/pt-query-digest.html
- execution-time profiling & wait analysis
- instrumentation
- plug for percona, touts time based measurement and instrumentation as being superior to MySQL

###### p74
- percona 5.0 slow query log. important causes poor performance: ex waiting for disk IO or row-level locks
- Amdahl's law?

###### p76
- Tom Kyte, guy who worked at Oracle for a long time and knows a lot about OracleDB

###### p77
- New Relic plugged
- PHP specific tools explained in great detail, which I mean... who TF wants to use PHP anymore

###### p80
- slow query log vs general query log
  - `long_query_time` can be set to 0
  - percona logs more to slow query log than MySQL

###### p81
- 2 alternative log gathering strategies:
  1. repeatedly poll SHOW FULL PROCESSLIST to see all queries (pt-query-digest --processlist)
  2. capture & inspect TCP network traffic (tcpdump && pt-query-digest --type=tpcdump)

###### p82
- pt-query-digest Query ID -> fingerprint -> hash of canonical (whitespace removed, lowercase)
- https://severalnines.com/database-blog/analyzing-your-sql-workload-using-pt-query-digest (other blog posts look interesting)

###### p85
- SHOW PROFILE
- SET profiling=1;
- https://dev.mysql.com/doc/sakila/en/
- SHOW PROFILES; SHOW PROFILE FOR QUERY n;

###### p88
- SHOW STATUS / SHOW GLOBAL STATUS, return counters. Session or global scoped
- Most important (expensive) counters:
  - handler counters + temporary file and table counters (Appendix B)

###### p89
- SHOW STATUS, counters of what server did. EXPLAIN is estimate of what server thinks it will do.

###### p91
- Authors pretty adamant performance_schema was only basic in MySQL 5.[5/6]. Investigate changes since then.

###### p95
- innotop https://www.percona.com/blog/2013/10/14/innotop-real-time-advanced-investigation-tool-mysql/

###### p96
- advice: get used to gnuplot or R for graphing quickly

###### p98
- pt-stalk, good threshhold monitoring tool
  - https://www.percona.com/doc/percona-toolkit/LATEST/pt-stalk.html
  - pt-sift useful for analyzing results

###### p99
- oprofile - https://en.wikipedia.org/wiki/OProfile
- strace https://en.wikipedia.org/wiki/Strace
- GDP stack traces
  - pt-pmp tool
- pt-collect, tool intended to be used from pt-stalk
- nm tool
- debug symbols https://stackoverflow.com/questions/3694900/can-you-explain-whats-symbols-and-debug-symbols-in-c-world

###### p100
- most useful things to look at: query/transaction behavior and server internals behavior
- pt-mysql-summary and pt-summary, 2 other useful tools

###### p101
- http://poormansprofiler.org/ & pt-pmp
- video on oprofile https://www.youtube.com/watch?v=-fjujEUJZuE

###### p104
- SHOW INNODB STATUS is now `SHOW ENGINE INNODB STATUS`

###### p104
- iostat, vmstat tools mentioned

###### p106
- furious flushing/checkpoint stall.
- InnoDB performs certain tasks in the background, including flushing of dirty pages from the buffer pool. Dirty pages are those that have been modified but are not yet written to the data files on disk.
  - https://dev.mysql.com/doc/refman/8.0/en/innodb-buffer-pool-flushing.html

###### p108
- use df-h and lsof to monitor disk space and file descriptors to see how much data MySQL is writing to disk
- MySQL writing 1.5GB to temporary tables

###### p109
- could be a "cache stampede"

###### p110
- INFORMATION_SCHEMA.USER_STATISTICS (%_STATISTICS) additional table included w/ Percona & MariaDB

###### p111
- using `_statistics` tables, can find unused indexes which are candidates for removal
- strace, profiles/intercepts system calls. Can filter by process id
  - similar but different to oprofile
  - strace will bring mysqld to a crawl when attached

###### p118
- floating-point vs exact math operations
- DECIMAL has higher space and computational costs. Use only when need exact results for fractional numbers.

###### p119
- varchar(255), any length greater than 255 will use 2 bits to store length

###### p121
- padding and trimming behavior consistent across all storage engines since
  it's handled above the storage engine by mysql

###### p122
- max_sort_length for BLOG/TEXT

###### p125
- primary key size affects other indexes (explained more next ch)

###### p127
- mysql 8 does support subsecond times https://dev.mysql.com/doc/refman/8.0/en/fractional-seconds.html
- inoodb no space saved using BIT data type, engine uses smallest INT type that works

###### p130
- UUID values good idea convert to 16 byte numbers with UNHEX() and store them in BINARY(16)
- UUID/SHA/MD5 slow INSERTs, value goes to random location in index

###### p132
- performance consequences many columns, storage engine and mysql server pass
  rows in "row buffer format" and some conversion occurs
- EAV design pattern bad choice w/ MySQL, limit 61 tables per join and it requires a lot of self joins

###### p135
- normalization, denormalization, 2NF, 3NF. Example of partial denormalization to avoid expensive joins

###### p136
- common denormalization - duplicate (cache) columns. Use triggers to keep in sync
- caching 'derived' values, ex: num_messages

###### p138
- technique: using different engines for cache tables. innodb -> myisam
- shadow tables for background rebuilding of cached tables, atomic rename to swap them
- materialized views, flexviews https://docs.huihoo.com/mysql/percona/live/mysql-conference-2015/Materialized-Views-for-MySQL-using-Flexviews.pdf
  - type of cache table, can calculate changes incrementally from source data (using deltas to compute)

###### p140
- query w/ USING() - alternative join method (opposed to "ON")
- trick to get higher concurrency with counter tables: use multiple rows
  instead of 1 to achieve higher update/write concurrency. A single row will
  have all writes serialized

###### p141
- alter table techniques.
  1. switch around servers (build new table in offline server then swap)
  2. "shadow copy" - build new table next to existing one, then rename and drop. Book lists tools useful for this.

###### p142
- ALTER TABLE -> ALTER column vs MODIFY column, ALTER can be much faster (no table rebuild) (see also CHANGE COLUMN)
- random aside reading on replacement for FRM files in MySQL8 https://www.percona.com/blog/2016/10/03/mysql-8-0-general-tablespaces-file-per-database-no-frm-files/
  - schema per customer ability sounds interesting

###### p143
- MyISAM trick: disable keys, load data, reenable keys (will load faster) (goal: build indexes by sorting)
  - myisam builds indexes in memory, very slow if no available memory. Unique indexes always built in memory

###### p148
- B-Tree indexes, usually the "default"
- InnoDB uses B+Tree indexes (each leaf has pointer to next leaf)

###### p149
- MyISAM indexes refers to rows by physical storage location, InnoDB refers to primary key

###### p151
- list of query types that use b-tree indexes (match full value, left prefix, range, one part exactly partial second part)
- column ordering in indexes important

###### p152
- only Memory tables support HASH indexes
- hash indexes are very compact (indexes only short hash values)

###### p153
- HASH indexes can't be used for sorting
- innodb creates "adaptive hash indexes" above b-tree indexes, frequently accessed values go in adaptive hash index (in memory)
    - https://www.percona.com/blog/2016/04/12/is-adaptive-hash-index-in-innodb-right-for-my-workload/

###### p154
- using a column that's value is a hash, you can create shorter b-tree indexes (example urls)
  - using triggers on insert/update you can automatically assign
- https://en.wikipedia.org/wiki/Cyclic_redundancy_check

###### p158
- Book recommendation: Relational Database Index Design and the Optimizers
- 3 star ranking system for database indexes
  - 1 star: places relevant rows adjacent to each other
  - 2 star: rows are sorted in the order the query needs
  - 3 star: contains all columns needed

###### p159
- terabyte scale indexes break down, per block metadata better (Infobright)
- columns must be "isolated" in queries to use indexes (not in a function or expression)

###### p160
- prefix indexes, saves space but less selective
- index selectivity: distinct indexed values (cardinality) / total rows (max = 1)

###### p162
- shows how to calculate average selectivity for prefix indexes of different lengths

###### p163
- prefix indexes can't be used for order by or group by or covering indexes (composite indexes)
- individual indexes on columns wont help most queries ("index merge" helps a bit)
  - https://www.percona.com/blog/2009/09/19/multi-column-indexes-vs-index-merge/
  - As a summary: Use multi column indexes is typically best idea if you use
    AND between such columns in where clause. Index merge does helps
    performance but it is far from performance of combined index in this case.
    In case you’re using OR between columns – single column indexes are
    required for index merge to work and combined indexes can’t be used for
    such queries.

###### p164
- using index merge, for AND conditions, when server intersects indexes a single covering index with relevant columns prob better
- union indexes for OR conditions can be CPU/memory intensive (buffering, sorting, merging), esp if indexes not very selective.
- optimizer only accounts for number of random page reads

###### p165
- EXPLAIN -> index merge, examine query and table structure
- column order very important multicolumn indexes
- Column order important. When no sorting/ordering, most selective column first good. When ordering, first column should aid ordering

###### p166
- orders table example with staff_id, customer_id. logically, staff_id appears
  a lot more often. for queries that select for staff_id and customer_id,
  create index with customer_id first (higher cardinality/selectivity)
- example situation with specific values higher than normal cardinality: guest
  users all sharing a single user id in a sessions table

###### p168
- clustered indexes. Really b-tree indexes that also contain row data. Oracle calls them index-organized tables.
  - normally primary key in innodb
  - works best I/O bound workloads, if data small enough to fit in memory it doesn't matter
  - insert speeds depend on insertion order.
  - secondary indexes require two lookups. Secondary indexes contain primary key values, not pointers too rows. Must do 2 b-tree lookups.

###### p178
- "covering index": index that contains all data requested by query
- EXPLAIN on query that is covered by covering index will have "Using Index" in Extra column

###### p180
- Deferred join optimization
```
EXPLAIN SELECT * FROM products WHERE actor='SEAN CARREY' AND title LIKE '%APOLLO%'
```
```
EXPLAIN SELECT *
FROM products
    JOIN (
      SELECT prod_id
      FROM products
      WHERE actor='SEAN CARREY' AND title LIKE '%APOLLO%'
    ) AS t1 ON (t1.prod_id=products.prod_id)
```
^ uses covering index in subquery, more efficient when reducing number of full rows that must be read, but not so few rows that
the subquery becomes a performance detriment

###### p181
- innodb secondary indexes leaf node values contain primary keys, so secondary indexes can also "cover" selections for primary key
```
# EXAMPLE: secondary index on "last name"
SELECT actor_id, last_name FROM sakila.actor WHERE last_name = 'HOPPER' # This is "covered" by index
```

###### p182
- MySQL 5.6 "index condition pushdown" https://dev.mysql.com/doc/refman/5.6/en/index-condition-pushdown-optimization.html
- two ways to produce ordered results, sort operation or scan an index in order
  - EXPLAIN -> type: "index" (not confused with "using index" in extra column)
- can use index to sort on non-leftmost column of an index if the leading columns are constants

###### p184
- MyISAM packed indexes, reduce index size
  ex: 1. perform
      2. 7,ance (performance)
  - downside, can't do binary searches. Must scan block from start

* random side blog post thought. 10 pages, 1 disk page my brain
* aware enough to point where I THINK I can have my brain ping me to come back to some of this stuff in the future when I need it
* sometimes read a whole page, touch every word, then stop and say "wait what did I just read"
* could devote entire life to thoroughly understanding this, I want to be a generalist tho

###### p185
- mysql will let you shoot yourself in the foot and create duplicate/redundant indexes
- redundant & duplicate bit different
  - index on (a, b) and index on (a) - redundant. leftmost prefix

###### p186
- example beneficial redundant index. One is covering one is not. Helps 2 queries. Cost is insertion speed

###### p187
- identify redundant/duplicate indexes:
  - common_schema by Shlomi Noach
  - pt-duplicate-key-checker tool https://www.percona.com/doc/percona-toolkit/LATEST/pt-duplicate-key-checker.html
  - pt-upgrade tool
  - find unused indexes tool: pt-index-usage

###### p188
- example how innodb can lock more rows than it needs by not being aware of filters from WHERE part of query

###### p189
- shared read locks on secondary indexes, but exclusive write locks only with primary key

###### p190
- index based sorting vs post-retrieval sorting
- thinking about queries on hypothetical dating site: idea create indexes prefixed with most common filter columns (sex,country,etc)
  - index can still by usable to query that doesn't use sex by adding `AND sex in('m', 'f')`

###### p191
- (sex, country, age) + (sex, country, region, age) + (sex, country, region, city, age) can reuse with IN() trick and scrap first 2
- general principle, keep range criterion (age) at end of index, optimizer use index as much as possible

###### p192
- range vs multiple equality conditions. Both appear in EXPLAIN:type as "range" - multiple equality doesn't ignore any further columns
- cron job compute "active" column to replace last_online range in query for better indexing

###### p193
- example query difficult for index to optimize
  `SELECT <cols> FROM profiles WHERE sex='M' ORDER BY LIMIT 100000000, 10;`
  high offset requires scanning a lot of data, index can't really help.
  deferred join potentially helpful. Use covering index to retrieve just primary key columns then join

###### p195
- innodb corruption should be rare. Usual culprit is someone messing with files outside mysql (rsync)
- ANALYZE TABLE regenerates index stats

###### p197
- innodb recalcs indexes ANALYZE TABLE || size change (1/16 or 2 billion rows whichever first)
- percona allows users to pause stats resampling `innodb_stats_auto_update`
- `innodb_analyze_is_persistent` option speed up system start by persisting statistics to disk
- 3 types of data fragmentation: row fragmentation, intra-row, free space

###### p198
- OPTIMIZE TABLE can defragment
- MySQL/Percona differ in index defragmentation, OPTIMIZE TABLE only defrags clustered index not secondary. `expand_fast_index_creation`

* Query/index/schema optimization confluence

###### p203
- fetching all columns precludes using covering indexes

###### p204
- simplest query cost metrics: response time, num rows examined, num rows returned
- response time = service time + queue time
- QUBE (quick upper bound estimate) - RDID and the O (wiley)

###### p205
- rows examined / rows returned. Joins require at least 2 rows examined for every row returned.
- row access methods are "type" column in EXPLAIN
  - full table scan > index scan > range scan > unique index lookups > constants

###### p208
- "chopping up" queries. Deleting in a loop (batch 10000) while affected rows > 0 minimize impact on server
- pt-archiver

###### p209
- merits of breaking up multi-table JOINS into discrete queries

* Learning how MySQL optimizes and executes queries lets you reason from principles

###### p210
- mysql protocol half-duplex (any given time sending or receiving)
- client sends query as single packet. `max_allowed_packet`

###### p215
- how expensive optimizer estimated query
```
SELECT SQL_NO_CACHE COUNT(*) FROM sakila.film_actor;
SHOW STATUS LIKE 'Last_query_cost';
```
^ estimate of disk page reads required to execute query

###### p216
- static & dynamic optimizations

###### p219
- IN() usually much faster than series of OR. IN() fast binary search. `O(log n)` vs `O(n)`

###### p220
- MySQL every query considered a "join"
- execution strategy involves nested loop for every join

###### p221
- good pseudocode logic INNER and LEFT OUTER joins

###### p222
- swim-lane diagram for visualizing an INNER JOIN
- subqueries create temporary tables for which there are no indexes (performance consideration)
- RIGHT OUTER JOIN queries are rewritten to equivalent LEFT OUTER JOINS
- FULL OUTER JOIN not supported by MySQL due to incompatibility of nested loop join logic
- `EXPLAIN EXTENDED` and `SHOW WARNINGS` https://dev.mysql.com/doc/refman/8.0/en/explain-extended.html

###### p223
- multitable join can be represented as a tree (ex on page)
- "join optimizer" most important part of query optimizer. Estimates cost diff plans, chooses the least expensive one.

###### p224
- Oracle lingo "driver into table"
- First row in EXPLAIN output is start point of query plan
- `STRAIGHT_JOIN` forces join to proceed in order specified in query

###### p226
- join over n tables will have n! combinations of join orders -> "search space" possible query plans. ex 10! = 3,628,800
- `optimizer_search_depth`, "greedy" searches shortcut
- when no index for sorting, "filesort" (can be in memory or on disk)
  - if on disk, sorts values in chunks. If in memory (sort buffer fit) quicksort
- filesort algos: two passes (old), single pass (4.1+ avoids reading rows twice and less random IO more sequential IO, uses more sort buffer mem)

###### p227
  - `max_length_for_sort_data` & size of query columns & order by columns determines which algo is used
- when sorting, mysql allocates fixed-size record for each sorted row that is the max possible size (think what happens if column type huge varchar)

* advantage of reading book over sequential days, brain state becomes "in the zone" or flow

###### p228
- execution plan not bytecode in mysql, data structure
- typically, optimization stage more complex than execution stage
- there's a "handler API" class for storage engine + table. For each table in query MySQL creates an instance of this class.
  - instance has API methods (ex read first row in index, next row)

###### p229
- `SQL_BUFFER_RESULT` config what interval computed results are sent at while computing
- each result row = different packet (mysql client/server protocol)
- correlated subqueries (pushes condition into subquery to help optimize)

###### p230
- to prevent rewrite to correlated subquery, can rewrite query as a JOIN

###### p232
- Extra: Not exists, optimization - stops processing current row when it finds a match
- subquery optimal situation DISTINCT INNER JOIN -> EXISTS() subquery (avoids temporary tables)

###### p233
- "push down" conditions like LIMIT or ORDER BY into UNION (mysql can't automatically)
- https://www.iheavy.com/2013/06/13/how-to-optimize-mysql-union-for-high-speed/

###### p234
- MySQL no parallel execution
- MySQL8 has limited parallel query execution now: https://www.percona.com/blog/2019/01/23/mysql-8-0-14-a-road-to-parallel-query-execution-is-wide-open/
- No hash joins, MySQL8 has hash joins
  - https://dev.mysql.com/doc/refman/8.0/en/hash-joins.html
  - https://mysqlserverteam.com/hash-join-in-mysql-8/

###### p237
- SELECT & UPDATE on same table not allowed, can get around using derived/temporary table

###### p241
- pt-upgrade tool, validate queries run well on new MySQL versions

###### p242
- always use `COUNT(*)` and don't name columns when you just want to know the number of rows in a result

###### p243
- `SELECT SUM (IF(color = 'blue', 1, 0)) AS blue ...`

###### p244
- COUNT() queries hard to optimize in general. Maybe use summary tables, memcached, etc

###### p245
- When grouping on lookup table, GROUP BY on primary key more efficient

###### p246
- Suggest `ONLY_FULL_GROUP_BY` prevent different GROUP BY columns from SELECT
- https://www.mysqltutorial.org/mysql-rollup/

###### p248
- LIMIT OFFSET query trick, instead of SQL_CALC_FOUND_ROWS (expensive) fetch 21 rows, if 21st row returned render "next" link

###### p249
- pt-query-advisor, sort of "lint checker" - parses log of queries

###### p251
- quadratic algorithms https://medium.com/@verdi/the-quadratic-sorting-algorithms-an-overview-of-bubble-selection-and-insertion-sort-266de2b26004

###### p253
- variable issues arise frequently when assigning and reading at different stages of a query. Rec is to assign and read in same stage.
- Book Rec: SQL and Relational Theory: How to Write Accurate SQL Code - CJ Date

###### p256
- SELECT .. FOR [SHARE/UPDATE] https://dev.mysql.com/doc/refman/8.0/en/innodb-locking-reads.html

###### p257
- CONNECTION_ID() function

###### p260
- the merits of "fudging" location proximity calculations

###### p261
- Index that isn't the whole truth but gets you close to the truth cheaply

###### p266
- Partitioned tables = multiple discrete tables (each handled by storage engine) logically connected by MySQL

###### p267
- MySQL syntax partition on column preferable to function over column

###### p273
- Can partition by expressions but must query/search by column

###### p274
- Merge storage engine https://dev.mysql.com/doc/refman/8.0/en/merge-storage-engine.html

###### p277
- querying VIEWS, 2 algos: MERGE and TEMPTABLE
  - TEMPTABLE used with VIEW def uses constructs that don't preserve 1-to-1
    relationship between the view rows and underlying table rows

###### p288
- View that doesn't preserve 1-to-1 not updateable (TEMPTABLE algo)
- CHECK OPTION: any rows changed must continue to match WHERE clause

###### p289
- Views useful refactoring schema in stages, code can continue accessing old schema

###### p282
- ISO/IEC 9075-4:2016 https://www.iso.org/standard/63557.html

###### p283
- common stored procedure example, transfer funds at bank. All done in transaction. No access to underlying tables.

###### p286
* Always checking limitation claims
- Only 1 trigger per event per table NO LONGER TRUE (5.7+) https://www.mysqltutorial.org/mysql-triggers/create-multiple-triggers-for-the-same-trigger-event-and-action-time/

###### p287
- Triggers occur within same transaction as op that triggered them

###### p288
- GET_LOCK() & RELEASE_LOCK()

###### p289
- Using version number comments to preserver stored procedure comments hack

###### p290
- important to note MySQL cursors execute entire query on open
- in memory tables do not support BLOB/TEXT, if cursor on table with these types temp table will go to disk

###### p291
- prepared statements: client sends "prototype" which is stored on server and server returns "handle"
  - only parameters sent for each execution
  - helps with security, don't need to escape values in client app

###### p295
- UDF (since "ancient" times)
- http://www.mysqludf.org/
- book doesn't show complete example of how to use or load a UDF

###### p299
- Character sets & encodings, reviewed this old blog post https://www.joelonsoftware.com/2003/10/08/the-absolute-minimum-every-software-developer-absolutely-positively-must-know-about-unicode-and-character-sets-no-excuses/
- "hello" -> U+0048 U+0065 U+006C U+006C U+006F (unicode code points) -> UTF8 (system for storing code points in memory)
-                                                 |--> encodings (characters to display per code point)

###### p303
- CHAR(10) encoded with UTF-8 uses 30 bytes (MySQL supports UTF8 subset using 3 bytes... at least it did when this book was written)

###### p304
- LENGTH() and CHAR_LENGTH() - return bytes and character lengths
- MySQL will automatically use prefix index if user specified index is too long
- Can see how multibyte character sets can use more space when sorting in memory and such

###### p306
- full text index -> full-text collection
- special B-tree index w/ 2 levels. L1 keywords. L2 associated "document pointers"
- prunes "stopwords" and short/long words

###### p308
- boolean full-text searches, specify increased/decreased relative importance of words

###### p310
- full-text indexes are expensive maintain for INSERT/UPDATE/DELETE ops
- full-text prone fragmentation

###### p311
- Techniques hacking full text indexes to more efficiently search by author and coordinates

###### p314
- distributed transactions
- aside on XA https://www.percona.com/blog/2018/05/16/mysql-xa-transactions/
- Internal vs External XA. Another way to sync data multiple servers besides replication
- query cache stores complete results of SELECT queries
- query cache is GONE in MySQL8 https://www.digitalocean.com/community/tutorials/how-to-optimize-mysql-with-query-cache-on-ubuntu-18-04
- https://mysqlserverteam.com/mysql-8-0-retiring-support-for-the-query-cache/

* Kinda glossed over query cache section after learning it's no longer in MySQL
* sounds weird to say, but reading through an 800 page book is kinda the lazy way to learn a subject. Someone else has gathered everything and laid it all out in a nice order.

###### p323
- good ratio of cache hits/inserts 10:1

###### p324
- `query_cache_min_res_unit` tradeoff wasting memory/cpu cycles

###### p329
- stored procedures vs stored functions https://www.a2hosting.com/kb/developer-corner/mysql/mysql-stored-functions-and-procedures
  - stored function can be used to compute a value in a query
  - stored procedure is CALL'd to do something

###### p330
- Ideal mysql configuration function of server hardward, workload, data, app requirements (not just hardware)
- using defaults = safety of numbers
- /etc/my.cnf || /etc/mysql/my.cnf

###### p334
- can set variable values with suffix for units in cli args or config file but not with SET command (ex 1M - 1 megabyte)
- SET DEFAULT can restore session scoped variables back to what they were when server started

* Attempting to focus like a machine in a book makes me aware of my impulses to do things like check instagram. Normally I'd just do the impulse not catch it. uncaughtException()
* When brain feels tired possible to literally force its circuits to "warm up" and absorb faster

###### p337
- Example of dynamically adjusting a variable `sort_buffer_size` for a specific query
- https://bugs.mysql.com/bug.php?id=37359

###### p340
- pt-log-player - replay query log against server for benchmarking
- tuning by ratio bad. tuning scripts also bad. "tuning" phrase also bad. Internet has tons of bad advice.
```
SET @crash_me_1 := REPEAT('a', @@max_allowed_packet);
SET @crash_me_2 := REPEAT('a', @@max_allowed_packet);
```
- what the hell is the joke in the footer

###### p342
- for even more than "few megabytes of data" need to configure mysql. It assumes it's not the only thing running on a system. (LAMP days)

###### p345
- not crucial to get right immediately, can start with something larger than default but still safe and test

###### p347
- managing memory important. Memory 2 categories: can control and can't control

###### p350
- innotop for monitoring innodb buffer pool https://www.percona.com/blog/2013/10/14/innotop-real-time-advanced-investigation-tool-mysql/
- `SHOW INNODB STATUS`
- `innodb_max_dirty_pages_pct`

###### p351
- feature to reload pages after restart helps server "warm up" faster. Especially useful on replicas.
- init_file setting, SQL to run when server starts up. Some people run full table scans to load indexes into buffer pool.

###### p352
- `CACHE INDEX` command

###### p353
- "read-around writes." Good idea make MyISAM block size same as OS

###### p354
- good setting for `thread_cache_size` is fluctuation in connections (`Threads_connected`) variable

###### p355
- table cache not used for much with innodb tables
- https://www.got-it.ai/solutions/sqlquerychat/sql-help/data-definition/how-to-optimize-mysql-table-cache/
- https://www.percona.com/blog/2012/03/23/how-flush-tables-with-read-lock-works-with-innodb-tables/A
- https://jameshfisher.com/2016/12/24/perror/

###### p356
- innodb per table cache: "table definition cache" || "data dictionary"

###### p357
- tradeoff data security / performance
- most important things: InnoDB log file size, how InnoDB flushes log buffer, how InnoDB performs I/O

###### p359
- `innodb_recovery_stats` opt w/ percona
- moonitor innodb log & log buffer I/O `SHOW INNODB STATUS` -> LOG section. `Innodb_os_log_writte`

###### p360
- `innodb_flush_log_at_trx_commit` 0, 1, 2
- important to know difference writing to log buffer vs flushing to durable storage

###### p361
- srsly need RAID volume with battery backed write cache
- percona, as usual, better because can offer `innodb_flush_log_at_trx_commit` at session level
- https://linux.die.net/man/2/fdatasync

###### p363
- importance of having write cache if using O_DIRECT

* visual representation of understanding of a topic before/after book. Block with missing bits vs solid.

###### p364
- innodb tablespace, virtual filesystem spanning one or more files on disk
```
innodb_data_home_dir = /var/lib/mysql/
innodb_data_file_path = ibdata1:1G;ibdata1:2G;ibdata3:1G;
```
- recommendation to disable autoextend feature

###### p366
- possibility of storing files on raw partition for minor performance gain, very annoying tho

###### p367
- `SHOW INNODB STATUS` command see number of unpurged transactions (sub trx id counter from purge done)
- sometimes necessary to throttle writes to let purge process keep up (avoid tablespace growing indefinitely)

###### p368
- doublewrite buffer. uses checksum to check for corruption

###### p369
- might want to disable doublewrite on replicas

###### p372
- innodb has "thread scheduler" controls how threads enter it's kernel. `innodb_thread_concurency`
  - theoretical good concurrency value: concurrency = Num CPUs * Num disks * 2 (in practice better lower)

###### p373
- innodb 2 phase process for thread entrance. helps avoids context switch overhead
- `innodb_concurrency_tickets`

###### p375
- particularity of BLOB workloads: server cannot use in-memory temporary tables
  - can convert temp to varchar w/ SUBSTRING() to get in-memory tables

###### p376
- memory-based filesystem https://www.jamescoyle.net/knowledge/951-the-difference-between-a-tmpfs-and-ramfs-ram-disk
- tmpdir - server setting controls where temporary tables are placed
- 768 byte prefix for variable length columns, for index. if too large innodb may use external storage (outside cluster index)

###### p377
- External storage use disables use of adaptive hash index

* Thinking back to that one project I did in college where I stored images in MySQL for no good reason

###### p378
- `tmp_table_size` & `max_heap_table_size` + counters: `Created_tmp_disk_tables` & `Created_tmp_tables`

###### p381
- Not documented, `max_allowed_packet` controls size of user defined variable
- max_connect_errors and host blacklisting

###### p389
- hyperthreading used to suck but it's ok now
- Fast CPU -> low latency
- Many CPI -> High throughput

###### p390
- innodb/myisam global locks may limit effectiveness of multiple CPUs
- stack traces help diagnose. `pt-pmp` tool
- x86 arch dominant, 64 bit cool

###### p391
- two types of concurrency problems:
  - logical concurrency issues
  - internal concurrency issues

###### p393
- spatial and temporal locality of reference

###### p394
- Random and sequential IO. Random I/O benefits most from caching

* training most efficient, quick reading of text as possible while achieving a high level of understanding
* forcing brain not to go off intended path into every little thought and
  impulse (instagram, thing I want to google) - sometimes brain does this
  without me even knowing it happened
* maximize "endurance" - how much info I can parse every day. I notice iterative improvement

- sequential IO can be 5000x faster (example in book)

###### p395
- The speedy nature of sequential writes is the impetus for "write-ahead logging" strategies
  - converts random IO to sequential

###### p396
- concept of "working set", can be defined as time-based percentage. cache units (16kb in innodb)

###### p397
- memory/disk ratio. cache miss rate. Can infer from % time CPU used.
- nonlinear relationship memory to cache miss rate (due to things like page size etc)

###### p398
- latecy/throughput of disks relates to access time and transfer speed
- non-linear relationship of system resources to performance
- good cache miss rate is a "fuzzy" concept

###### p400
- solid state drives faster at random and sequential IO. This is important to normalized databases
  - slight improvement in sequential IO, huge improvement in random IO (since spin disks suck at random IO)

###### p401
- SSDs have interesting characteristics due to physics. Write amplification, wear leveling. garbage collection. Can slow down when full.
  - writes in big blocks (512kb). writes take hundreds of microseconds, erases take milliseconds (slower)

###### p402
- single-level-cell vs multi-level-cell. how many bits of storage per cell. Trade off density&speed vs reliability/durability

###### p408
- flashcache

###### p411
- innodb logs better suited for RAID w/ battery backed write cache and spinning disks... it's sequential write only and occasionally sequential read (recoveries)

###### p412
- suggestion to disable doublewrite buffer on flash storage because it's not necessary... interesting

###### p414
- hardware considerations for replica depend on whether it's being used for failover or additional read capacity

###### p417
- RAID doesn't eliminate need for backups.
- 2 disks can actually fail simultaneously

* Actually reading out loud seems to be helpful

###### p418
- background patrol read with RAID, avoid latent data loss
- hot spare hard drive for when failures occur
- raid cache batteries sometimes schedule discharge/recharge learn cycles and disable the cache, can affect perf

###### p419
- hardware limitations, bugs or misconfiguration can cause serialization of IO instead of parallel, leading to shitty performance
- use sysbench to benchmark hardware to validate expected performance

###### p420
- good value in aligning all components of filesystem, innodb page size, filesystem blocks, LVM, partition offset, etc...

###### p421
- RAID controllers have a cache. It's pretty scarce. Using it for reads is
  usually pointless since caching at higher levels accomplishes the same thing
  but better. Most useful for caching writes.

###### p423
- SAN's have interesting perf characteristics. Good at sequential reads/writes, bad at random. High latency due to distance.

###### p424
- SAN's can be a good match for MySQL provided storage/retrieval carefully planned (lots of memory, avoid random IO)

###### p430
- https://en.wikipedia.org/wiki/Sysctl
- https://en.wikipedia.org/wiki/Multi_Router_Traffic_Grapher
- network trunking (connect multiple NICs)

###### p435
- MySQL needs kernel threads > userland threads (access multiple CPUs)

###### p436
- For obvious reasons, swapping is bad for MySQL

###### p438
- MySQL memlock option can prevent paging but will crash if memory runs out
- vmstat useful tool. Most output in blocks/sec

###### p441
- "Little's law" https://en.wikipedia.org/wiki/Little%27s_law
- using `iostat` to detect serialization

###### p442
- other good tools: dstat, collectl, mpstat, blktrace, pt-diskstats

###### p443
- write requests can be buffered, read requests are synchronous

###### p447
- replicas key for H/A and performance
- statement-based replication and row-based replication
- all replication is asynchronous, can't guarantee up-to-date

###### p448
- replication backward compatible, newer versions can be replicas of older versions

###### p449
- replication: master binary log -> replica copies to relay log -> replica replays relay log to its own data

###### p450
- replica initiates binlog dump process on master, non-polling

###### p451
- replication is serialized on replica (single thread)

###### p460
- statement-based replication simple, low bandwidth, doesn't work well with triggers

###### p461
- mysql can switch row-based and statement based replication dynamically depending on situation

###### p462
- statement based replication can be better if working with different schemas

###### p465
- `log_slave_updates` use replica as master to other replica
- globally unique server ID's important preventing infinite loop

###### p466
- replication filters, binary log on master, relay log on replica

###### p470
- master-master fraught with merge conflicts, super annoying

###### p471
- master-master active-passive mode 比較好用, a

###### p473
- master-master w/ replicas
- master-master is special case of *ring replication*
- ring replication depends on every node in ring working, otherwise failure (infinite loop possible if server removed)

###### p475
- each replica creates new thread on master that runs `binlog dump`
- distribution master, offload serving binary logs
  - distribution master can use blackhold storage engine (doesn't perform queries on disk)
  - `slave_compressed_protocol`

###### p476
- distribution replicas wil make promoting replica to master more tricky, binlog coordinates will be different
  - `pt-heartbeat` can help (measure mysql replication delay)
  - https://www.percona.com/doc/percona-toolkit/LATEST/pt-heartbeat.html
- pyramid design useful

###### p477
- selective replication, horizontal partitioning but for replicas only (master has all data)

###### p478
- selective replication, easiest way: separate databases on master
- OLTP vs OLAP (large, low, don't require 100% up-to-date data)

* Break from reading to learn about MySQL at Facebook & MyRocks
* https://engineering.fb.com/core-data/myrocks-a-space-and-write-optimized-mysql-database/

###### p482
- replication capacity scaling diminishing returns. Hard to scale writes. Reads easier.

###### p487
- measure replication lag with heartbeat record. `pt-heartbeat` useful

###### p489
- `pt-table-sync` - fix replication servers when out of sync

###### p493
- recommendation, if master crashes:
  1. promote replica to new master, 2. grab binlog off disk of old master to catch up new master, 3. allow client connections to new master
- also consider using SANs or DRBD (distributed replicated block device) to preserve data

###### p497
- examining bytes of binlogs to solve corruption issues (check out more info on this)

###### p499
- issue arise replication mixing transactional/non-transactional tables
- nondeterministic statements can lead to out of sync data on replicas

###### p502
- temporary tables don't work well with replication
- workaround for temporary tables, use separate database w/ connection namespaced tables (`pt-find` helps remove)

###### p503
- `INSERT ... SELECT` locking can be contentious, alternative `SELECT INTO OUTFILE` and `LOAD DATA INFILE`

###### p505
- close quote: "writing to both masters in master-master replication is a terrible idea"
- possible to "interleave" auto increments on 2 servers

###### p507
- good idea design apps tolerate replication lag
- https://www.slideshare.net/JeanFranoisGagn/how-bookingcom-avoids-and-deals-with-replication-lag

###### p508
- avoiding large group-by queries on all replicas by moving off master to single replica and then pushing results back to master

###### p515
- semisynchronous replication (transaction commit doesn't complete on master until at least one replica has the data)

###### p517
- tungsten replicator https://docs.continuent.com/
- tungsten allows mutli-master clusters, certain nodes tagged as system of record for certain data

###### p525
- linear scaling diagram, most systems slightly less than linear scaling
- Dr. Neil J. Gunther's Universal Scalability Law
- Guerrila Capacity Planning
- Forecasting MySQL Scalability with the Universal Scalability Law
- https://www.percona.com/resources/white-papers/forecasting-mysql-scalability-universal-scalability-law
- Amdahl's Law https://en.wikipedia.org/wiki/Amdahl%27s_law
  - predict the theoretical speedup when using multiple processors.

###### p526
- Eliyahu M. Goldratt "The Goal"

###### p531
- scale-out tactics: replication, partitioning, sharding
- simplest & most common: replication (use replicas for read queries)
- nodes can be: master-master replication pair, master & many replicas, active server distributed replicated block device, SAN

###### p532
- functional partitioning, separate areas of functionality to different MySQL servers (ex forums, news, support, users)
  - obvious issue here is each "functional area" must scale vertically if trying to use only single MySQL node

###### p534
- to scale write capacity must partition data

###### p535
- primary sharding challenge is choosing a partitioning key

###### p536
- good sharind partitioning key candidates avoid chunks with disproportionately large data (think sharding by states, california is big)

###### p541
- fixed and dynamic shard allocation strategies
- fixed: modulus and CRC32 are examples
- fixed requires reallocation to expand number of "buckets"
- dynamic allocation, example -> mapping table (user id to shard id)

###### p543
- can combine fixed & dynamic sharding, example take data and hash into fixed number of buckets, then dynamically assign those buckets to nodes
- explicit allocation, advantage is objectid contains shard id

###### p545
- dynamically toggle "full" property on shards in application
- also can use many replicas of a shard for improved reads
- `auto_increment_increment` and `auto_increment_offset` help generate globally unique IDs multiple servers

###### p546
- repeat that GUIDs aren't good primary keys since non-sequential (slower inserts) (also large)

###### p547
- advocate running multiple mysql instances per server to better use available resources

###### p550
- MySQL cluster: NDB database (network database) and MySQL

* cli "title bar" tool for tmux replication demo blog post

###### p552
- "akiban" - can replicate mysql and serve read queries, interleaves rows from different tables for more performant joins

###### p556
- https://www.proxysql.com/ interesting
* use mysql proxying in blog post?

###### p557
- load balancing in application potentially very efficient

###### p560
- DNS not controllable, therefore not great for load balancing 
- ARP Address Resolution Protocol for virtual IPs

###### p563
- When adding new MySQL server to load balanced setup, can mirror SELECT queries to warm cache before adding to rotation

###### p567
- 99.999% uptime = 5 mins downtime a year

###### p569
- top downtime source; operating environment -> disk space
- high availability achievement: proactively prevent downtime, prepare to recover quickly from downtime
- MTBF mean time between failures
- MTTR mean time to recovery

###### p570
- large dummy file you can delete if filesystem becomes full... lol

###### p571
- Malcolm Gladwell essay on Challenger in What the Dog Saw
- https://web.mit.edu/2.75/resources/random/How%20Complex%20Systems%20Fail.pdf

###### p572
- All downtime incidents are caused by multiple failures in combination
  - example, data loss: loss of data and not having backups

###### p574
- DRBD - duplicates disk on another machine over network, linux kernel feature

###### p576
- *favorite* way to use DRBD replicate device with binary logs

###### p576
- NDB Cluster / MySQL cluster offers a form of synchronous MySQL replication

###### p577
- NDBC can do eventual-consistency (per transaction conflict detection and resolution)

###### p578
- https://en.wikipedia.org/wiki/Optimistic_concurrency_control
- probably sexiest thing about Percona XtraDB: synchronous replication. All nodes can handle reads

###### p584
- need to read up on virtual IP addresses, Address Resolution Protocol
- STONITH - shoot the other node in the head

###### p590
- George Reese "Cloud Application Architecture" (I skimmed the TOC, looks too simple and kinda old)
- Accessing additional resources from min guaranteed capacity can make capacity planning more complex

###### p600
- AWS RDS basically MySQL, EC2, EBS (at least it was at time of writing, looks like its SSD based now)

###### p607
- persistent connections "generally bad idea" except if slow network, making so many connections running out of ports
- connection pooling and persistent connections a little confusing
- connection pools can queue requests when pool is busy (makes req wait on app server)

###### p611
- `pt-tcp-model` https://www.percona.com/doc/percona-toolkit/2.1/pt-tcp-model.html

* I should explore the whole percona toolkit suite in detail

###### p617
- reminder that cache importance can creep up as time goes on, making it a critical point of the infrastructure

* Mind wants to turn the channel a lot, is this ADD or just normal?

###### p618
- HandlerSocket, bypass upper layer of MySQL server and access innodb handler directly

###### p619
- Brian Aker has articles about storage engines
- appears to be this: https://krow.livejournal.com/513146.html

###### p621
- Curtis Preston "Backup & Recovery"
- hot, cold, warm backups - overloaded terms
- restore: fetching files and placing them in right spots, recover: entire process of rescuing a system

###### p623
- RPO recovery point objective
- RTO recovery time objective

###### p624
- replicas and RAIDs are not backups
- logical backups too slow large databases, need raw backups (good idea to periodically extract logical from raw)
- binary logs essential for PIT recovery

###### p625
- hard vs soft point-in-time recovery

###### p629
- innodb raw files bigger than logical backups because unused space

###### p630
- differential and incremental backups

###### p634
- replicas very useful for backups
- `pt-slave-delay` intentional replica delay for backups

###### p636
- good idea to have mysql itself handle purging binary logs from disk instead of using `rm`

###### p637
- logical backups: sql dumps & delimited fies

###### p640
- filesystem snapshots/backups are great

###### p656
- hardware most common inndodb corruption source, lying fsync

###### p657
- `innodb_force_recovery` https://dev.mysql.com/doc/refman/8.0/en/forcing-innodb-recovery.html

###### p661
- mysqldump can only lock and dump 1 database at a time

###### p673
- innotop cli monitoring tool

###### p680
- percona server easy to use, just stop mysql process and start percona (uses same files)

###### p682
- drizzle is abandoned
