[My notes from reading MySQL Crash Course](https://gist.github.com/cflynn07/41b79ad0c2fbc5d2ec4475f6aeefefba)

<img src="/static/images/mysql-crash-course.png" style="width:50%;">

"MySQL Crash Course" by Ben Forta is a somewhat outdated, yet still useful
overview of MySQL functionality for a beginner or someone with experience that
wants a review. The book was originally published in July 2006, during the
heydays of MySQL 5.0. 14 years later the latest version of MySQL is 8, many
things have changed but plenty remains the same.

<a href="/static/images/mysql-crash-course-screenshot.png">
  <img src="/static/images/mysql-crash-course-screenshot.png" style="width:100%;">
</a>
<span class="font-weight-bold">My tmux setup for reading MySQL Crash Course</span>

I have a few years of experience working with MySQL and relational databases
(SQL from 10 years ago that makes me
[cringe](https://github.com/cflynn07/clubbingowl/blob/master/main/application/models/model_team_guest_lists.php#L508)
today) and therefore I didn't come to this book with a beginner's perspective.
I wanted an overview of MySQL to serve as a "warm up" for my planned subsequent
adventure into <span class="font-weight-bold">"High Performance MySQL"</span>
by Baron Schwartz, Peter Zaitsev, and Vadim Tkachenko.  Therefore, for me,
this book was useful as it helped guide a focused review over core features
of MySQL. For every topic the book touched on, I frequently embarked on a
Google side quest to read supplementary blog posts or documentation. The
sample database schema and example queries were useful. I found that executing
the queries against my own MySQL server as I read along with the text led to
experimentations on variations of the queries and further bolstered my
comprehension. 

The book is about 280 pages. The first 100 pages cover the basics of SELECT,
INSERT, UPDATE and DELETE. Anyone that isn't completely new to MySQL or RDBMS's
likely already has a good grasp of these. I breezed through this part of the
book in one sitting. My pace through the latter half of the book was much
slower due to all the side reading for each topic. Being a "crash course" book
that's only ~280 pages, each topic is minimally covered. For example, chapter
13 "Grouping Data" suggests using a WHERE clause to find `orders` placed within
the last 9 months, but doesn't provide an example for how to do this. The lack
of an example led to one of my side quests to the MySQL documentation on data
and time functions
([date-and-time-functions](https://dev.mysql.com/doc/refman/8.0/en/date-and-time-functions.html#function_subdate))
to gather necessary the necessary knowledge to actually write the query. Since
the provided sample data contains datetime values from 2005, I used an interval
of 15 years instead of 9 months.

<pre class="prettyprint">
$ date
Fri Jan 31 19:16:15 CST 2020
$ cat<<-EOF | mysql -u root -ppassword -h 127.0.0.1 crashcourse 2>/dev/null
  SELECT DATE_SUB(CURDATE(), INTERVAL 9 MONTH) AS 9_months_ago;
EOF 
--------------
SELECT DATE_SUB(CURDATE(), INTERVAL 9 MONTH) AS 9_months_ago
--------------
+--------------+
| 9_months_ago |
+--------------+
| 2019-05-09   |
+--------------+

-- The Query I came up with (example not provided)
-- Since all the example data orders are from 2005, I can't use
-- 9 months so I use 20 years
SELECT cust_id, COUNT(*) AS orders
FROM orders
WHERE DATE(order_date) > DATE_SUB(CURDATE(), INTERVAL 20 YEAR)
GROUP BY cust_id
HAVING COUNT(*) >= 2;
</pre>

The book does a good job of touching on full-text searching, views, stored
procedures, cursors, triggers and transactions. Reviewing these spurred my mind
to think of specific examples of problems encountered at previous startups
where we as a team could have implemented solutions using the database layer
but instead for whatever reasons missed these opportunities and implemented
more convoluted and less robust solutions in the application layer.

Being about 14 years old, some of the book examples are outdated. For instance,
an example on page 167 suggests aliasing a column to 'rank' - however in MySQL
8 'rank' is a reserved word.
<pre class="prettyprint">
-- No longer works, 'rank' is a reserved word
SELECT note_text,
       Match(note_text) Against('rabbit') AS rank
FROM productnotes;
</pre>

And the newest features of MySQL 5.1-5.6 and MySQL 8 are naturally not covered,
such as document stores (designed to counter popularity of NoSQL databases like
mongodb). And by virtue of being a crash course, other topics like referential
integrity and foreign keys are also not covered.

There are better books out there that benefit simply by being newer. "MySQL 8
Cookbook" by Karthik Appigatla appears to be a superior text, judging by the
table of contents.
