##### Tracking progress, in general
  
I like measuring progress. I'm better at staying focused and motivated if I can
quantify my rate of progression when working towards goals.
  
I first figured this out in high school. One of my friends introduced me to the
practice of *writing down* my exercise sets and reps at the gym. Before this,
my approach to working out was: **go to the gym, do a bunch of random stuff for
a while, repeat.** I didn't plan ahead what exercises I would perform or how
much weight or how many reps I would try to lift, and I didn't measure how much
I lifted at each session. I had no idea if I was getting stronger or just
treading water. I started bringing a notebook and a pen with me to the gym. I
planned out 5 different workouts, each focusing on different body areas,  and
rotated through them. I wrote down how many sets and how many reps at what
weight I was able to do for each exercise, and each time I went to the gym I'd
try to lift slightly more than I was able to the previous time. The gain in
efficiency was substantial. Every time I went to the gym I had a plan for what
workouts and sets to perform; no more randomly picking exercises. It was great
to be able to flip through the notebook and actually observe progression.
  
In 2020 I'm still measuring and planning, except instead of a notebook I'm
using the iOS app [strong][1] which basically does the same thing, except it
has all the goodness of being 21st century software like cloud storage and
progress analytics.
  
##### Tracking reading progress
I tend to think of reading and comprehending books on technical topics as
another form of "working out." The strain is on my neural networks' ability to
form connections rather than my muscle fibers' ability to contract. I've ported
the same progress measuring approach from the gym to my reading.
  
A few months ago I bought a dozen or so technical books and I'm diligently
working through all of them. Books are underrated. Many of the topics covered
in books could probably be learned through reading a series of online blog
posts (I'm very confident I've learned more from blog posts than college), but
books have the advantage of pulling together a lot of information and
organizing it for more efficient learning. Also, you can flex your reading
accomplishments with physical books by creating an awesome looking bookshelf.
My apartment benefits from the added color.
  
Like in the gym, measuring and tracking reading progress has benefits for focus
and motivation. When you actually track your reading progress you go from some
impalpable, unquantifiable sense of progress to knowing things like: *I read #x
amount over the last 5 days, with #y pages per day and at this rate it will
take me #z days to finish this book.* You get a little dopamine hit from
quantifiable achievement.
  
I'm currently reading [Fullstack React][2]. I decided to break the book up
into segments of 10 pages; you can think of these as being akin to gym "sets."
I record the start and finish time of each set as I read. I started out just
jotting down measurements on paper; eventually I transitioned to using a
spreadsheet.
  
![](/static/images/fullstack-react-time-tracking.png)
[My Spreadsheet on Google Docs][3]
  
Columns `A` through `F` are my page block (~"reps") tracking input area. For
each block of 10 pages I read I input `A` the date, `B` block start, `C` Block
end, `D` start time, `E` end time and the duration `F` is automatically
calculated.
  
![](/static/images/fullstack-react-time-tracking-1.png)
  
Just as a nice visual aid, I apply alternating colors to the `A` column rows by
date with using a conditional formating formula. The formula applies the green
color depending on whether the day of the month is even or odd (this means if I
were to skip a day the alternation would be a little messed up).
  
<pre class="prettyprint">
=and(iseven(day(A2)),isblank(A2)=false)
</pre>

![](/static/images/fullstack-react-time-tracking-2.png)
  
The fun part of having this data in a tabular format is you can write queries
on it to get analytics. Google Sheets has its own flavor of SQL which is
similar enough to "real" databases like MySQL. You can write queries that use
aggregator functions and SQL statements like `GROUP BY` to find out stats like,
how much total time per day spent reading, average time per block per day,
total number of blocks per day, etc. Unfortunately, Google Sheet's aggregator
functions `MAX`, `SUM` and `AVG` don't work with non-numeric values like time
durations, so I came up with a bit of a workaround by converting all of the
duration values in column `F` into seconds in column `G` (which I've hidden). I
then ran a query using column `G`.
  
The values in column `G` are the values in the "duration" column `F` converted to
seconds displayed as an integer using this function:
<pre class="prettyprint">
=if(ISBLANK(F3)," ",F3*24*3600)
</pre>
Multiplying the value in `F` converts the duration to float representation of the
fraction of 1 day, so 0:30:00 minutes becomes `0.5/24=0.02083333333`.
Multiplying by `24*3600` converts to seconds.
  
Columns `I` through `L` display the output of the following query:
<pre class="prettyprint" style="white-space: pre-wrap">
=QUERY(A3:G, "SELECT A, (AVG(G)/86400), (SUM(G)/86400), (MAX(C)-MIN(B)) GROUP BY A LABEL A 'Date', (AVG(G)/86400) 'Average Time/Block', (SUM(G)/86400) 'Total Time', (MAX(C)-MIN(B)) 'Total Pages'")
</pre>

![](/static/images/fullstack-react-time-tracking-3.png)
  
Now as I input my reading progress, I can get a sense for how well I'm doing
each day relative to other days.
  
By using the output of the previous query I can get a few more stats, such as
total days, total time, etc. In my sheet I've placed these in columns `M` and `N`.
  
###### Total Days:
<pre class="prettyprint">
=count(unique(A2:A100))
</pre>
  
###### Total Time:
<pre class="prettyprint">
=SUM(F2:F)
</pre>
  
###### Total Average Time/Block:
<pre class="prettyprint">
=AVERAGE(F2:F)
</pre>
  
###### Average Total Time/Day:
<pre class="prettyprint">
=N4/counta(unique(A2:A100))
</pre>
  
###### Average Pages/Day:
<pre class="prettyprint">
=FLOOR(SUM(L:L)/counta(unique(A2:A100)))
</pre>
  
###### Extrapolated Days Remaining:
<pre class="prettyprint">
=ROUND((N2-MAX(C:C))/N7,2)
</pre>
  
###### Extrapolated Completion Date:
<pre class="prettyprint">
=MAX(A:A) + CEILING(N8)
</pre>
  
[1]: https://www.strong.app/
[2]: https://www.newline.co/fullstack-react/
[3]: https://docs.google.com/spreadsheets/d/1yr6W_tK-W4uowAwJnWekOfjmmG9ZbJV-OlCnasdZ72g/edit?usp=sharing
