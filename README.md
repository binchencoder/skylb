# Overview

SkyLB: an external gRPC load balancer based on SkyDNS. The design follows the gRPC
Load Balancer Architecture (https://github.com/grpc/grpc/blob/master/doc/load-balancing.md).
SkyLB can be used both in and out of Kubernetes.

## Developer Guide

### Prepare Source Code

```bash
$ git clone https://github.com/binchencoder/skylb.git
```

### Build with Bazel

```bash
$ bazel build skylb/...
```

### Start Dev Run with Docker Compose

Docker compose is our recommended way to start dev run of SkyLB. So first make
sure you installed docker and docker compose on your machine.

Check out our docker projects:

```bash
$ sc track bld_tools/docker/etcd
$ sc track bld_tools/docker/ubuntu
```

Next you should download Ubuntu base package ubuntu-base-16.04-core-amd64.tar
from www.ubuntu.com (as described in bld_tools/docker/ubuntu/BUILD) and put
it in folder bld_tools/docker/ubuntu.

Run the following command to build docker images:

(bazel version must be >= **0.4.4** in order for docker-compose to work properly)

```bash
$ bazel run bld_tools/docker/etcd:latest
$ bazel run skylb/cmd/skylb:latest
$ bazel run skylb/cmd/webserver:latest
```

Run docker compose to start the dev run:

```bash
$ docker-compose -f docker-compose/dev/skylb/docker-compose.yml up
```

