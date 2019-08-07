#!/bin/sh

export ETCDCTL_API=3
export ENDPOINTS=http://192.168.64.176:2381,http://192.168.64.176:2383,http://192.168.64.176:2385

LOGDIR=`pwd`/logs
mkdir -p $LOGDIR
bin/skylb -etcd-endpoints="$ENDPOINTS" -v 3 -stderr-to-file -log_dir=$LOGDIR &

