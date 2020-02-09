Notes
- deprecation of mysql query browser for mysql workbench
- use of mycli
- p65 tips for using wildcards, avoid wildcard at start of string - performance
- p91 SOUDNEX, never heard of this before
-
- p102 COUNT(*) all rows, COUNT(column) only rows with NON-NULL values for column
- p113 "WITH ROLLUP"

- p114 * WHERE does not work with groups (WHERE doesn't know what a group is)
  - HAVING <-- this works on groups, I didn't previously know this
  - WHERE filters ROWS, HAVING filters GROUPSA
    - p115 WHERE filters before data is grouped, HAVING filters after data is grouped

# Get me all the customers that have at least 2 orders
SELECT cust_id, COUNT(*) AS orders
FROM orders
GROUP BY cust_id
HAVING COUNT(*) >= 2;

---
p115 book suggests using WHERE clause to get all orders past 6 months, but doesn't provide example.
So I figured it out myself...
# Found this documentation on DATE_SUB/SUBDATE function
# https://dev.mysql.com/doc/refman/8.0/en/date-and-time-functions.html#function_subdate
- The DATE_SUB function is not listed in the functions on page 93
- According to W3 schools it's bene around since MySQL 4
`SELECT DATE_SUB(CURDATE(), INTERVAL 9 MONTH) AS 9_months_ago`
# Getting the date 9 months ago
```
$ date
Fri Jan 31 19:16:15 CST 2020
$ echo "SELECT DATE_SUB(CURDATE(), INTERVAL 9 MONTH) AS 9_months_ago;" | mysql -u root -ppassword -h 127.0.0.1 crashcourse 2>/dev/null
9_months_ago
2019-04-30
```
# The Query I came up with (example not provided)
# Since all the example data orders are from 2005, I can't use
# 9 months so I use 20 years
```
SELECT cust_id, COUNT(*) AS orders
FROM orders
WHERE DATE(order_date) > DATE_SUB(CURDATE(), INTERVAL 20 YEAR)
GROUP BY cust_id
HAVING COUNT(*) >= 2;
```
---

p117 has good comparison table of ORDER BY and GROUP BY
RULE: anytime GROUP BY, use ORDER BY

---
p123
```
SELECT cust_id
FROM orders
WHERE order_num IN (SELECT order_num
                    FROM orderitems
                    WHERE prod_id = 'TNT2');
```
This returns:
cust_id
10001
10004
My thought here was, what if customer 10001 had 2 orders that
had the item TNT2. To test this I looked at the sample data and
saw customer 10001 has another order, 20009, in `orders`. That
order doesn't contain an orderitem of TNT2 but we can pretend it does:
```
SELECT cust_id
FROM orders
WHERE order_num IN (20005,20007,20009); //20005, 20007 belong to 10001 and 10004 respectively
```
returns:
cust_id
10001
10001
10004

is the solution to use DISTINCT? Will the book mention this in the next few pages?
Interestingly, when using this as a subquery to a SELECT on customers, customer data
is not duplicated
```
SELECT cust_name, cust_contact
FROM customers
WHERE cust_id IN (10001,10001,10004);

cust_name       cust_contact
Coyote Inc.     Y Lee
Yosemite Place  Y Sam
```

p129 shows using subquery to calculate a field, mentions this might not always
be post efficient. Notes later chapters will discuss.

---
p 143 aliases for table names, example uses `AS`: tablename AS tn
  - I never really used this, I wonder if dropping AS has always been supported
three other joins, self join, natural join, outer join -- this is where my need
for review comes in





p145 To think about self join, I literally printed out two copies of the table,
placed them side by side, and worked my way down the line with a pen to my screen
testing like the DBMS
ps <3 heredoc
````
cat <<- EOF | mysql -u root -ppassword -h 127.0.0.1 -v --table crashcourse
SELECT prod_id, vend_id, prod_name, prod_id, vend_id, prod_name
FROM products;
EOF
```
--------------
SELECT prod_id, vend_id, prod_name, prod_id, vend_id, prod_name FROM products
--------------

+---------+---------+----------------+---------+---------+----------------+
| prod_id | vend_id | prod_name      | prod_id | vend_id | prod_name      |
+---------+---------+----------------+---------+---------+----------------+
| ANV01   |    1001 | .5 ton anvil   | ANV01   |    1001 | .5 ton anvil   |
| ANV02   |    1001 | 1 ton anvil    | ANV02   |    1001 | 1 ton anvil    |
| ANV03   |    1001 | 2 ton anvil    | ANV03   |    1001 | 2 ton anvil    |
| DTNTR   |    1003 | Detonator      | DTNTR   |    1003 | Detonator      |
| FB      |    1003 | Bird seed      | FB      |    1003 | Bird seed      |
| FC      |    1003 | Carrots        | FC      |    1003 | Carrots        |
| FU1     |    1002 | Fuses          | FU1     |    1002 | Fuses          |
| JP1000  |    1005 | JetPack 1000   | JP1000  |    1005 | JetPack 1000   |
| JP2000  |    1005 | JetPack 2000   | JP2000  |    1005 | JetPack 2000   |
| OL1     |    1002 | Oil can        | OL1     |    1002 | Oil can        |
| SAFE    |    1003 | Safe           | SAFE    |    1003 | Safe           |
| SLING   |    1003 | Sling          | SLING   |    1003 | Sling          |
| TNT1    |    1003 | TNT (1 stick)  | TNT1    |    1003 | TNT (1 stick)  |
| TNT2    |    1003 | TNT (5 sticks) | TNT2    |    1003 | TNT (5 sticks) |
+---------+---------+----------------+---------+---------+----------------+
when the book started touching on self joins this is when I paused to start
looking at youtube videos and blog posts talking about visualizing self joins
# https://www.mysqltutorial.org/mysql-self-join/ 
great tutorial from MySQL

typo on page 147 (OI)

Went to bed, woke up again (Sunday morning) re-reviewed self-joins and somehow
I feel like it makes a little bit more sense now.

```
SELECT p1.prod_id, p1.prod_name
FROM products AS p1, products AS p2
WHERE p1.vend_id = p2.ved_id
  AND p2.prod_id = 'DTNTR'; <-- there will be only 1 record from p2 compared to every row in p1*
```

p 157 haven't written a UNION/UNION ALL query in a long time...
* UNION ALL does what can't be done w/ multiple WHERE clauses

----
page 167 example is outdated
example query suggests aliasing generated column to 'rank' -- in MySQL 8 'RANK' is a reserved word
# https://dev.mysql.com/doc/refman/8.0/en/mysql-nutshell.html
"Window functions.  MySQL now supports window functions that, for each row from
a query, perform a calculation using rows related to that row. These include
functions such as RANK(), LAG(), and NTILE(). In addition, several existing
aggregate functions now can be used as window functions (for example, SUM() and
AVG()). For more information, see Section 12.21, “Window Functions”."

