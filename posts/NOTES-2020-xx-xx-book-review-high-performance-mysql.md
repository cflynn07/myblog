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
