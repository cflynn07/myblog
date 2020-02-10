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

