# ddshop: the distributed database workshop

In this lab, you'll be running a Go-based web application that can store data in
both PostgreSQL and CockroachDB. Your facilitators have only tested this lab on
macOS and Linux machines, but it should be possible to run this lab on a Windows
machine as well.

## Step 1. Download `ddshop`

We've made precompiled binaries available for your convenience. Navigate to
a convenient directory and run the appropriate command for your platform.

On 64-bit macOS:

```shell
$ curl -L https://github.com/cockroachdb/ddshop/releases/download/v0.0.5/ddshop-darwin-amd64 > ddshop
$ chmod +x ddshop
```

On 64-bit Linux:

```shell
$ curl -L https://github.com/cockroachdb/ddshop/releases/download/v0.0.5/ddshop-linux-amd64 > ddshop
$ chmod +x ddshop
```

On 64-bit Windows, download the [ddshop-windows-amd64] binary.

On other platforms, you can compile from source if you have Go installed:

```shell
$ go get github.com/cockroachdb/ddshop
$ cd $(go env GOPATH)/src/github.com/cockroachdb/ddshop
$ go build
```

[ddshop-windows-amd64]: https://github.com/cockroachdb/ddshop/releases/download/v0.0.5/ddshop-windows-amd64.exe

## Step 2. Install and launch PostgreSQL (optional)

On macOS:

```shell
$ brew install postgresql
$ brew services start postgresql
```

On Linux, you'll have to consult your package manager. On Windows, don't bother
with this step or the next.

## Step 2. Launch the `ddshop` server

Launch the `ddshop` web server against PostgreSQL by running:

```shell
$ ./ddshop -postgres
```

The web application will be available at http://localhost:26256. If your
PostgreSQL server is running on a nonstandard port or under a nonstandard user,
you can explicitly specify a database URL for `ddshop` to connect to:

```shell
$ ./ddshop postgres://user@localhost:1234/specialdb
```

Spend a few moments familiarizing yourself with `ddshop` by adding, completing,
and removing todos. Close and reload your browser to convince yourself that
todos are getting saved.

## Step 3. Turn off PostgreSQL

Simulate a failure of your database machine by turning off PostgreSQL.

On macOS:

```shell
$ brew services stop postgresql
```

On Linux, the command is usually something like:

```shell
$ sudo service stop postgres
```

Now try to use the ddshop application. You'll notice that it's impossible to
perform any updates while the database is down.

## Step 5. Install and launch CockroachDB

Install CockroachDB by following the [installation
instructions][cockroachdb-install] for your platform. Then launch a three-node
cluster. In this workshop, you'll run all three nodes on one machine. In a real
deployment, you'd run each node on a separate server for fault tolerance.

1. In a new terminal, start node 1:

    ```shell
    $ ./cockroach start \
        --insecure \
        --store=node1 \
        --host=localhost \
        --port=26257 \
        --http-port=8080
    ```

2. In another terminal, start node 2:

    ```shell
    $ ./cockroach start \
        --insecure \
        --store=node2 \
        --host=localhost \
        --port=26258 \
        --http-port=8081 \
        --join=localhost:26257,localhost:26258,localhost:26259
    ```

3. In another terminal, start node 3:

    ```shell
    $ ./cockroach start \
        --insecure \
        --store=node3 \
        --host=localhost \
        --port=26259 \
        --http-port=8082 \
        --join=localhost:26257,localhost:26258,localhost:26259
    ```


## Step 6: Relaunch `ddshop` against CockroachDB

Relaunch `ddshop` against CockroachDB by passing the `-cockroach` flag:

```shell
$ ./ddshop -cockroach
```

`ddshop` will first try to contact the CockroachDB server at port 26257. If this
fails, it contacts the server at port 26258, finally falling back to the server
on port 26259.

Go back to the `ddshop` todo list and verify that it still works.

## Step 7: Check out the web UI

CockroachDB ships with a web UI that provides visibility into your cluster's
health. Check it out at http://localhost:8080.

## Step 8: Take down a CockroachDB node

Simulate a failure of the machine running node 3 by quitting that CockroachDB
instance:

```shell
$ ./cockroach quit --insecure --port=26258
```

You'll notice that the web UI has detected this failure and displays that one
node is "suspect."

Now go back to the `ddshop` todo list. You'll notice that it's still online,
even though the database has gone down!

## Step 9: Take down another CockroachDB node

Now simulate a failure of the machine running node 2 by quitting that
CockroachDB instance:

```shell
$ ./cockroach quit --insecure --port=26258
```

This time, the `ddshop` todo list will become unavailable. CockroachDB can
tolerate one node failure in its default configuration, but not two. For
particularly critical data, you can increase CockroachDB's replication to 5x to
tolerate up to two node failures, or 7x to tolerate three node failures, and so
on.

## Step 10: Bring a CockroachDB node back online

Simulate a recovery of the machine running node 2 by restarting that CockroachDB
instance:

```shell
$ ./cockroach start \
    --insecure \
    --store=node2 \
    --host=localhost \
    --port=26258 \
    --http-port=8081 \
    --join=localhost:26257,localhost:26258,localhost:26259
```

The `ddshop` todo list should come back to life.

## Concluding thoughts

Distributed databases provide you with all the flexibility and power of standard
SQL databases while being fault-tolerant and horizontally scalable.

[cockroachdb-install]: https://www.cockroachlabs.com/docs/stable/install-cockroachdb.html
