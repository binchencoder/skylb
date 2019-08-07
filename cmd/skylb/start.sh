#!/bin/sh
cd `dirname $0`

/skylb/skylb \
    -within-k8s=$WITHIN_K8S \
    -etcd-endpoints="$ENDPOINTS" \
    -v $LOG_LEVEL -log_dir=$LOG_DIR &

wait
