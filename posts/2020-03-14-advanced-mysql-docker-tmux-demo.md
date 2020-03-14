Recently I finished reading "High Performance MySQL" by Baron Schwartz, Peter
Zaitsev and Vadim Tkachenko. I write a
[review](/posts/2020-03-13-book-review-high-performance-mysql) of my thoughts on the book.

While I read technical books I like to dabble in code and go off on "side
quests" with Google for additional information. HPM covers topics including
benchmarking (io, cpu, OLTP workloads), system statistic information
monitoring, replication, etc. I got the idea to use docker, tmux, and
tmuxinator to set up a few "one-click" demos of working with these.

Tmuxinator is a useful tool. It lets you define tmux sessions consisting of
windows and panes in yaml. One cli command and you've got a fully set up tmux
session.

##### Setting up Master-Slave MySQL replication with Docker
[cflynn07/hpm-sandbox](https://github.com/cflynn07/hpm-sandbox)
<script id="asciicast-wMJ22UmOaahU0XO60OGRvPDDt" src="https://asciinema.org/a/wMJ22UmOaahU0XO60OGRvPDDt.js" async></script>

The above demo is a tmux session, initialized by tmuxinator, that shows two
docker containers running MySQL in a master-slave replication pair. The top two
tmux panes show the master's logs, and commands that are run to initialize the
server as the master. The bottom two panes show the slave and the commands to
initialize it as a replica.

###### The YML definition of the tmux session
<pre class="prettyprint">
# replication/tmuxinator-replication.yml

name: tmuxinator-replication
windows:
  - demo:
      root: /root
      layout: 83ed,214x120,0,0[214x59,0,0{107x59,0,0,0,106x59,108,0,2},214x60,0,60{107x60,0,60,4,106x60,108,60,5}]
      panes:
        - while true; do docker logs -f mysql_master; sleep 1; done
        - ./master-setup.sh
        - while true; do docker logs -f mysql_slave; sleep 1; done
        - ./slave-setup.sh
</pre>

First, `replication/start.sh` uses docker-compose to start the two MySQL containers.
###### docker-compose.yml
<pre class="prettyprint">
version: "3.7"
services:
  mysql_master:
    container_name: mysql_master
    image: "mysql:8"
    environment:
      MYSQL_ROOT_PASSWORD: password
    expose:
      - 3306
    command: --server-id=1 --log-bin=master-bin.log
  mysql_slave:
    container_name: mysql_slave
    image: "mysql:8"
    environment:
      MYSQL_ROOT_PASSWORD: password
    expose:
      - 3306
    command: --server-id=2 --relay-log-index=slave-relay-bin.index --relay-log=slave-relay-bin
</pre>
Note the `command:` values. Normally `mysqld` looks to a `my.cnf` file for
configuration but it's also possible to pass configuration values to the binary
as CLI arguments. These CLI arguments set the server IDs and enable binary
logging which is required for replication.

Next, `start.sh` creates a new container with the docker CLI client, tmux and
tmuxinator and connects it to the network created by docker-compose so it can
communicate with the two servers. It also mounts the dockerd socket into the
container's filesystem so the container's docker client can communicate with
the host's docker daemon.

###### tmux/tmuxinator container
<pre class="prettyprint">
$ docker run \
  --rm \
  -it \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v "$(pwd):/root" \
  --network "${DIR}_default" \
  --name tmuxinator \
  tmuxinator:latest \ #built and tagged locally
  tmuxinator start -p /root/tmuxinator-replication.yml
</pre>

The tmuxinator container connects to the master and the slave servers and
configures each with a bash script.

###### Master initialization commands
<pre class="prettyprint">
#!/bin/bash
# replication/master-setup.sh

set -e
set -f #noglob

. ./shared.sh

shared::init "mysql_master"
shared::wait_on_mysql
shared::query "CREATE USER 'repl_user'@'%' IDENTIFIED WITH mysql_native_password BY 'password'"
shared::query "GRANT REPLICATION SLAVE ON *.* TO repl_user@'%'"
shared::query "CREATE DATABASE repl_demo"
shared::query "CREATE TABLE repl_demo.table1 (id int not null primary key)"

count=1
while true; do
  shared::query "INSERT INTO repl_demo.table1 VALUES ($count)"
  count=$((count + 1))
  sleep 1
done
</pre>

###### Slave initialization commands
<pre class="prettyprint">
#!/bin/bash
# replication/slave-setup.sh

# set -e
set -f #noglob

. ./shared.sh

shared::init "mysql_slave"
shared::wait_on_mysql
shared::query "CHANGE MASTER TO MASTER_HOST='mysql_master',
  MASTER_USER='repl_user',
  MASTER_PASSWORD='password',
  MASTER_LOG_FILE='master-bin.log',
  MASTER_LOG_POS=0"

while true; do
  shared::query "SELECT * FROM repl_demo.table1 ORDER BY id DESC LIMIT 10"
  sleep 1
done
</pre>

The end result is an easy to run and visualize demo of multiple services
running in a docker managed network. I'm planning on creating more demos like
this to visualize and document setting up services, running tests, etc.
