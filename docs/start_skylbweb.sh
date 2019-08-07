#!/bin/sh

cd `dirname $0`

export ETCDCTL_API=3
# etcd endpoints:
export ENDPOINTS=http://192.168.38.6:2377,http://192.168.38.6:2377
# skylb endpoints
export SKYLB=skylbserver:1900,skylbserver:1900

# Uncomment the following two lines to modify listen port and metrics port.
#PORT_OPT=" -host-port=:8091 "
#METRICS_OPT=" -scrape-port=8093 "

LOGDIR=`pwd`/logs
mkdir -p $LOGDIR
bin/skylbweb -etcd-endpoints="$ENDPOINTS" --skylb-endpoints=$SKYLB $PORT_OPT $METRICS_OPT -v 3 -stderr-to-file -log_dir=$LOGDIR \
  &

