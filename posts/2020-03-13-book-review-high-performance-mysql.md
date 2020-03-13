[My notes from reading High Performance MySQL](https://gist.github.com/cflynn07/46c564935607c444e2258b23f490168f)

<img src="/static/images/high-performance-mysql.png" style="width:50%;">

I've just finished reading "High Performance MySQL" by Baron Schwartz, Peter
Zaitsev and Vadim Tkachenko. The printed version of the book is approximately
790 pages. I actually read the whole thing, daily, sequentially... and it took
me over a month. I feel like I've got a million data points and thoughts from
the book swirling around in my head. In this post I want to focus on some of
the content but also my thoughts on digesting dense technical books.

I'm working my way through a pile of technical books as part of a new years
resolution to "level up" in software engineering. I'm concentrating on my
ability to concentrate. Over the past month, while reading this book, I
progressively got better at catching myself "changing the channel" in my brain.
Whenever some thought or thing would enter my brain while I was working on
downloading the information in this book into my mind I would actively fight
back and restore my concentration. And over the course of a month, I got a lot
better at it. I subjectively feel like my mind tries less often to "change the
channel" and when it does I find it easier to overcome the urge. My learning
pace quickened as I progressed through the book and I'm eager to see how I
perform when I dive into the next book in my queue.

<img src="/static/images/high-performance-mysql-time-blocks.png" style="width:50%;">

One trick I stumbled upon was a simple way to quantify my progress. I read the
book in blocks of 10 pages and I marked the start and end time for every block.
My average time spent per block of 10 pages was approximately 1 hour. I started
noticing I could get through blocks of 10 pages in 30-45 minutes in the final
quarter of the book. To read the entire book and comprehend it well probably
took on the order of 70-80 hours.

I took a pretty comprehensive database class back in college and I'm no
stranger to RDBMSs but my work as a "full stack" web application software
engineer over the past few years has precluded me from doing a really "deep
dive" in MySQL for some time. High Performance MySQL is comprehensive and well
organized. From reading it I've substantially increased my [tree trunk of
understanding](https://waitbutwhy.com/2016/05/mailbag-1.html) around
RDBMSs/MySQL.

The third edition was published in 2012, around the time of MySQL 5.5, so some
of the information is out of date. I found myself constantly googling claims
made throughout the book thinking "is that still the case?" Things move quickly
in the world of software so often the answer to my question was: no. That being
said, much of the information is still good even 8 years after publication.
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
That's in large part to the detailed explanations of b-tree+ indexes, covering
indexes, hash indexes, clustering indexes, composite indexes, query parsing,
index selectivity, selective denormalization strategies, MySQL statistics and
logging, etc.

The book even recommends other books for guidance on how to go deeper. One book
the book recommends that I plan to read is "Relational Database Index Design
and the Optimizers" by Mike Leach and Tapio Lahdenmaki. I particularly thought
the "3 star" ranking system for quantifying the utility of indexes was
interesting.

###### 3 Star Index Ranking
- 1 star: places relevant rows adjacent to each other
- 2 star: rows are sorted in the order the query needs
- 3 star: contains all columns needed (a "covering" index)

The book covers scaling in great detail. Both scaling up and scaling out
(vertical/horizontal scaling) through replication, partitioning and sharding.
With intelligent adaptation of these tools, the authors explain how MySQL can
scale to match the read/write capacity needs of even Facebook.


# the covering of how to scale up & out (partitioning and sharding
