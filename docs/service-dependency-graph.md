Service Dependency Graph
========================

Objective
=========

Service dependency graph is used to track relationship between grpc clients and
services.

Implementation
==============

Storage structure
-----------------

The dependency is recorded by SkyLB when the grpc client calls
NewServiceCli(callerServiceId) and Resolve(serviceSpec). The data is stored in
etcd, in the path with following format:

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
/skylb/graph/<namespace>/<serviceName1>/clients/
                                               |_<clientServiceName1>/timestamp1
                                               |_<clientServiceName2>/timestamp2
                                               |_...
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The timestamp is to achieve some level of "versioning": when client-server
dependencies changes over time, the calculated dependency will be up-to-date,
instead of stacking up obsolete information. Also, with a reasonable TTL, the
obsolete dependencies will vanish eventually. See more explanations under “Query
service by client” section.

An example of the stored data looks like:

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
/skylb/graph/default/recency-service/clients/windows-client
1478926168

/skylb/graph/default/idm-service/clients/windows-client
1478926168
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Queries
-------

### Query clients by service

Simply list `/skylb/graph/<namespace>/<serviceName1>/clients/`, and the client
names will be obtained. Note that these may contain obsolete clients which no
longer rely on this service. If it is required to filter the obsolete clients,
similar process to the following section will help:

### Query service by client

To obtain a service-to-client dependency graph - and an up-to-date one, it is
needed to traverse `/skylb/graph/` to get all server-to-client pairs, and for
each pair relating to a given client, only keep the ones with the latest
timestamp, filtering the old ones. In this way we can eliminate the obsolete
information.

#### Example:

Time 1: client 1 depends on server 1.

Graph is stored like below, here ignoring the prefix /skylb/graph:

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
s1 -> clients
        |_ c1 (t1)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Calculated result:

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
c1 -> s1
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Time 2: client 1 depends on service 1 and service 2:

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
s1 -> clients
        |_ c1 (t2)
s2 -> clients
        |_ c1 (t2)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Calculated result:

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
c1 -> s1
   -> s2
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Time 3: client 1 depends only on service 2, not service 1.

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
s1 -> clients
        |_ c1 (t2)
s2 -> clients
        |_ c1 (t3)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Now for pair `s1 -> c1(t2)` and `s2 -> c1(t3)`, we know the latter is more
up-to-date, so `s1 -> c1(t2)` will be filtered, and the result will be:

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
c1 -> s2
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
