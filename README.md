Go MySQL Replay
===============

Replays statements from a traffic dump.

WARNING: This tool will execute operations that can change your data if those
 operations are in your dump file. Don't run it against your production DB!

Workflow
========

1. Capture data

    # tcpdump -i eth0 -w dump.pcap port 3306

2. Convert your data to a tab dilimtered file with tshark
3. Replay the statements

    $ ./go-mysql-replay -f my_workload.dat

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

TODO
====

* The parameters for the MySQL connection are hardcoded.
