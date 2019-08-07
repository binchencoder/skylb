#!/bin/sh
cd `dirname $0`

export ETCDCTL_API=3
export ENDPOINTS=http://192.168.64.176:2381,http://192.168.64.176:2383,http://192.168.64.176:2385

LOGDIR=`pwd`/logs
mkdir -p $LOGDIR

/skylb/skylbweb \
    -etcd-endpoints="$ENDPOINTS" -v 2 -log_dir=$LOGDIR &

wait
