Go MySQL Replay
===============

Replays statements from a traffic dump or captures from Performance Schema.

WARNING: This tool will execute operations that can change your data if those
 operations are in your dump file. Don't run it against your production DB!

Requirements
============

* MySQL
* go (if building from source)
* tshark (optional, to capturing and convert network traffic)
* tcpdump (optional, to capture network traffic)
* Access to the network interface to capture.
* Access to `performance_schema`. (if capturing from PS)

Note that wireshark/tshark should be able to decode SSL/TLS if you give it
 your key and configure MySQL to not use Diffie-Hellman.

Workflow
========

1. Capture data

```
# tcpdump -i eth0 -w dump.pcap port 3306
```

2. Convert your data to a tab dilimtered file with tshark
3. Replay the statements

```
$ ./go-mysql-replay -f my_workload.dat
```

To combine steps 2 and 3:

    $ sudo tshark -i lo -Y mysql.query -d tcp.port==5709,mysql -e tcp.stream -e frame.time_epoch -e mysql.query -Tfields
    Running as user "root" and group "root". This could be dangerous.
    Capturing on 'Loopback'
    0	1445166898.745198000	select @@version_comment limit 1
    0	1445166898.745338000	SELECT VERSION()
    0	1445166898.745516000	SELECT CURRENT_TIMESTAMP
    1	1445166923.496890000	select @@version_comment limit 1
    1	1445166923.497021000	SELECT VERSION()
    1	1445166923.497140000	SELECT CURRENT_TIMESTAMP
    ^C1 packet dropped
    6 packets captured

Using Performance Schema
========================

You need to enable the consumer for one or more events statements tables. You
might want to adjust `performance_schema_events_statements_history_long_size`
to control how many statements you're capturing

    UPDATE setup_consumers SET ENABLED='YES' WHERE NAME='events_statements_history_long';

The to capture the statements:

    SELECT THREAD_ID, TIMER_START*10e-13, SQL_TEXT FROM events_statements_history_long
    WHERE SQL_TEXT IS NOT NULL INTO OUTFILE '/tmp/statements_from_ps.dat';

TODO
====

* The parameters for the MySQL connection are hardcoded.
* Better handle the default database of each connection (hardcoded now).
